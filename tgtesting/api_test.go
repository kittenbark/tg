package tgtesting

import (
	"context"
	"github.com/kittenbark/tg"
	"net/http"
	"testing"
)

func TestOK(t *testing.T) {
	t.Parallel()

	t.Run("getMe",
		MakeTestOK(
			&tg.User{Id: 123456, Username: "testbot"},
			tg.GetMe,
			Stub{Url: "/getMe", Result: StubResultOK(http.StatusOK, &tg.User{Id: 123456, Username: "testbot"})},
		),
	)

	t.Run("sendPhoto",
		MakeTestOK(
			&tg.Message{
				MessageId: 123,
				Chat:      &tg.Chat{Id: 123456},
				Photo: []*tg.PhotoSize{{
					FileId:   "abc123",
					Width:    600,
					Height:   800,
					FileSize: 1024,
				}}},
			func(ctx context.Context) (*tg.Message, error) {
				result, err := tg.SendPhoto(ctx, int64(123), tg.FromCloud("photo_id"))
				return result, err
			},
			Stub{Url: "/sendPhoto", Result: StubResultOK(http.StatusOK, &tg.Message{
				MessageId: 123,
				Chat:      &tg.Chat{Id: 123456},
				Photo: []*tg.PhotoSize{{
					FileId:   "abc123",
					Width:    600,
					Height:   800,
					FileSize: 1024,
				}}}),
			},
		),
	)
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("getMe#429",
		MakeTestError(
			tg.GetMe,
			Stub{Url: "/getMe", Result: StubResultError(http.StatusTooManyRequests,
				"too many requests",
				"retry_after", "5.66",
			)},
		),
	)
}

func MakeTestOK[Expected any, Request any](expected *Expected, request func(ctx context.Context) (*Request, error), stubs ...Stub) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := NewTestingEnv(t, &Config{
			Stubs: stubs,
		})
		result, err := request(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	}
}

func MakeTestError[Request any](request func(ctx context.Context) (*Request, error), stubs ...Stub) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := NewTestingEnv(t, &Config{
			Stubs: stubs,
		})
		_, err := request(ctx)
		require.Error(t, err)
	}
}
