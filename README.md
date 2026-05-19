# kittenbark/tg

> Go is an open source programming language that makes it easy to build simple, reliable, and
> efficient software. _(From Golang README.md)_

This package aims to be exactly that — trivial and straightforward.

- Consistent, simple API.
- Declarative updates flow, filtering, and branches.
- Compile-time safety: zero `interface{}` exposed to user-space.
- Zero dependencies.

Inspired by [teloxide](https://github.com/teloxide/teloxide) and its declarative nature.

## Quickstart

Scaffold a new project:

```sh
go run github.com/kittenbark/tg/init@latest
```

Or set `KITTENBARK_TG_TOKEN=<your bot token>` and write it yourself:

```go
package main

import (
    "context"
    "github.com/kittenbark/tg"
)

func main() {
    tg.NewFromEnv().
        Filter(tg.OnPrivateMessage).
        Command("/start", tg.CommonTextReply("hello, this echo bot is made with @kittenbark_tg")).
        Command("/help", tg.CommonTextReply("just send a message")).
        Branch(tg.OnMessage, func(ctx context.Context, upd *tg.Update) error {
            msg := upd.Message
            _, err := tg.CopyMessage(ctx, msg.Chat.Id, msg.Chat.Id, msg.MessageId)
            return err
        }).
        Start()
}
```

## Pipeline

The bot is built as a pipeline of filters, branches, and handlers. Each update flows through it top to bottom.

- **`Filter`** drops updates that don't match — everything below only sees updates that pass.
- **`Branch`** handles matching updates and lets the rest fall through to the next branch.
- **`Command`** is shorthand for `Branch(tg.OnCommand("/cmd"), handler)`.
- **`Handle`** is a catch-all at the end.

```go
tg.NewFromEnv().
    Filter(tg.OnPrivateMessage).       // only private messages from here on
    Command("/start", ...).            // handled, others fall through
    Branch(tg.OnPhoto, ...).           // handled, others fall through
    Handle(tg.CommonTextReply("?")).   // everything else
    Start()
```

## Filters

Filters are plain `func(ctx, upd) bool` values and compose freely:

```go
tg.All(tg.OnText, tg.OnPublicMessage)        // both must match
tg.Either(tg.OnPhoto, tg.OnVideo)            // either matches
tg.Not(tg.OnForwarded)                       // negation
tg.OnTextRegexp(`^\d+$`)                     // regex
tg.OnCommand("/ban")                         // specific command
tg.OnChat(chatId1, chatId2)                  // whitelist chats
tg.OnSender(userId)                          // whitelist senders
tg.OnCallbackWithData[MyStruct](predicate)   // typed callback data
```

## Handlers

Handlers are plain `func(ctx, upd) error` values. Because they're just functions, they can be defined anywhere and composed:

```go
// reusable component, no bot reference needed
func RequireAdmin(ctx context.Context, upd *tg.Update) error {
    // check admin status, return error to stop the chain
}

func ProcessMedia(ctx context.Context, upd *tg.Update) error {
    // heavy processing
}

bot.Branch(tg.OnVideo, tg.Chain(
    RequireAdmin,
    tg.CommonReactionReply(":eyes:"),
    tg.Synced(ProcessMedia),       // only one at a time
    tg.CommonReactionReply(":ok:"),
))
```

`tg.Chain` runs handlers in sequence and stops on the first error. `tg.Synced` wraps any handler in a mutex. For a shared mutex across multiple handlers use `tg.SyncedGroup`:

```go
group := &tg.SyncedGroup{}

bot.
    Branch(tg.OnVideo, group.Synced(HandleVideo)).
    Branch(tg.OnAudio, group.Synced(HandleAudio))
```

## Common helpers

Ready-made handlers for typical patterns:

```go
tg.CommonTextReply("hello!")                          // send a text reply
tg.CommonTextReply("hello!", true)                    // as a reply to the message
tg.CommonTextReplyExpiring(5*time.Second, "wait...")  // send then delete after duration
tg.CommonReactionReply(":ok:")                        // react to the message
tg.CommonDeleteMessage                                // delete the message
tg.CommonRestrictSender(permissions)                  // restrict the sender
tg.CommonArgs[Args](func(ctx, upd, args *Args) error) // parse command arguments
```

`tg.CommonReaction` understands named aliases like `":ok:"`, `":fire:"`, `":eyes:"` — see `tg.CommonReactionEmojiMap` for the full list.

## More examples and documentation

- [https://kittenbark.com/tg](https://kittenbark.com/tg) — more info and docs
- [tgdeploy](https://github.com/kittenbark/tgdeploy) — Dockerfile/compose.yml templates for deploying your bots with ease

## Contributing

Make a PR or create an issue. Submit your examples to be listed above. Enjoy.

## Todo

- [ ] update to 9.0 API version (problem: no InputFile for thumbnail in InputMediaVideo)