package tg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Chain allows convenient chaining of handlers.
// Example:
//
//	tg.NewFromEnv().
//		OnError(tg.OnErrorLog).
//		Branch(tg.OnVideo, tg.Chain(
//			tg.CommonReactionReply("👀"),
//			tg.Synced(SomeHeavyVideoConverterHandler),
//			tg.CommonReactionReply("👌")),
//		).
//		Start()
func Chain(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		for _, handler := range handlers {
			if err := handler(ctx, upd); err != nil {
				return err
			}
		}
		return nil
	}
}

func Fallback(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx context.Context, upd *Update) (err error) {
		for _, handler := range handlers {
			if err = handler(ctx, upd); err == nil {
				return nil
			}
		}
		return err
	}
}

func FallbackWithMessage(handler HandlerFunc, msg string) HandlerFunc {
	return Fallback(handler, CommonTextReply(msg, true))
}

// Synced wraps handler making it, ugh, synced.
func Synced(handlerFunc HandlerFunc) HandlerFunc {
	mutex := &sync.Mutex{}
	return func(ctx context.Context, upd *Update) error {
		mutex.Lock()
		defer mutex.Unlock()
		return handlerFunc(ctx, upd)
	}
}

type SyncedGroup struct {
	mutex sync.Mutex
}

func (group *SyncedGroup) Synced(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		group.mutex.Lock()
		defer group.mutex.Unlock()
		return handler(ctx, upd)
	}
}

func Either(fn ...FilterFunc) FilterFunc {
	return func(ctx context.Context, upd *Update) bool {
		for _, f := range fn {
			if f(ctx, upd) {
				return true
			}
		}
		return false
	}
}

func All(fn ...FilterFunc) FilterFunc {
	return func(ctx context.Context, upd *Update) bool {
		for _, f := range fn {
			if !f(ctx, upd) {
				return false
			}
		}
		return true
	}
}

func Not(fn FilterFunc) FilterFunc {
	return func(ctx context.Context, upd *Update) bool {
		return !fn(ctx, upd)
	}
}

func CommonTextReply(text string, asReply ...bool) HandlerFunc {
	isReply := at(asReply, 0, false)
	return func(ctx context.Context, upd *Update) error {
		opts := []*OptSendMessage{}
		if isReply {
			opts = append(opts, &OptSendMessage{ReplyParameters: &ReplyParameters{
				MessageId:                upd.Message.MessageId,
				AllowSendingWithoutReply: true,
			}})
		}
		_, err := SendMessage(ctx, upd.Message.Chat.Id, text, opts...)
		return err
	}
}

func CommonReaction(emoji string, big ...bool) *OptSetMessageReaction {
	return &OptSetMessageReaction{
		Reaction: []ReactionType{&ReactionTypeEmoji{Type: "emoji", Emoji: emoji}},
		IsBig:    at(big, 0, false),
	}
}

func CommonReactionReply(emoji string, big ...bool) HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		_, err := SetMessageReaction(ctx, upd.Message.Chat.Id, upd.Message.MessageId, CommonReaction(emoji, big...))
		return err
	}
}

func CommonDeleteMessage(ctx context.Context, upd *Update) error {
	if upd == nil || upd.Message == nil || upd.Message.Chat == nil {
		return errors.New("tg#CommonDeleteMessage: bad update, no upd/message/chat")
	}
	_, err := DeleteMessage(ctx, upd.Message.Chat.Id, upd.Message.MessageId)
	return err
}

type CommonArgsHandlerFunc[T any] func(ctx context.Context, upd *Update, args T) error

func CommonArgs[T any](fn CommonArgsHandlerFunc[*T]) HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		fields := slices.Concat(strings.Fields(upd.Message.Text), strings.Fields(upd.Message.Caption))[1:]
		args := new(T)
		val := reflect.ValueOf(args).Elem()
		for i := range min(val.NumField(), len(fields)) {
			arg := val.Field(i)
			switch arg.Kind() {
			case reflect.String:
				arg.SetString(fields[i])
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				parsed, err := strconv.ParseInt(fields[i], 10, 64)
				if err != nil {
					return err
				}
				arg.SetInt(parsed)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				parsed, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return err
				}
				arg.SetUint(parsed)
			case reflect.Float32, reflect.Float64:
				parsed, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					return err
				}
				arg.SetFloat(parsed)
			case reflect.Bool:
				parsed, err := strconv.ParseBool(fields[i])
				if err != nil {
					return err
				}
				arg.SetBool(parsed)
			default:
				return fmt.Errorf("tg#CommonArgs: unexpected field type %s", arg.Kind().String())
			}
		}

		return fn(ctx, upd, args)
	}
}

func OnCommand(command string) FilterFunc {
	const botCommandEntity = "bot_command"
	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}

	return func(ctx context.Context, upd *Update) bool {
		if !OnText(ctx, upd) || len(upd.Message.Entities) == 0 {
			return false
		}

		for _, entity := range upd.Message.Entities {
			if entity == nil || entity.Type != botCommandEntity {
				continue
			}

			entityText, _, _ := strings.Cut(
				upd.Message.Text[entity.Offset:entity.Offset+entity.Length],
				"@",
			)

			if entityText == command {
				return true
			}
		}

		return false
	}
}

func OnPrivate(ctx context.Context, upd *Update) bool {
	return OnPrivateMessage(ctx, upd) ||
		upd.EditedMessage != nil && isMessagePrivate(upd.EditedMessage) ||
		upd.MessageReaction != nil && upd.MessageReaction.Chat != nil && upd.MessageReaction.User != nil && upd.MessageReaction.Chat.Id == upd.MessageReaction.User.Id ||
		upd.MessageReactionCount != nil && upd.MessageReactionCount.Chat != nil && upd.MessageReaction.User != nil && upd.MessageReaction.Chat.Id == upd.MessageReaction.User.Id ||
		upd.InlineQuery != nil && upd.InlineQuery.ChatType == "private"
}

func OnMessage(ctx context.Context, upd *Update) bool {
	return upd != nil && upd.Message != nil
}

func OnPrivateMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && isMessagePrivate(upd.Message)
}

func OnPublicMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Chat != nil && upd.Message.From != nil &&
		upd.Message.Chat.Id != upd.Message.From.Id
}

func OnChat(chatId ...int64) FilterFunc {
	chats := map[int64]struct{}{}
	for _, id := range chatId {
		chats[id] = struct{}{}
	}
	return func(ctx context.Context, upd *Update) bool {
		_, ok := chats[getChatId(upd)]
		return ok
	}
}

func OnSender(senderId ...int64) FilterFunc {
	whitelist := map[int64]struct{}{}
	for _, id := range senderId {
		whitelist[id] = struct{}{}
	}
	return func(ctx context.Context, upd *Update) bool {
		_, ok := whitelist[getSenderId(upd)]
		return ok
	}
}

// OnChannelPostInCommentsChat uses OnSender with the id of the hidden Telegram channel->comments reposting bot.
func OnChannelPostInCommentsChat() FilterFunc { return OnSender(777000) }

func OnForwarded(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.ForwardOrigin != nil
}

func OnAutomaticForward(ctx context.Context, upd *Update) bool {
	return upd != nil && upd.Message != nil && upd.Message.IsAutomaticForward
}

func OnReply(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.ReplyToMessage != nil
}

func OnEdited(ctx context.Context, upd *Update) bool {
	return upd != nil && (upd.EditedMessage != nil || upd.EditedChannelPost != nil || upd.EditedBusinessMessage != nil)
}

func OnChance(chance float64) FilterFunc {
	return func(ctx context.Context, upd *Update) bool {
		return rand.Float64() < chance
	}
}

func OnAddedToGroup(ctx context.Context, upd *Update) bool {
	if upd == nil || upd.Message == nil {
		return false
	}

	msg := upd.Message
	if msg.GroupChatCreated || msg.SupergroupChatCreated {
		return true
	}
	newMembers := msg.NewChatMembers
	if len(newMembers) == 0 {
		return false
	}

	token, err := tryGetTokenFromContext(ctx)
	if err != nil {
		slog.Warn("tg.OnAddedToGroup#no_token_in_context", "err", err)
		return false
	}
	identifier, _, _ := strings.Cut(token, ":")
	id, _ := strconv.ParseInt(identifier, 10, 64)

	for _, newMember := range newMembers {
		if newMember != nil && newMember.Id == id {
			return true
		}
	}
	return false
}

func OnText(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Text != ""
}

func OnPhoto(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && len(upd.Message.Photo) != 0
}

func OnVideo(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Video != nil
}

func OnAnimation(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Animation != nil
}

func OnDocument(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Document != nil
}

func OnVideoNote(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.VideoNote != nil
}

func OnVoice(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Voice != nil
}

func OnAudio(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Audio != nil
}

func OnMedia(ctx context.Context, upd *Update) bool {
	return OnAnimation(ctx, upd) || OnAudio(ctx, upd) || OnDocument(ctx, upd) ||
		OnPhoto(ctx, upd) || OnVideoNote(ctx, upd) || OnVideo(ctx, upd) || OnVoice(ctx, upd)
}

func OnSticker(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Sticker != nil
}

func OnTextRegexp(regex string) FilterFunc {
	re := regexp.MustCompile(regex)
	return func(ctx context.Context, upd *Update) bool {
		return OnText(ctx, upd) && re.MatchString(upd.Message.Text)
	}
}

func OnUrl(ctx context.Context, upd *Update) bool {
	if !OnText(ctx, upd) {
		return false
	}

	URL, err := url.Parse(upd.Message.Text)
	if err != nil || URL == nil {
		return false
	}
	if URL.Host == "" {
		return false
	}
	return true
}

func CallbackData[T any](upd *Update) (*T, error) {
	if upd == nil || upd.CallbackQuery == nil {
		return nil, errors.New("tg: callback query is nil")
	}

	var value T
	if err := json.Unmarshal([]byte(upd.CallbackQuery.Data), &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func OnCallback(ctx context.Context, upd *Update) bool {
	return upd != nil && upd.CallbackQuery != nil
}

func OnCallbackWithData[T any](pred ...func(value *T) bool) FilterFunc {
	predicate := at(pred, 0, func(value *T) bool { return true })
	return func(ctx context.Context, upd *Update) bool {
		if upd == nil || upd.CallbackQuery == nil {
			return false
		}

		value, err := CallbackData[T](upd)
		if err != nil {
			return false
		}
		return predicate(value)
	}
}

func OnChatJoinRequest(ctx context.Context, upd *Update) bool {
	return upd.ChatJoinRequest != nil
}

func OnNewChatMember(filters ...func(user *User) bool) FilterFunc {
	return func(ctx context.Context, upd *Update) bool {
		if upd == nil || upd.Message == nil || len(upd.Message.NewChatMembers) == 0 {
			return false
		}
		if len(filters) == 0 {
			return true
		}
		// For any users & all filters.
		for _, user := range upd.Message.NewChatMembers {
			ok := true
			for _, filter := range filters {
				if !filter(user) {
					ok = false
					break
				}
			}
			if ok {
				return true
			}
		}
		return false
	}
}

func AsReplyTo(msg *Message) *ReplyParameters {
	if msg == nil {
		return nil
	}
	return &ReplyParameters{MessageId: msg.MessageId}
}

// HandleEditedAsMessage basically reuses upd.EditedMessage as upd.Message,
// useful for getting track of any message changes in chats.
// It's a bit hacky:
//
//	tg.NewFromEnv().
//		Filter(tg.HandleEditedAsMessage).
//		...your handlers, filters, branches...
func HandleEditedAsMessage(ctx context.Context, upd *Update) bool {
	if upd != nil {
		if upd.EditedMessage != nil && upd.Message == nil {
			upd.Message = upd.EditedMessage
		}
		if upd.EditedBusinessMessage != nil && upd.BusinessMessage == nil {
			upd.BusinessMessage = upd.EditedBusinessMessage
		}
	}
	return true
}

func (impl *Message) TextOrCaption() string {
	switch {
	case impl == nil:
		return ""
	case impl.Text != "":
		return impl.Text
	case impl.Caption != "":
		return impl.Caption
	default:
		return ""
	}
}

func isMessagePrivate(msg *Message) bool {
	return msg.Chat != nil && msg.From != nil && msg.Chat.Id == msg.From.Id
}

func getChatId(upd *Update) int64 {
	defer func() { _ = recover() }()

	switch {
	case upd == nil:
		return 0
	case upd.Message != nil:
		return upd.Message.Chat.Id
	case upd.EditedMessage != nil:
		return upd.EditedMessage.Chat.Id

	case upd.ChannelPost != nil:
		return upd.ChannelPost.Chat.Id
	case upd.EditedChannelPost != nil:
		return upd.EditedChannelPost.Chat.Id

	case upd.InlineQuery != nil:
		return upd.InlineQuery.From.Id
	case upd.ChosenInlineResult != nil:
		return upd.ChosenInlineResult.From.Id
	case upd.CallbackQuery != nil:
		return upd.CallbackQuery.From.Id
	case upd.ShippingQuery != nil:
		return upd.ShippingQuery.From.Id
	case upd.PreCheckoutQuery != nil:
		return upd.PreCheckoutQuery.From.Id

	case upd.ChatJoinRequest != nil:
		return upd.ChatJoinRequest.Chat.Id

	case upd.MyChatMember != nil:
		return upd.MyChatMember.Chat.Id
	case upd.ChatMember != nil:
		return upd.ChatMember.Chat.Id

	case upd.ChatBoost != nil:
		return upd.ChatBoost.Chat.Id
	case upd.RemovedChatBoost != nil:
		return upd.RemovedChatBoost.Chat.Id

	case upd.MessageReaction != nil:
		return upd.MessageReaction.Chat.Id
	case upd.MessageReactionCount != nil:
		return upd.MessageReactionCount.Chat.Id

	case upd.PurchasedPaidMedia != nil:
		return upd.PurchasedPaidMedia.From.Id
	case upd.BusinessMessage != nil:
		return upd.BusinessMessage.Chat.Id
	case upd.EditedBusinessMessage != nil:
		return upd.EditedBusinessMessage.Chat.Id
	case upd.DeletedBusinessMessages != nil:
		return upd.DeletedBusinessMessages.Chat.Id
	}
	return 0
}

func getSenderId(upd *Update) int64 {
	defer func() { _ = recover() }()

	switch {
	case upd == nil:
		return 0
	case upd.Message != nil:
		return upd.Message.From.Id
	case upd.EditedMessage != nil:
		return upd.EditedMessage.From.Id

	case upd.ChannelPost != nil:
		return upd.ChannelPost.From.Id
	case upd.EditedChannelPost != nil:
		return upd.EditedChannelPost.From.Id

	case upd.InlineQuery != nil:
		return upd.InlineQuery.From.Id
	case upd.ChosenInlineResult != nil:
		return upd.ChosenInlineResult.From.Id
	case upd.CallbackQuery != nil:
		return upd.CallbackQuery.From.Id
	case upd.ShippingQuery != nil:
		return upd.ShippingQuery.From.Id
	case upd.PreCheckoutQuery != nil:
		return upd.PreCheckoutQuery.From.Id

	case upd.ChatJoinRequest != nil:
		return upd.ChatJoinRequest.From.Id

	case upd.MyChatMember != nil:
		return upd.MyChatMember.From.Id
	case upd.ChatMember != nil:
		return upd.ChatMember.From.Id

	case upd.ChatBoost != nil:
		return upd.ChatBoost.Chat.Id
	case upd.RemovedChatBoost != nil:
		return upd.RemovedChatBoost.Chat.Id

	case upd.MessageReaction != nil:
		return upd.MessageReaction.User.Id
	case upd.MessageReactionCount != nil:
		return upd.MessageReactionCount.Chat.Id

	case upd.PurchasedPaidMedia != nil:
		return upd.PurchasedPaidMedia.From.Id
	case upd.BusinessMessage != nil:
		return upd.BusinessMessage.From.Id
	case upd.EditedBusinessMessage != nil:
		return upd.EditedBusinessMessage.From.Id
	case upd.DeletedBusinessMessages != nil:
		return upd.DeletedBusinessMessages.Chat.Id
	}
	return 0
}
