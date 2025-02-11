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
	PluginHookOnHandle
	PluginHookOnError
)

type Plugin interface {
	Hooks() []PluginHookType
	Apply(hook PluginHookType, ctx context.Context)
}

var (
	_ Plugin = (*pluginOnError)(nil)
	_ Plugin = (*pluginLogger)(nil)
)

type pluginOnError OnErrorFunc

func (plugin pluginOnError) Hooks() []PluginHookType {
	return []PluginHookType{PluginHookOnError}
}

func (plugin pluginOnError) Apply(hook PluginHookType, ctx context.Context) {
	if hook != PluginHookOnError {
		return
	}

	// Note: this could panic, but only of context.Context is ill-formed.
	plugin(ctx, ctx.Value(ContextPluginHooksError).(error))
}

func PluginOnError(fn OnErrorFunc) Plugin {
	return pluginOnError(fn)
}

type pluginLogger struct {
	logger *slog.Logger
}

func (plugin *pluginLogger) Hooks() []PluginHookType {
	return []PluginHookType{PluginHookOnUpdate, PluginHookOnFilter, PluginHookOnHandle}
}

func (plugin *pluginLogger) Apply(hook PluginHookType, ctx context.Context) {
	switch hook {
	case PluginHookOnUpdate:
		if plugin.logger.Enabled(ctx, slog.LevelInfo) {
			plugin.logger.InfoContext(ctx, "bot#update", "update", ctx.Value(ContextPluginHooksUpdate))
		}
	case PluginHookOnFilter:
		if plugin.logger.Enabled(ctx, slog.LevelDebug) {
			plugin.logger.DebugContext(ctx, "bot#filter", "func", getFuncName(ctx.Value(ContextPluginHooksFilter)))
		}
	case PluginHookOnHandle:
		if plugin.logger.Enabled(ctx, slog.LevelDebug) {
			plugin.logger.DebugContext(ctx, "bot#handle", "func", getFuncName(ctx.Value(ContextPluginHooksHandle)))
		}
	case PluginHookOnError:
		if plugin.logger.Enabled(ctx, slog.LevelWarn) {
			plugin.logger.WarnContext(ctx, "bot#error", "error", ctx.Value(ContextPluginHooksError))
		}
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
