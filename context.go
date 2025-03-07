package tg

import (
	"context"
)

const (
	ContextToken            = contextPrefix + "token"
	ContextTestToken        = contextPrefix + "test_token"
	ContextHttpClient       = contextPrefix + "http_client"
	ContextApiUrl           = contextPrefix + "api_url"
	ContextFileDownloadType = contextPrefix + "file_downloader"
	ContextScheduler        = contextPrefix + "scheduler"

	contextPrefix = "kittenbark_"
)

func tryGetTokenFromContext(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(ContextToken).(string); ok {
		return token, nil
	}
	if token, ok := ctx.Value(ContextTestToken).(string); ok {
		return token, nil
	}
	return "", &Error{Description: "api token not found in context"}
}
