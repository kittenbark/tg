package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log/slog"
	"sync"
)

type ButtonI interface {
	Build() *InlineKeyboardButton
	HandlerFunc() HandlerFunc
}

var (
	_ ButtonI = (*CallbackButton)(nil)
	_ ButtonI = (*InlineKeyboardButton)(nil)
)

type Button = InlineKeyboardButton

func (b *Button) Build() *Button { return b }

func (b *Button) HandlerFunc() HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		slog.Warn("tg.Button.HandlerFunc was called, use tg.CallbackButton if you need a button with a callback")
		return nil
	}
}

type CallbackButton struct {
	Text       string
	Handler    HandlerFunc
	OnComplete *OptAnswerCallbackQuery
}

func (c *CallbackButton) Build() *InlineKeyboardButton {
	return &InlineKeyboardButton{
		Text: c.Text,
	}
}

func (c *CallbackButton) HandlerFunc() HandlerFunc {
	return func(ctx context.Context, upd *Update) error {
		if err := c.Handler(ctx, upd); err != nil {
			return err
		}
		if c.OnComplete != nil && upd.CallbackQuery != nil {
			query := upd.CallbackQuery
			_, _ = AnswerCallbackQuery(ctx, query.Id, c.OnComplete)
		}
		return nil
	}
}

type Keyboard struct {
	Layout [][]ButtonI
	idOnce sync.Once
	id     string
}

func (k *Keyboard) init() {
	k.idOnce.Do(func() {
		result := make([][]*InlineKeyboardButton, len(k.Layout))
		for i, buttons := range k.Layout {
			result[i] = make([]*InlineKeyboardButton, len(buttons))
			for j, button := range buttons {
				result[i][j] = button.Build()
			}
		}
		data, _ := json.Marshal(result)
		hash := fnv.New32()
		_, _ = hash.Write(data)
		k.id = fmt.Sprintf("%x", hash.Sum32())
	})
}

type keyboardSchema struct {
	KeyboardId string `json:"K"`
	ButtonId   int    `json:"B"`
	Data       string `json:"D"`
}

func (k *Keyboard) Branch() (FilterFunc, HandlerFunc) {
	return k.FilterFunc(), k.HandlerFunc()
}

func (k *Keyboard) Build() *InlineKeyboardMarkup {
	k.init()
	result := make([][]*InlineKeyboardButton, len(k.Layout))
	buttonId := 0
	for i, buttons := range k.Layout {
		result[i] = make([]*InlineKeyboardButton, len(buttons))
		for j, button := range buttons {
			buttonId++
			built := button.Build()
			data, _ := json.Marshal(keyboardSchema{k.id, buttonId, built.Text})
			built.CallbackData = string(data)
			result[i][j] = built
		}
	}
	return &InlineKeyboardMarkup{
		InlineKeyboard: result,
	}
}

func (k *Keyboard) FilterFunc() FilterFunc {
	k.init()
	return All(OnCallback, func(ctx context.Context, upd *Update) bool {
		var data keyboardSchema
		if err := json.Unmarshal([]byte(upd.CallbackQuery.Data), &data); err != nil {
			return false
		}
		return k.id == data.KeyboardId
	})
}

func (k *Keyboard) HandlerFunc() HandlerFunc {
	k.init()
	buttonId := 0
	hashToButton := map[int]HandlerFunc{}
	for _, buttons := range k.Layout {
		for _, button := range buttons {
			buttonId++
			hashToButton[buttonId] = button.HandlerFunc()
		}
	}
	return func(ctx context.Context, upd *Update) error {
		var data keyboardSchema
		if err := json.Unmarshal([]byte(upd.CallbackQuery.Data), &data); err != nil {
			return err
		}
		handler, ok := hashToButton[data.ButtonId]
		if !ok {
			return fmt.Errorf("unknown button %d", data.ButtonId)
		}
		if upd.CallbackQuery != nil {
			upd.CallbackQuery.Data = data.Data
		}
		return handler(ctx, upd)
	}
}
