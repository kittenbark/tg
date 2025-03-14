# kittenbark/tg

> Go is an open source programming language that makes it easy to build simple, reliable, and
> efficient software. _(From Golang README.md)_

This package should be as trivial and straightforward as Golang is to me.

- Consistent simple api.
- Trivial declarative updates flow, filtering and branches.
- Compile-time safety: zero `interface{}` exposed to user-space.
- Zero dependencies.

Inspired by [teloxide](https://github.com/teloxide/teloxide) and its declarative nature.

## Example: echo bot with commands and filtering

```go
package main

import (
    "context"
    "github.com/kittenbark/tg"
)

func main() {
    tg.NewFromEnv().
        Filter(tg.OnPrivateMessage).
        Command("/start", tg.CommonTextReply("hello this echo bot is made with @kittenbark_tg")).
        Command("/help", tg.CommonTextReply("just send a message")).
        Branch(tg.OnMessage, func(ctx context.Context, upd *tg.Update) error {
            msg := upd.Message
            _, err := tg.CopyMessage(ctx, msg.Chat.Id, msg.Chat.Id, msg.MessageId)
            return err
        }).
        Start()
}
```

## More complex examples, documentation and etc

- [https://kittenbark.com/tg/docs](https://kittenbark.com/tg/docs) — more info and docs,
- [tgdeploy](https://github.com/kittenbark/tgdeploy) — repository with Dockerfile/compose.yml templates for deploying
  your bots with ease,

## Contributing

Well, make a PR or create an issue. Btw, submit your examples to be listed above. Enjoy.