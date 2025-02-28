package tgtesting

import (
	"context"
	"github.com/kittenbark/tg"
	"net/http"
	"testing"
	"time"
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

func TestGetUpdatesCallbackQuery(t *testing.T) {
	t.Parallel()

	ctx := NewTestingEnv(t, &Config{
		Stubs: []Stub{
			{Url: "/getUpdates", Result: StubResultOK(200, []*tg.Update{
				{
					UpdateId: 123456,
					CallbackQuery: &tg.CallbackQuery{
						Id:              "abc123",
						From:            &tg.User{Id: 123456, Username: "testuser"},
						Message:         &tg.Message{Chat: &tg.Chat{Id: 123456, Username: "testchat"}, Text: "testtext"},
						InlineMessageId: "inlineabc123",
						ChatInstance:    "testinsance",
						Data:            ``,
					},
				},
				{
					UpdateId: 123457,
					CallbackQuery: &tg.CallbackQuery{
						Id:              "def456",
						From:            &tg.User{Id: 123456, Username: "testuser"},
						Message:         &tg.Message{Chat: &tg.Chat{Id: 123456, Username: "testchat"}, Text: "testtext"},
						InlineMessageId: "inlineabc123",
						ChatInstance:    "testinsance",
						Data:            `{"id": 123, "value": "testvalue"}`,
					},
				},
				{
					UpdateId: 123457,
					CallbackQuery: &tg.CallbackQuery{
						Id:              "def456",
						From:            &tg.User{Id: 123456, Username: "testuser"},
						Message:         &tg.Message{Chat: &tg.Chat{Id: 123456, Username: "testchat"}, Text: "testtext"},
						InlineMessageId: "inlineabc123",
						ChatInstance:    "testinsance",
						Data:            `{"id": 123, "value": "other"}`,
					},
				},
			})},
		},
	})

	type CallbackData struct {
		Id    int    `json:"id"`
		Value string `json:"value"`
	}

	bot := tg.New(&tg.Config{
		Token:        ctx.Value(tg.ContextToken).(string),
		ApiURL:       ctx.Value(tg.ContextApiUrl).(string),
		SyncHandling: true,
	})
	time.AfterFunc(time.Second*5, bot.Stop)
	isOther := func(value *CallbackData) bool {
		return value.Value == "other"
	}
	met1, met2, met3, met4 := 0, 0, 0, 0
	bot.
		OnError(tg.OnErrorLog).
		Branch(tg.OnCallbackWithData[CallbackData](isOther), func(ctx context.Context, upd *tg.Update) error {
			met1++
			value, err := tg.CallbackData[CallbackData](upd)
			if err != nil {
				return err
			}
			require.Equal(t, 123, value.Id)
			require.Equal(t, "other", value.Value)
			return nil
		}).
		Branch(tg.OnCallbackWithData[CallbackData](), func(ctx context.Context, upd *tg.Update) error {
			met2++
			value, err := tg.CallbackData[CallbackData](upd)
			if err != nil {
				return err
			}
			require.Equal(t, 123, value.Id)
			require.Equal(t, "testvalue", value.Value)
			return nil
		}).
		Branch(tg.OnCallback, func(ctx context.Context, upd *tg.Update) error {
			met3++
			require.Equal(t, int64(123456), upd.UpdateId)
			require.Equal(t, "testchat", upd.CallbackQuery.Message.Chat.Username)
			return nil
		}).
		Handle(func(ctx context.Context, upd *tg.Update) error {
			met4++
			return nil
		}).
		Start()

	require.Equal(t, true, met1*met2*met3 > 0)
	require.Equal(t, true, met4 == 0)
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

func TestRegexes(t *testing.T) {
	t.Parallel()

	onText := tg.OnTextRegexp("kitten")
	require.True(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kitten"}}))
	require.True(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kittenbark"}}))
	require.False(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kit"}}))

	onUrl := tg.OnUrl()
	require.True(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "http://google.com"}}))
	require.False(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "google.com"}}))
	require.True(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "https://google.com?url=me"}}))
	require.False(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "kittenbark"}}))
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
