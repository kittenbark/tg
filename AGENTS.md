# tg — Agent Reference

`github.com/kittenbark/tg` — zero-dependency Go Telegram Bot framework. Pipeline-style update handling, code-generated Telegram API, plugin system, rate-limit scheduler. Go 1.23, no external deps.

---

## Mental model

Updates flow through a linear pipeline top-to-bottom:

```
NewFromEnv()
  .Filter(pred)       // drops non-matching updates; everything below only sees matches
  .Command("/start", h)  // shorthand for Branch(OnCommand("/start"), h)
  .Branch(pred, h)    // handles matches, lets rest fall through
  .Handle(h)          // catch-all
  .Default(h)         // runs only if no node above handled the update
  .Start()            // blocks; long-polls Telegram
```

`FilterFunc = func(ctx context.Context, upd *Update) bool`
`HandlerFunc = func(ctx context.Context, upd *Update) error`

Both are plain functions — composable, definable anywhere, no framework coupling.

---

## Bot construction

```go
tg.NewFromEnv()           // reads KITTENBARK_TG_TOKEN (or KBTG_TOKEN); panics on failure
tg.New(cfg)               // from *Config struct; panics
tg.NewFromFile(path)      // reads JSON file; panics
tg.TryNew(cfg)            // same but returns (Bot, error)
tg.TryNewFromEnv()
tg.TryNewFromFile(path)
```

### Config fields

```go
type Config struct {
    Token         string
    TokenTesting  string
    ApiURL        string            // default: https://api.telegram.org
    TimeoutHandle time.Duration     // default: 1h
    TimeoutPoll   time.Duration     // default: 100ms
    SyncHandling  bool
    DownloadType  DownloadType      // Classic / LocalMove / LocalCopy
    OnError       OnErrorFunc
    OnErrorByType string            // "ignore" | "log" | "exit"
    ExtraHeaders  map[string]string
}
```

### Env var names (both `KITTENBARK_TG_` and `KBTG_` prefixes work)

| Suffix | Maps to |
|--------|---------|
| `TOKEN` | Bot token (required) |
| `TEST_TOKEN` | Fallback test token |
| `TEST_CHAT` | Test chat ID |
| `TEST_GROUP_CHAT` | Test group chat ID |
| `API_URL` | Custom API base URL |
| `EXTRA_HEADERS` | JSON `{"Key":"Value"}` headers |
| `DOWNLOAD_TYPE` | `classic` / `local_move` / `local_copy` |
| `ON_ERROR` | `ignore` / `log` / `exit` |
| `SYNCED_HANDLE` | `true` to serialize handlers |
| `TIMEOUT_HANDLE` | seconds (default 3600) |
| `TIMEOUT_POLL` | seconds (default 0.1) |

---

## Pipeline builder methods

```go
bot.Filter(pred ...FilterFunc) *Bot
bot.Branch(pred FilterFunc, h HandlerFunc) *Bot
bot.Handle(h HandlerFunc) *Bot
bot.Default(h HandlerFunc) *Bot
bot.Command(command string, h HandlerFunc) *Bot  // "/" prefix added automatically
bot.Start(updates ...*Update)                     // blocks; optional pre-loaded updates for testing
bot.Stop()                                        // graceful stop
bot.StopImmediately()
bot.Help(commands ...string) *Bot                 // alternating: cmd, description, cmd, description ...
bot.HelpScoped(scope BotCommandScope, commands ...string) *Bot
bot.OnError(fn OnErrorFunc) *Bot
bot.Plugin(plugins ...Plugin) *Bot
bot.Scheduler(s ...Scheduler) *Bot
bot.ExtraContext(fn ...ExtraContext) *Bot
```

---

## Filters — built-in predicates

### Composition
```go
tg.All(f1, f2, ...)      // logical AND
tg.Either(f1, f2, ...)   // logical OR
tg.Not(f)                // logical NOT
```

### Message type
```go
tg.OnMessage             // any message
tg.OnPrivateMessage      // private/DM only
tg.OnPublicMessage       // group/channel
tg.OnText                // has text content
tg.OnPhoto
tg.OnVideo
tg.OnAnimation
tg.OnDocument
tg.OnVideoNote
tg.OnVoice
tg.OnAudio
tg.OnSticker
tg.OnMedia               // any media type
```

### Commands & text
```go
tg.OnCommand("/start")           // matches that command
tg.OnTextRegexp(`^\d+$`)        // regex on message text
tg.OnUrl                         // text is a valid URL
```

### Callbacks
```go
tg.OnCallback                                        // any callback query
tg.OnCallbackWithData[T](pred ...func(*T) bool)     // typed JSON callback data
```

### Chat & user
```go
tg.OnChat(chatIds ...int64)       // whitelist chats
tg.OnSender(userIds ...int64)     // whitelist senders
tg.OnPrivate                      // private message/inline query
tg.OnForwarded
tg.OnAutomaticForward
tg.OnReply
tg.OnEdited
tg.OnAddedToGroup
tg.OnNewChatMember(filters ...func(*User) bool)
tg.OnChatJoinRequest
tg.OnChannelPostInCommentsChat    // hidden Telegram channel bot (ID 777000)
```

### Misc
```go
tg.OnChance(0.5)            // random pass with probability
tg.HandleEditedAsMessage    // treat EditedMessage as Message (also a HandlerFunc)
```

---

## Handler composition

```go
tg.Chain(h1, h2, h3)                        // run in sequence; stop on first error
tg.Fallback(h1, h2, h3)                     // try in sequence; return first success
tg.FallbackWithMessage(h, "sorry")          // fallback that sends message as last resort
tg.Synced(h)                                // wrap in mutex (new mutex per call)
group := &tg.SyncedGroup{}
group.Synced(h)                             // shared mutex across multiple handlers
```

---

## Common ready-made handlers

```go
tg.CommonTextReply("text")                          // send text
tg.CommonTextReply("text", true)                    // send as reply
tg.CommonTextReplyExpiring(5*time.Second, "text")   // send then delete after duration
tg.CommonReactionReply(":ok:")                      // react with emoji alias
tg.CommonReactionReply("🔥", true)                  // big reaction
tg.CommonDeleteMessage                              // delete the triggering message
tg.CommonRestrictSender(permissions)                // restrict sender permanently
tg.CommonRestrictSenderUntil(dur, permissions)      // restrict until duration
tg.CommonArgs[MyArgs](func(ctx, upd *tg.Update, args *MyArgs) error { ... })
// CommonArgs parses command arguments into struct fields (string, int*, uint*, float*, bool)
```

Emoji aliases (see `tg.CommonReactionEmojiMap`): `":ok:"`, `":fire:"`, `":eyes:"`, `":heart:"`, 150+ more.

---

## Callback query routing

```go
// Typed data: encode any JSON-serializable struct as callback data
tg.OnCallbackWithData[MyStruct](func(v *MyStruct) bool { return v.Id == 42 })
data, err := tg.CallbackData[MyStruct](upd)

// Keyboard with built-in routing
keyboard := &tg.Keyboard{
    Layout: [][]tg.ButtonI{
        {&tg.CallbackButton{Text: "Click me", Handler: myHandler}},
    },
}
markup := keyboard.BuildRegister(ctx)           // builds + registers handlers in bot
filter, handler := keyboard.Branch()            // or use manually with bot.Branch
```

---

## API methods (131 total)

All methods are package-level functions taking `ctx context.Context` as first arg. Optional params go in a variadic `*Opt{Method}` struct.

```go
// Signature pattern
func SendMessage(ctx context.Context, chatId int64, text string, opts ...*OptSendMessage) (*Message, error)

// Common methods
tg.SendMessage(ctx, chatId, text, &tg.OptSendMessage{ParseMode: tg.ParseModeHTML})
tg.EditMessageText(ctx, text, &tg.OptEditMessageText{ChatId: chatId, MessageId: msgId})
tg.DeleteMessage(ctx, chatId, msgId)
tg.CopyMessage(ctx, chatId, fromChatId, msgId)
tg.SendPhoto(ctx, chatId, photo, opts...)
tg.SendVideo(ctx, chatId, video, opts...)
tg.SendDocument(ctx, chatId, doc, opts...)
tg.GetChat(ctx, chatId)
tg.RestrictChatMember(ctx, chatId, userId, permissions, opts...)
tg.SetMessageReaction(ctx, chatId, msgId, opts...)
tg.AnswerCallbackQuery(ctx, callbackQueryId, opts...)
tg.GetFile(ctx, fileId)
```

Parse modes: `tg.ParseModeHTML`, `tg.ParseModeMarkdown`, `tg.ParseModeMarkdownV2`

Escape helpers:
```go
tg.Md("text with *special* chars")     // escape for MarkdownV2
tg.HTML("<b>text</b>")                  // escape for HTML
tg.EscapeParseMode(mode, text)          // generic
```

---

## File input/output

```go
// Upload from disk
tg.FromDisk("/path/to/file.jpg")
tg.FromDisk("/path/to/file.jpg", "custom_name.jpg")  // with override name

// Use cloud file (Telegram file_id or URL)
tg.FromCloud("BQACAgIA...")
tg.FromCloud("https://example.com/file.pdf")

// Download files
tg.GenericDownload(ctx, "./save/path.jpg", fileId)
path, err := tg.GenericDownloadTemp(ctx, fileId)             // saves to temp file
path, err := tg.GenericDownloadTemp(ctx, fileId, "/tmp", "prefix_*.jpg")

// Photo wrapper (picks largest size automatically)
photos := tg.TelegramPhoto(upd.Message.Photo)
photos.FileId()
photos.Download(ctx, "./photo.jpg")
tmpPath, err := photos.DownloadTemp(ctx)
```

---

## Album / media group handling

```go
bot.Branch(tg.OnPhoto, tg.HandleAlbum(func(ctx context.Context, updates []*tg.Update) error {
    // all updates with same MediaGroupId arrive together after ~550ms
    for _, upd := range updates {
        // process each photo
    }
    return nil
}))

// With config
tg.HandleAlbum(handler, &tg.ConfigHandleAlbum{HandlingTimeout: 200 * time.Millisecond})
```

---

## Scheduler / rate limiting

```go
// Attach to bot (uses Telegram-safe defaults)
bot.Scheduler(tg.NewScheduler())
bot.Scheduler(tg.NewScheduler(
    tg.SchedulerClauseGlobal(30, 1500*time.Millisecond),
    tg.SchedulerClauseChat(20, 60500*time.Millisecond),
    tg.SchedulerClauseUser(100, 30500*time.Millisecond),
))

// Use inside handlers
func myHandler(ctx context.Context, upd *tg.Update) error {
    tg.ContextSchedule(ctx, upd.Message.Chat.Id, 1)
    defer tg.ContextScheduleDone(ctx, upd.Message.Chat.Id, 1)
    _, err := tg.SendMessage(ctx, upd.Message.Chat.Id, "hello")
    return err
}
```

Default clauses (applied when `NewScheduler()` called with no args):
- Global: 30 req / 1500ms
- Chat: 20 req / 60500ms, 10 req / 10500ms
- User: 100 req / 30500ms, 30 req / 5500ms

---

## Plugin / middleware

```go
// Built-in plugins
bot.Plugin(tg.PluginLogger(slog.LevelInfo))
bot.Plugin(tg.PluginLoggerFrom(myLogger))
bot.Plugin(tg.PluginOnError(myOnErrorFn))

// Error handling shortcuts
bot.OnError(tg.OnErrorLog)    // log and continue
bot.OnError(tg.OnErrorExit)   // log and exit(1)
bot.OnError(tg.OnErrorPanic)  // panic

// Custom plugin
type MyPlugin struct{}
func (p *MyPlugin) Hooks() []tg.PluginHookType {
    return []tg.PluginHookType{tg.PluginHookOnError, tg.PluginHookOnHandleFinish}
}
func (p *MyPlugin) Apply(ctx tg.PluginHookContext) {
    switch c := ctx.(type) {
    case *tg.PluginHookContextOnError:
        // c.Context, c.Bot, c.Error
    case *tg.PluginHookContextOnHandleFinish:
        // c.Context, c.Bot, c.Update, c.Handler, c.Error
    }
}
```

Hook types: `PluginHookOnUpdate`, `PluginHookOnFilter`, `PluginHookOnHandleStart`, `PluginHookOnHandleFinish`, `PluginHookOnError`

---

## Context

```go
// Canonical way to create a context (respects TimeoutHandle config)
ctx, cancel := bot.ContextWithCancel()
defer cancel()

// Add custom values
bot.ExtraContext(tg.WithCustomHttpClient(myClient))
bot.ExtraContext(func(ctx context.Context) context.Context {
    return context.WithValue(ctx, myKey, myValue)
})

// Context keys (rarely needed directly)
// tg.ContextBotInstance, ContextToken, ContextHttpClient, ContextApiUrl,
// ContextExtraHeaders, ContextFileDownloadType, ContextScheduler
```

---

## Error types

```go
// Check error kind
tg.IsApiError(err)           // is a Telegram API error
tg.IsTooManyRequests(err)    // is 429

// Error structs
var e *tg.Error              // .Code, .Description
var r *tg.ErrorTooManyRequests  // .Description, .RetryAfter
errors.As(err, &e)
errors.As(err, &r)
```

---

## Message helpers

```go
msg.TextOrCaption()                  // Text if set, else Caption, else ""
msg.TextOrCaptionEntities()          // iter.Seq[*MessageEntity]
tg.AsReplyTo(msg)                    // *ReplyParameters for replying
```

---

## Testing package (`tgtesting`)

```go
import "github.com/kittenbark/tg/tgtesting"

ctx := tgtesting.NewTestingContext(t, &tgtesting.Config{
    Stubs: []tgtesting.Stub{
        {
            Url:       "/bot{token}/sendMessage",
            Validator: func(req *http.Request) bool { return true },
            Result:    tgtesting.StubResultOK(200, &tg.Message{MessageId: 1}),
        },
    },
})
// Use ctx with any tg.* method; calls hit local stub server instead of Telegram
```

```go
// Assertions
var require = tgtesting.Require
require.NoError(t, err)
require.Error(t, err)
require.Equal(t, expected, actual)
require.True(t, condition)
require.NotNil(t, value)
```

---

## Scaffolding new projects

```sh
go run github.com/kittenbark/tg/init@latest
go run github.com/kittenbark/tg/init@latest -module github.com/myorg/mybot
go run github.com/kittenbark/tg/init@latest -overwrite
```

Generates: `cmd/main.go`, `Dockerfile`, `compose.yml`, `go.mod`, `.env`, `.gitignore`

---

## Key Update fields

```go
upd.Message          *Message
upd.EditedMessage    *Message
upd.CallbackQuery    *CallbackQuery
upd.InlineQuery      *InlineQuery
upd.ChannelPost      *Message

upd.Message.Chat.Id  int64
upd.Message.From     *User
upd.Message.Text     string
upd.Message.Caption  string
upd.Message.Photo    []*PhotoSize      // wrap with tg.TelegramPhoto(upd.Message.Photo)
upd.Message.MessageId int64
upd.Message.MediaGroupId string        // non-empty when part of album
```

---

## Code generation

Generated files (`api_types.go`, `api_methods.go`, `api_types_unmarshalers.go`) come from `tgcodegen/` — do not edit by hand. Schema at `tgcodegen/data/schema.json`. Custom overrides at `tgcodegen/data/custom.json`. Run `go generate` in the module root to regenerate.

`tgcodegen2/` is a work-in-progress rewrite of the codegen; its `gen/` output is experimental.
