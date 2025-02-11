package tg

import (
	"context"
	"slices"
	"strings"
)

func CommonFilterCommand(command string) FilterFunc {
	const botCommandEntity = "bot_command"
	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}

	return func(ctx context.Context, update *Update) bool {
		if update.Message == nil || len(update.Message.Entities) == 0 && len(update.Message.CaptionEntities) == 0 {
			return false
		}

		var commandEntity *MessageEntity
		pred := func(entity *MessageEntity) bool { return entity != nil && entity.Type == botCommandEntity }
		if pos := slices.IndexFunc(update.Message.Entities, pred); pos != -1 {
			commandEntity = update.Message.Entities[pos]
		} else if pos = slices.IndexFunc(update.Message.Entities, pred); pos != -1 {
			commandEntity = update.Message.Entities[pos]
		}

		if commandEntity == nil {
			return false
		}
		offset := commandEntity.Offset
		length := commandEntity.Length
		return update.Message.Text[offset:offset+length] == command
	}
}

func CommonTextReply(text string, asReply ...bool) HandlerFunc {
	isReply := at(asReply, 0, false)
	return func(ctx context.Context, update *Update) error {
		opts := []*OptSendMessage{}
		if isReply {
			opts = append(opts, &OptSendMessage{ReplyParameters: &ReplyParameters{
				MessageId:                update.Message.MessageId,
				AllowSendingWithoutReply: true,
			}})
		}
		_, err := SendMessage(ctx, update.Message.Chat.Id, text, opts...)
		return err
	}
}

func OnMessage(ctx context.Context, update *Update) bool {
	return update != nil && update.Message != nil
}

func OnText(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Text != ""
}

func OnPhoto(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && len(update.Message.Photo) != 0
}

func OnVideo(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Video != nil
}

func OnAnimation(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Animation != nil
}

func OnDocument(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Document != nil
}

func OnVideoNote(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.VideoNote != nil
}

func OnVoice(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Voice != nil
}

func OnAudio(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Audio != nil
}

func OnSticker(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Sticker != nil
}

func IsPrivateMessage(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Chat != nil && update.Message.SenderChat != nil && update.Message.Chat.Id == update.Message.SenderChat.Id
}

func IsPublicMessage(ctx context.Context, update *Update) bool {
	return OnMessage(ctx, update) && update.Message.Chat != nil && update.Message.SenderChat != nil && update.Message.Chat.Id != update.Message.SenderChat.Id
}

func Or(fn ...FilterFunc) FilterFunc {
	return func(ctx context.Context, update *Update) bool {
		for _, f := range fn {
			if f(ctx, update) {
				return true
			}
		}
		return false
	}
}

func All(fn ...FilterFunc) FilterFunc {
	return func(ctx context.Context, update *Update) bool {
		for _, f := range fn {
			if !f(ctx, update) {
				return false
			}
		}
		return true
	}
}
