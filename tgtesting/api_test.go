package tgtesting

import (
	"context"
	"errors"
	"fmt"
	"github.com/kittenbark/tg"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
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
	hg := &tg.SyncedGroup{}
	bot.
		OnError(tg.OnErrorLog).
		Branch(tg.OnCallbackWithData[CallbackData](isOther), hg.Synced(func(ctx context.Context, upd *tg.Update) error {
			met1++
			value, err := tg.CallbackData[CallbackData](upd)
			if err != nil {
				return err
			}
			require.Equal(t, 123, value.Id)
			require.Equal(t, "other", value.Value)
			return nil
		})).
		Branch(tg.OnCallbackWithData[CallbackData](), hg.Synced(func(ctx context.Context, upd *tg.Update) error {
			met2++
			value, err := tg.CallbackData[CallbackData](upd)
			if err != nil {
				return err
			}
			require.Equal(t, 123, value.Id)
			require.Equal(t, "testvalue", value.Value)
			return nil
		})).
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

func TestStartAndStop(t *testing.T) {
	SetTestingEnv(t, &Config{
		Stubs: []Stub{
			{
				Url: "/getUpdates",
				Result: StubResultOK(200, []*tg.Update{{
					UpdateId: 1,
					Message:  &tg.Message{MessageId: 1, Chat: &tg.Chat{Id: 1}, Text: "testtext"},
				}}),
			},
			{
				Url:    "/sendMessage",
				Result: StubResultOK(200, &tg.Message{MessageId: 2, Chat: &tg.Chat{Id: 1}, Text: "testtext"}),
			},
		},
	})

	bot := tg.NewFromEnv()
	time.AfterFunc(time.Millisecond*10, bot.Stop)
	bot.Start()
	time.AfterFunc(time.Millisecond*10, bot.Stop)
	bot.Start()
	time.AfterFunc(time.Millisecond*10, bot.StopImmediately)
	bot.Start()

	var i int = 0
	start := time.Now()
	sent := &atomic.Int64{}
	bot.
		OnError(tg.OnErrorPanic).
		Branch(tg.OnMessage, func(ctx context.Context, upd *tg.Update) error {
			if time.Since(start).Seconds() > 0.05*float64(i+1) {
				bot.Stop()
			}

			msg := upd.Message
			_, err := tg.SendMessage(ctx, msg.Chat.Id, time.Now().String())
			sent.Add(1)
			return err
		}).
		Start()

	for i = range 10 {
		bot.Start()

		require.Geq(t, 1, sent.Load(), fmt.Sprintf("i=%d", i))
		sent.Store(0)
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
			if time.Since(start).Seconds() > 0.05 {
				bot.Stop()
			}
			if rand.Float32() > 0.5 {
				return errors.New("random error")
			}
			return nil
		}).
		Start()

	for _, count := range counter.Calls {
		require.Equal(t, true, count > 0)
	}
}

func TestCommonArgs(t *testing.T) {
	type Args struct {
		String string
		Int    int
		Bool   bool
		Float  float64
	}
	type Param struct {
		Text     string
		Expected Args
	}
	for _, param := range []Param{
		{"/start hello 42 true 3.14", Args{"hello", 42, true, 3.14}},
		{"/start 42", Args{String: "42"}},
		{"/start", Args{}},
	} {
		t.Run(param.Text, func(t *testing.T) {
			SetTestingEnv(t, &Config{
				Stubs: []Stub{
					{
						Url: "/getUpdates",
						Result: StubResultOK(200, []*tg.Update{{
							UpdateId: 1,
							Message: &tg.Message{
								MessageId: 1,
								Chat:      &tg.Chat{Id: 1},
								Text:      param.Text,
								Entities: []*tg.MessageEntity{{
									Offset: 0,
									Length: 6,
									Type:   "bot_command",
								}},
							},
						}}),
					},
				},
			})

			bot := tg.NewFromEnv()
			time.AfterFunc(time.Millisecond*100, bot.Stop)
			counter := &atomic.Int64{}
			bot.
				OnError(tg.OnErrorPanic).
				Command("/start", tg.CommonArgs[Args](func(ctx context.Context, upd *tg.Update, args *Args) error {
					require.Equal(t, param.Expected, *args)
					counter.Add(1)
					return nil
				})).
				Start()
			require.Geq(t, 1, counter.Load())
		})
	}
}

func TestScheduler(t *testing.T) {
	SetTestingEnv(t, &Config{
		Stubs: []Stub{
			{
				Url:    "/sendMessage",
				Result: StubResultOK(200, &tg.Message{MessageId: 2, Chat: &tg.Chat{Id: 1}, Text: "testtext"}),
			},
			{
				Url: "/sendMediaGroup",
				Result: StubResultOK(200, []*tg.Message{
					{MessageId: 2, Chat: &tg.Chat{Id: 1}, Photo: tg.TelegramPhoto{}},
					{MessageId: 3, Chat: &tg.Chat{Id: 1}, Video: &tg.TelegramVideo{}},
				}),
			},
			{
				Url:    "/getMe",
				Result: StubResultOK(200, &tg.User{Id: 123}),
			},
		},
	})

	Pressure := 10 * time.Microsecond
	GlobalQuota := 30
	GlobalTimeout := time.Millisecond
	PerChatQuota := 20
	PerChatTimeout := 60 * time.Millisecond
	testTime := time.Millisecond * 180

	t.Run("global", func(t *testing.T) {
		bot := tg.NewFromEnv().Scheduler(tg.NewSchedulerVerbose(
			Pressure,
			tg.SchedulerClauseGlobal(GlobalQuota, GlobalTimeout),
			tg.SchedulerClauseChat(PerChatQuota, PerChatTimeout),
		))
		sent := atomic.Int64{}
		for i := 0; i < 256; i++ {
			go func() {
				_, err := tg.GetMe(bot.Context())
				require.NoError(t, err)
				sent.Add(1)
			}()
			time.Sleep(time.Microsecond)
		}
		time.Sleep(testTime)

		require.LessOrEqualInt(t, int64(GlobalQuota)*(1+int64(testTime/GlobalTimeout)), sent.Load())
	})

	t.Run("per_chat", func(t *testing.T) {
		bot := tg.NewFromEnv().Scheduler(tg.NewSchedulerVerbose(
			Pressure,
			tg.SchedulerClauseGlobal(GlobalQuota, GlobalTimeout),
			tg.SchedulerClauseChat(PerChatQuota, PerChatTimeout),
		))
		sent := atomic.Int64{}
		for i := 0; i < 256; i++ {
			go func() {
				_, err := tg.SendMessage(bot.Context(), -100, "testtext")
				require.NoError(t, err)
				sent.Add(1)
			}()
			time.Sleep(time.Microsecond)
		}
		time.Sleep(testTime)

		require.LessOrEqualInt(t, int64(PerChatQuota)*(1+int64(testTime/PerChatTimeout)), sent.Load())
	})

	t.Run("per_user", func(t *testing.T) {
		bot := tg.NewFromEnv().Scheduler(tg.NewSchedulerVerbose(
			Pressure,
			tg.SchedulerClauseGlobal(GlobalQuota, GlobalTimeout),
			tg.SchedulerClauseUser(PerChatQuota, PerChatTimeout),
		))
		sent := atomic.Int64{}
		for i := 0; i < 256; i++ {
			go func() {
				_, err := tg.SendMessage(bot.Context(), 100, "testtext")
				require.NoError(t, err)
				sent.Add(1)
			}()
			time.Sleep(time.Microsecond)
		}
		time.Sleep(testTime)

		require.LessOrEqualInt(t, int64(PerChatQuota)*(1+int64(testTime/PerChatTimeout)), sent.Load())
	})

	t.Run("per_chat_album", func(t *testing.T) {
		bot := tg.NewFromEnv().Scheduler(tg.NewSchedulerVerbose(
			Pressure,
			tg.SchedulerClauseGlobal(GlobalQuota, GlobalTimeout),
			tg.SchedulerClauseChat(PerChatQuota, PerChatTimeout),
		))
		sent := atomic.Int64{}
		for i := 0; i < 256; i++ {
			go func() {
				_, err := tg.SendMediaGroup(bot.Context(), -100, tg.Album{
					&tg.Photo{},
					&tg.Video{},
					&tg.Photo{},
					&tg.Video{},
				})
				require.NoError(t, err)
				sent.Add(1)
			}()
			time.Sleep(time.Microsecond)
		}
		time.Sleep(testTime)

		require.LessOrEqualInt(t, int64(PerChatQuota)*(1+int64(testTime/PerChatTimeout))/4, sent.Load())
	})
}

func TestFilters(t *testing.T) {
	t.Parallel()

	onText := tg.OnTextRegexp("kitten")
	require.True(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kitten"}}))
	require.True(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kittenbark"}}))
	require.False(t, onText(nil, &tg.Update{Message: &tg.Message{Text: "kit"}}))

	onUrl := tg.OnUrl
	require.True(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "http://google.com"}}))
	require.False(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "google.com"}}))
	require.True(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "https://google.com?url=me"}}))
	require.False(t, onUrl(nil, &tg.Update{Message: &tg.Message{Text: "kittenbark"}}))

	require.True(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "kitten /start",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 7, Length: 6}},
	}}))
	require.True(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "kitten /start@somebot",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 7, Length: 14}},
	}}))
	require.False(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "kitten /stopp@somebot",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 7, Length: 14}},
	}}))
	require.False(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{Text: "kitten start"}}))
	require.True(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "/start",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
	}}))
	require.False(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "/starting",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 0, Length: 9}},
	}}))
	require.True(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "/start@somebot",
		Entities: []*tg.MessageEntity{{Type: "bot_command", Offset: 0, Length: 14}},
	}}))
	require.False(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "/start@somebot",
		Entities: []*tg.MessageEntity{{Type: "blockquote", Offset: 0, Length: 14}},
	}}))
	require.False(t, tg.OnCommand("/start")(nil, &tg.Update{Message: &tg.Message{
		Text:     "/start",
		Entities: []*tg.MessageEntity{{Type: "blockquote", Offset: 0, Length: 9}},
	}}))

	filters := []tg.FilterFunc{
		tg.OnText,
		tg.OnUrl,
		tg.OnMessage,
		tg.OnMedia,
		tg.OnAudio,
		tg.OnVideo,
		tg.OnPhoto,
		tg.OnCallback,
		tg.OnPrivate,
		tg.OnPublicMessage,
		tg.OnChat(1),
		tg.OnPrivateMessage,
		tg.OnSticker,
		tg.OnForwarded,
		tg.OnAutomaticForward,
		tg.OnVideoNote,
		tg.OnReply,
		tg.OnEdited,
		tg.OnChatJoinRequest,
		tg.OnChance(0),
		tg.OnTextRegexp("kitten"),
		tg.OnCallbackWithData[int](),
	}
	for _, filter := range filters {
		require.False(t, filter(nil, &tg.Update{}), "func", getFuncName(filter))
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
