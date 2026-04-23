package tg

import (
	"context"
	"net/http"
)

const (
	ContextBotInstance      = contextPrefix + "bot_instance"
	ContextToken            = contextPrefix + "token"
	ContextTestToken        = contextPrefix + "test_token"
	ContextHttpClient       = contextPrefix + "http_client"
	ContextApiUrl           = contextPrefix + "api_url"
	ContextFileDownloadType = contextPrefix + "file_downloader"
	ContextScheduler        = contextPrefix + "scheduler"

	contextPrefix = "kittenbark_"
)

type ExtraContext func(context.Context) context.Context

func WithCustomHttpClient(client *http.Client) ExtraContext {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ContextHttpClient, client)
	}
}

func tryGetTokenFromContext(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(ContextToken).(string); ok {
		return token, nil
	}
	if token, ok := ctx.Value(ContextTestToken).(string); ok {
		return token, nil
	}
	return "", &Error{Description: "api token not found in context"}
}
