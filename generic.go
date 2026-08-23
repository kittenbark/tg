package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"time"
)

const DefaultTelegramApiUrl = "https://api.telegram.org"

func GenericRequest[Request any, Result any](ctx context.Context, method string, request *Request) (result Result, err error) {
	token, err := tryGetTokenFromContext(ctx)
	if err != nil {
		return
	}

	prepared := defaults(request)
	url := fmt.Sprintf("%s/bot%s/%s", getOrDefault(ctx, ContextApiUrl, DefaultTelegramApiUrl), token, method)

	for {
		var body bytes.Buffer
		if err = json.NewEncoder(&body).Encode(prepared); err != nil {
			return
		}

		httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
		if err != nil {
			return result, err
		}
		httpRequest.Header.Set("Content-Type", "application/json")
		if headers := getOrDefault(ctx, ContextExtraHeaders, map[string]string{}); headers != nil {
			for name, value := range headers {
				httpRequest.Header.Add(name, value)
			}
		}

		httpResponse, err := getOrDefault(ctx, ContextHttpClient, http.DefaultClient).Do(httpRequest)
		if err != nil {
			return result, err
		}

		type HttpResult struct {
			Ok          bool                   `json:"ok"`
			ErrorCode   int                    `json:"error_code,omitempty"`
			Description string                 `json:"description,omitempty"`
			Parameters  map[string]interface{} `json:"parameters,omitempty"`
			Result      Result                 `json:"result,omitempty"`
		}
		var httpResult HttpResult
		decodeErr := json.NewDecoder(httpResponse.Body).Decode(&httpResult)
		_ = httpResponse.Body.Close()
		if decodeErr != nil {
			return result, decodeErr
		}

		if !httpResult.Ok {
			apiErr := newTelegramError(httpResult.ErrorCode, httpResult.Description, httpResult.Parameters)
			var tooMany *ErrorTooManyRequests
			if errors.As(apiErr, &tooMany) {
				ContextScheduleThrottle(ctx, 0, tooMany.RetryAfter)
				select {
				case <-ctx.Done():
					return result, ctx.Err()
				case <-time.After(tooMany.RetryAfter):
					continue
				}
			}
			return result, apiErr
		}

		return httpResult.Result, nil
	}
}

func newTelegramError(code int, description string, parameters map[string]interface{}) error {
	switch code {
	case http.StatusTooManyRequests:
		var retryAfter time.Duration
		if retryAfterOpt, ok := parameters["retry_after"]; ok {
			retryAfterFloat, _ := retryAfterOpt.(float64)
			retryAfter = time.Duration(retryAfterFloat * float64(time.Second))
		}
		return &ErrorTooManyRequests{Description: description, RetryAfter: retryAfter}
	default:
		return &Error{Code: code, Description: description}
	}
}

func getOrDefault[T any](ctx context.Context, key string, defaultValue T) T {
	if ctx == nil {
		return defaultValue
	}
	result, ok := ctx.Value(key).(T)
	if !ok {
		return defaultValue
	}
	return result
}

func containsAll(container []string, items []string) bool {
	for _, item := range items {
		if !slices.Contains(container, item) {
			return false
		}
	}
	return true
}

func containsAny(container []string, items []string) bool {
	for _, item := range items {
		if slices.Contains(container, item) {
			return true
		}
	}
	return false
}

func deref[T any](ptr *T) (val T) {
	if ptr == nil {
		return
	}
	return *ptr
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	default:
		return false
	}
}

func getFuncName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func at[T any](list []T, pos int, defaultValue T) T {
	if len(list) <= pos {
		return defaultValue
	}
	return list[pos]
}
