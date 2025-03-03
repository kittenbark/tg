package tgtesting

import (
	"context"
	"errors"
	"github.com/kittenbark/tg"
	"math/rand/v2"
	"net/http"
	"sync"
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

	ctx := NewTestingContext(t, &Config{
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

type CounterPlugin struct {
	Calls map[string]int
	lock  *sync.Mutex
}

func (plugin *CounterPlugin) Hooks() []tg.PluginHookType {
	return []tg.PluginHookType{
		tg.PluginHookOnHandleStart,
		tg.PluginHookOnHandleFinish,
		tg.PluginHookOnFilter,
		tg.PluginHookOnUpdate,
		tg.PluginHookOnError,
	}
}

func (plugin *CounterPlugin) Apply(ctx tg.PluginHookContext) {
	plugin.lock.Lock()
	defer plugin.lock.Unlock()

	switch ctx.(type) {
	case *tg.PluginHookContextOnUpdate:
		plugin.Calls["update"]++
	case *tg.PluginHookContextOnFilter:
		plugin.Calls["filter"]++
	case *tg.PluginHookContextOnHandleStart:
		plugin.Calls["handle_start"]++
	case *tg.PluginHookContextOnHandleFinish:
		plugin.Calls["handle_finish"]++
	case *tg.PluginHookContextOnError:
		plugin.Calls["error"]++
	}
}

func TestPlugins(t *testing.T) {
	SetTestingEnv(t, &Config{
		Stubs: []Stub{
			{
				Url: "/getUpdates",
				Result: StubResultOK(200, []*tg.Update{
					{
						UpdateId: 1,
						Message:  &tg.Message{MessageId: 1, Chat: &tg.Chat{Id: 1}, Text: "testtext"},
					},
					{
						UpdateId: 2, CallbackQuery: &tg.CallbackQuery{Id: "callback_query", Data: `{"id": 1, "value": "other"}`},
					},
				}),
			},
			{
				Url:    "/sendMessage",
				Result: StubResultOK(200, &tg.Message{MessageId: 2, Chat: &tg.Chat{Id: 1}, Text: "testtext"}),
			},
		},
	})

	counter := &CounterPlugin{
		Calls: map[string]int{
			"update":        0,
			"filter":        0,
			"handle_start":  0,
			"handle_finish": 0,
			"error":         0,
		},
		lock: &sync.Mutex{},
	}

	bot := tg.NewFromEnv()
	start := time.Now()
	bot.
		Plugin(counter).
		Filter(tg.OnMessage).
		Handle(func(ctx context.Context, upd *tg.Update) error {
			if time.Since(start).Seconds() > 1 {
				bot.Stop()
			}
			if rand.Float32() > 0.5 {
				return errors.New("random error")
			}
			return nil
		}).
		Start()

	for hook := range counter.Calls {
		require.Equal(t, true, counter.Calls[hook] > 0)
	}
}

func MakeTestOK[Expected any, Request any](expected *Expected, request func(ctx context.Context) (*Request, error), stubs ...Stub) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := NewTestingContext(t, &Config{
			Stubs: stubs,
		})
		result, err := request(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	}
}

func MakeTestError[Request any](request func(ctx context.Context) (*Request, error), stubs ...Stub) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := NewTestingContext(t, &Config{
			Stubs: stubs,
		})
		_, err := request(ctx)
		require.Error(t, err)
	}
}
