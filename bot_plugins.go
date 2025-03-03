package tg

import (
	"context"
	"log/slog"
	"os"
)

type PluginHookType int

const (
	PluginHookOnUpdate PluginHookType = iota
	PluginHookOnFilter
	PluginHookOnHandleStart
	PluginHookOnHandleFinish
	PluginHookOnError
)

type (
	PluginHookContext = interface {
		OnHook() PluginHookType
	}

	PluginHookContextOnUpdate struct {
		Context context.Context
		Bot     *Bot
		Update  *Update
	}
	PluginHookContextOnFilter struct {
		Context context.Context
		Bot     *Bot
		Filter  FilterFunc
	}
	PluginHookContextOnHandleStart struct {
		Context context.Context
		Bot     *Bot
		Update  *Update
		Handler HandlerFunc
	}
	PluginHookContextOnHandleFinish struct {
		Context context.Context
		Bot     *Bot
		Update  *Update
		Handler HandlerFunc
	}
	PluginHookContextOnError struct {
		Context context.Context
		Bot     *Bot
		Error   error
	}
)

func (p *PluginHookContextOnUpdate) OnHook() PluginHookType       { return PluginHookOnUpdate }
func (p *PluginHookContextOnFilter) OnHook() PluginHookType       { return PluginHookOnFilter }
func (p *PluginHookContextOnHandleStart) OnHook() PluginHookType  { return PluginHookOnHandleStart }
func (p *PluginHookContextOnHandleFinish) OnHook() PluginHookType { return PluginHookOnHandleFinish }
func (p *PluginHookContextOnError) OnHook() PluginHookType        { return PluginHookOnError }

var (
	_ PluginHookContext = (*PluginHookContextOnUpdate)(nil)
	_ PluginHookContext = (*PluginHookContextOnFilter)(nil)
	_ PluginHookContext = (*PluginHookContextOnHandleStart)(nil)
	_ PluginHookContext = (*PluginHookContextOnHandleFinish)(nil)
	_ PluginHookContext = (*PluginHookContextOnError)(nil)
)

// Plugin allows minor modifications in Bot flow, i.e. logging requests, handling errors or even orchestrating bulk requests.
type Plugin interface {
	Hooks() []PluginHookType
	Apply(ctx PluginHookContext)
}

var (
	_ Plugin = (*pluginOnError)(nil)
	_ Plugin = (*pluginLogger)(nil)
)

type pluginOnError OnErrorFunc

func (plugin pluginOnError) Hooks() []PluginHookType {
	return []PluginHookType{PluginHookOnError}
}

func (plugin pluginOnError) Apply(ctx PluginHookContext) {
	switch errorContext := ctx.(type) {
	case *PluginHookContextOnError:
		plugin(errorContext.Context, errorContext.Error)
	}
}

func PluginOnError(fn OnErrorFunc) Plugin {
	return pluginOnError(fn)
}

type pluginLogger struct {
	logger *slog.Logger
}

func (plugin *pluginLogger) Hooks() []PluginHookType {
	return []PluginHookType{PluginHookOnUpdate, PluginHookOnFilter, PluginHookOnHandleStart, PluginHookOnHandleFinish}
}

func (plugin *pluginLogger) Apply(ctx PluginHookContext) {
	switch ctx := ctx.(type) {
	case *PluginHookContextOnUpdate:
		plugin.logger.InfoContext(ctx.Context, "bot#update", "update", ctx.Update)
	case *PluginHookContextOnFilter:
		plugin.logger.DebugContext(ctx.Context, "bot#filter_out", "func", getFuncName(ctx.Filter))
	case *PluginHookContextOnHandleStart:
		plugin.logger.DebugContext(ctx.Context, "bot#handle_start", "update_id", ctx.Update.UpdateId, "func", getFuncName(ctx.Handler))
	case *PluginHookContextOnHandleFinish:
		plugin.logger.DebugContext(ctx.Context, "bot#handle_finish", "update_id", ctx.Update.UpdateId, "func", getFuncName(ctx.Handler))
	case *PluginHookContextOnError:
		plugin.logger.ErrorContext(ctx.Context, "bot#error", "err", ctx.Error)
	}
}

func PluginLogger(level slog.Level) Plugin {
	return PluginLoggerFrom(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	})))
}

func PluginLoggerFrom(logger *slog.Logger) Plugin {
	return &pluginLogger{
		logger: logger,
	}
}
