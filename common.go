package tg

import (
	"context"
	"math/rand/v2"
	"slices"
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

func OnPrivateMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Chat != nil && upd.Message.From != nil &&
		upd.Message.Chat.Id == upd.Message.From.Id
}

func OnPublicMessage(ctx context.Context, upd *Update) bool {
	return OnMessage(ctx, upd) && upd.Message.Chat != nil && upd.Message.From != nil &&
		upd.Message.Chat.Id != upd.Message.From.Id
}

func OnChance(chance float64) bool {
	return rand.Float64() < chance
}
