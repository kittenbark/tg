package tg

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
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
func Chain(handlerFunc ...HandlerFunc) HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		for _, handler := range handlerFunc {
			if err := handler(ctx, upd); err != nil {
				return err
			}
		}
		return nil
	}
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

func OnCommand(command string) FilterFunc {
	const botCommandEntity = "bot_command"
	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}

	return func(ctx context.Context, upd *Update) bool {
		if !OnMessage(ctx, upd) {
			return false
		}

		var commandEntity *MessageEntity
		pred := func(entity *MessageEntity) bool { return entity != nil && entity.Type == botCommandEntity }
		if pos := slices.IndexFunc(upd.Message.Entities, pred); pos != -1 {
			commandEntity = upd.Message.Entities[pos]
		} else if pos = slices.IndexFunc(upd.Message.Entities, pred); pos != -1 {
			commandEntity = upd.Message.Entities[pos]
		}

		if commandEntity == nil {
			return false
		}
		offset := commandEntity.Offset
		length := commandEntity.Length
		return upd.Message.Text[offset:offset+length] == command
	}
}

func OnMessage(ctx context.Context, upd *Update) bool {
	return upd != nil && upd.Message != nil
}

func OnPrivateMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Chat != nil && upd.Message.From != nil &&
		upd.Message.Chat.Id == upd.Message.From.Id
}

func OnPublicMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Chat != nil && upd.Message.From != nil &&
		upd.Message.Chat.Id != upd.Message.From.Id
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

func OnUrl() FilterFunc {
	// Warning: this was generated by Claude 3.7, because I hate regexes. Enjoy.
	const pattern = `^(https?|ftp):\/\/` + // Protocol (http, https, ftp)
		`(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*` + // Subdomains
		`([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])` + // Domain name
		`(:\d+)?` + // Port (optional)
		`(\/[-a-zA-Z0-9\._~:/?#[\]@!$&'()*+,;=%]*)??` // Path (optional)
	return OnTextRegexp(pattern)
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
