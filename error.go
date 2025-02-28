package tg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func onErrorIgnore(ctx context.Context, err error) {}

func OnErrorLog(ctx context.Context, err error) {
	slog.ErrorContext(ctx, "tg.Bot#on_error", "err", err)
}

func OnErrorExit(ctx context.Context, err error) {
	OnErrorLog(ctx, err)
	os.Exit(1)
}

func OnErrorPanic(ctx context.Context, err error) { panic(err) }

type Error struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("telegram: %s (%d)", err.Description, err.Code)
}

type ErrorTooManyRequests struct {
	Description string        `json:"description"`
	RetryAfter  time.Duration `json:"retry_after"`
}

func (err *ErrorTooManyRequests) Error() string {
	return fmt.Sprintf("telegram: %s (%d)", err.Description, http.StatusTooManyRequests)
}

func IsTooManyRequests(err error) bool {
	var errErrorTooManyRequests *ErrorTooManyRequests
	return errors.As(err, &errErrorTooManyRequests)
}

func IsApiError(err error) bool {
	var errError *Error
	var errErrorTooManyRequests *ErrorTooManyRequests
	switch {
	case errors.As(err, &errError):
		return true
	case errors.As(err, &errErrorTooManyRequests):
		return true
	default:
		return false
	}
}
