package tgtesting

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kittenbark/tg"
)

func TestRedactToken(t *testing.T) {
	var gotBody string

	cfg := &Config{
		Stubs: []Stub{
			{
				Url: "/sendMessage",
				Validator: func(req *http.Request) bool {
					b, _ := io.ReadAll(req.Body)
					gotBody = string(b)
					return true
				},
				Result: StubResultOK(200, &tg.Message{MessageId: 1, Chat: &tg.Chat{Id: 1}}),
			},
		},
	}
	ctx := NewTestingContext(t, cfg)
	token := cfg.Token

	t.Run("enabled by default", func(t *testing.T) {
		_, err := tg.SendMessage(ctx, 1, "leaked: "+token)
		require.NoError(t, err)
		require.Equal(t, false, strings.Contains(gotBody, token))
		require.Equal(t, true, strings.Contains(gotBody, "123456:**************"))
	})

	t.Run("disabled explicitly", func(t *testing.T) {
		disabledCtx := context.WithValue(ctx, tg.ContextRedactToken, false)
		_, err := tg.SendMessage(disabledCtx, 1, "leaked: "+token)
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(gotBody, token))
	})
}
