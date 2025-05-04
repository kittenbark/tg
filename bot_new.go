package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvToken            = envPrefix + "TOKEN"
	EnvTokenTesting     = envPrefix + "TEST_TOKEN"
	EnvTestingChat      = envPrefix + "TEST_CHAT"
	EnvTestingGroupChat = envPrefix + "TEST_GROUP_CHAT"
	EnvApiURL           = envPrefix + "API_URL"
	EnvDownloadType     = envPrefix + "DOWNLOAD_TYPE"
	// EnvOnError is either ignore/log/exit.
	EnvOnError = envPrefix + "ON_ERROR"
	envPrefix  = "KITTENBARK_TG_"

	EnvSyncedHandle      = envPrefix + "SYNCED_HANDLE"
	EnvTimeoutHandle     = envPrefix + "TIMEOUT_HANDLE"
	defaultHandleTimeout = time.Hour

	EnvTimeoutPolling     = envPrefix + "TIMEOUT_POLL"
	defaultPollingTimeout = 100 * time.Millisecond
)

type DownloadType int

const (
	DownloadTypeUnspecified DownloadType = iota // calls default strategy (classic)
	DownloadTypeClassic                         // calls fileDownloadClassic
	DownloadTypeLocalMove                       // calls fileDownloadLocalMove
	DownloadTypeLocalCopy                       // calls fileDownloadLocalCopy
)

type Config struct {
	Token         string        `json:"token"`
	TokenTesting  string        `json:"token_testing"`
	ApiURL        string        `json:"api_url,omitempty"`
	TimeoutHandle time.Duration `json:"timeout,omitempty"`
	TimeoutPoll   time.Duration `json:"timeout_poll,omitempty"`
	SyncHandling  bool          `json:"sync,omitempty"`
	DownloadType  DownloadType  `json:"download_type,omitempty"`
	OnError       OnErrorFunc   `json:"-"`
	OnErrorByType string        `json:"on_error,omitempty"`

	buildType int
}

func New(cfg *Config) *Bot {
	bot, err := TryNew(cfg)
	if err != nil {
		panic(err)
	}
	return bot
}

func NewFromEnv() *Bot {
	bot, err := TryNewFromEnv()
	if err != nil {
		panic(err)
	}
	return bot
}

func NewFromFile(path string) *Bot {
	bot, err := TryNewFromFile(path)
	if err != nil {
		panic(err)
	}
	return bot
}

func TryNew(cfg *Config) (*Bot, error) {
	ctx := context.Background()

	if cfg.Token != "" {
		ctx = context.WithValue(ctx, ContextToken, cfg.Token)
	} else if cfg.TokenTesting != "" {
		slog.Warn("from_env: api token not found, using test token")
		ctx = context.WithValue(ctx, ContextTestToken, cfg.TokenTesting)
	} else {
		return nil, buildError(cfg.buildType,
			fmt.Errorf("config: missing bot api token or token_testing"),
			fmt.Errorf("env: missing bot api token (at '%s' (or for testing '%s'))", EnvToken, EnvTokenTesting),
		)
	}

	if cfg.ApiURL != "" {
		ctx = context.WithValue(ctx, ContextApiUrl, cfg.ApiURL)
	}

	switch cfg.DownloadType {
	case DownloadTypeUnspecified:
		if cfg.ApiURL == "" {
			ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadClassic)
			break
		}
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadClassic)
		link, err := url.Parse(cfg.ApiURL)
		if err != nil {
			return nil, fmt.Errorf("env: error '%s' while parsing '%s'", err.Error(), EnvApiURL)
		}
		switch strings.ToLower(link.Hostname()) {
		case "localhost", "127.0.0.1":
			ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadLocalMove)
		}
	case DownloadTypeClassic:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadClassic)
	case DownloadTypeLocalMove:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadLocalMove)
	case DownloadTypeLocalCopy:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadLocalCopy)
	default:
		return nil, fmt.Errorf("config: invalid download type: %#v", cfg.DownloadType)
	}

	onError, err := buildPluginsOnError(cfg)
	if err != nil {
		return nil, err
	}

	return &Bot{
		context:           ctx,
		contextCancelFunc: func() {},
		contextTimeout:    withDefault(cfg.TimeoutHandle, defaultHandleTimeout, 0),
		plugins: map[PluginHookType][]Plugin{
			PluginHookOnUpdate:       {},
			PluginHookOnFilter:       {},
			PluginHookOnHandleStart:  {},
			PluginHookOnHandleFinish: {},
			PluginHookOnError:        onError,
		},
		defaultHandler: nil,
		syncHandling:   cfg.SyncHandling,
		pollTimeout:    withDefault(cfg.TimeoutPoll, defaultPollingTimeout, 0),
		updatesOffset:  0,
	}, nil
}

func buildPluginsOnError(cfg *Config) ([]Plugin, error) {
	onError := []Plugin{}
	if cfg.OnError != nil {
		onError = append(onError, PluginOnError(cfg.OnError))
	} else if cfg.OnErrorByType != "" {
		switch strings.TrimSpace(cfg.OnErrorByType) {
		case "ignore":
			onError = append(onError, PluginOnError(onErrorIgnore))
		case "log":
			onError = append(onError, PluginOnError(OnErrorLog))
		case "exit":
			onError = append(onError, PluginOnError(OnErrorLog))
		default:
			return nil, buildError(cfg.buildType,
				fmt.Errorf("config: unknown on_error value '%s'", cfg.OnErrorByType),
				fmt.Errorf("env: unknown onError value '%s' (at '%s')", cfg.OnErrorByType, EnvOnError),
			)
		}
	}
	return onError, nil
}

func TryNewFromEnv() (*Bot, error) {
	config, err := configFromEnv()
	if err != nil {
		return nil, err
	}
	return TryNew(config)
}

func TryNewFromFile(path string) (*Bot, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(contents, &config); err != nil {
		return nil, err
	}

	return TryNew(&config)
}

const (
	buildTypeConfig = 0
	buildTypeEnv    = 1
)

func buildError[T any](buildType int, config T, env T) T {
	switch buildType {
	case buildTypeConfig:
		return config
	case buildTypeEnv:
		return env
	default:
		panic("unknown build type")
	}
}

func configFromEnv() (config *Config, err error) {
	var downloadType DownloadType
	if env, ok := os.LookupEnv(EnvDownloadType); ok {
		switch strings.ToLower(strings.TrimSpace(env)) {
		case "classic":
			downloadType = DownloadTypeClassic
		case "local_move":
			downloadType = DownloadTypeLocalMove
		case "local_copy":
			downloadType = DownloadTypeLocalCopy
		default:
			return nil, fmt.Errorf("env: unknown '%s' (at %s)", env, EnvDownloadType)
		}
	}

	config = &Config{
		Token:         os.Getenv(EnvToken),
		TokenTesting:  os.Getenv(EnvTokenTesting),
		ApiURL:        os.Getenv(EnvApiURL),
		DownloadType:  downloadType,
		OnError:       nil,
		OnErrorByType: strings.ToLower(os.Getenv(EnvOnError)),
		buildType:     buildTypeEnv,
	}
	if config.SyncHandling, err = parseFromEnvBool(EnvSyncedHandle, false); err != nil {
		return nil, err
	}
	if config.TimeoutHandle, err = parseFromEnvDuration(EnvTimeoutHandle, -1); err != nil {
		return nil, err
	}
	if config.TimeoutPoll, err = parseFromEnvDuration(EnvTimeoutPolling, -1); err != nil {
		return nil, err
	}

	return config, nil
}

func parseFromEnvDurationMust(env string, otherwise time.Duration) time.Duration {
	result, err := parseFromEnvDuration(env, otherwise)
	if err != nil {
		panic(err)
	}
	return result
}

func parseFromEnvDuration(env string, otherwise time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(EnvSyncedHandle)
	if !ok {
		return otherwise, nil
	}

	seconds, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return otherwise, fmt.Errorf("env: invalid '%s' (at %s), err '%s'",
			value, env, err.Error(),
		)
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

func parseFromEnvBool(env string, otherwise bool) (bool, error) {
	value, ok := os.LookupEnv(env)
	if !ok {
		return otherwise, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return otherwise, fmt.Errorf("env: invalid '%s' (at %s), err '%s'",
			value, env, err.Error(),
		)
	}
	return result, nil
}

func withDefault[T ~float64 | ~int64](value T, onZero T, onNegative T) T {
	switch {
	case value < 0:
		return onNegative
	case value == 0:
		return onZero
	default:
		return value
	}
}
