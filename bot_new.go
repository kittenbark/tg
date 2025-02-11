package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvToken        = envPrefix + "TG_TOKEN"
	EnvTokenTesting = envPrefix + "TG_TEST_TOKEN"
	EnvTestingChat  = envPrefix + "TG_TEST_CHAT"
	EnvSyncedHandle = envPrefix + "SYNCED_HANDLE"
	EnvApiURL       = envPrefix + "API_URL"
	EnvDownloadType = envPrefix + "DOWNLOAD_TYPE"
	// EnvOnError is either ignore/log/exit.
	EnvOnError = envPrefix + "ON_ERROR"
	envPrefix  = "KITTENBARK_"

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
	DownloadTypeLocalCopy                       // calls fileDownloadLocalMove
)

type Config struct {
	Token         string  `json:"token"`
	TokenTesting  string  `json:"token_testing"`
	ApiURL        string  `json:"api_url,omitempty"`
	HandleTimeout float64 `json:"timeout,omitempty"`
	PollTimeout   float64 `json:"poll_timeout,omitempty"`
	SyncHandling  bool    `json:"sync,omitempty"`
	DownloadType  DownloadType
	OnError       OnErrorFunc `json:"-"`
	OnErrorByType string      `json:"on_error,omitempty"`

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

	var pollTimeout time.Duration
	if cfg.PollTimeout < 0 {
		pollTimeout = 0
	} else if cfg.PollTimeout == 0 {
		pollTimeout = defaultPollingTimeout
	} else {
		pollTimeout = time.Duration(float64(time.Second) * cfg.PollTimeout)
	}

	var responseTimeout time.Duration
	if cfg.HandleTimeout < 0 {
		responseTimeout = 0
	} else if cfg.HandleTimeout == 0 {
		responseTimeout = defaultHandleTimeout
	} else {
		responseTimeout = time.Duration(float64(time.Second) * cfg.HandleTimeout)
	}

	switch cfg.DownloadType {
	case DownloadTypeUnspecified:
	case DownloadTypeClassic:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadClassic)
	case DownloadTypeLocalMove:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadLocalMove)
	case DownloadTypeLocalCopy:
		ctx = context.WithValue(ctx, ContextFileDownloadType, fileDownloadLocalCopy)
	default:
		return nil, fmt.Errorf("config: invalid download type: %#v", cfg.DownloadType)
	}

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

	ctx, ctxCancel := context.WithCancel(ctx)
	return &Bot{
		context:        ctx,
		contextTimeout: responseTimeout,
		contextCancel:  ctxCancel,
		plugins: map[PluginHookType][]Plugin{
			PluginHookOnUpdate: {},
			PluginHookOnFilter: {},
			PluginHookOnHandle: {},
			PluginHookOnError:  onError,
		},
		pipeline:       &pipe{},
		defaultHandler: nil,
		syncHandling:   cfg.SyncHandling,
		pollTimeout:    pollTimeout,
		updatesOffset:  0,
	}, nil
}

func TryNewFromEnv() (*Bot, error) {
	var timeoutHandle float64
	if env, ok := os.LookupEnv(EnvTimeoutHandle); ok {
		timeout, err := strconv.ParseFloat(env, 64)
		if err != nil {
			return nil, fmt.Errorf("env: invalid '%s' (at %s), err '%s'",
				env, EnvTimeoutHandle, err.Error(),
			)
		}
		timeoutHandle = timeout
	}

	var timeoutPoll float64
	if env, ok := os.LookupEnv(EnvTimeoutPolling); ok {
		timeout, err := strconv.ParseFloat(env, 64)
		if err != nil {
			return nil, fmt.Errorf("env: invalid '%s' (at %s), err '%s'",
				env, EnvTimeoutPolling, err.Error(),
			)
		}
		timeoutPoll = timeout
	}

	var syncHandling bool
	if env, ok := os.LookupEnv(EnvSyncedHandle); ok {
		sync, err := strconv.ParseBool(env)
		if err != nil {
			return nil, fmt.Errorf("env: invalid '%s' (at %s), err '%s'",
				env, EnvSyncedHandle, err.Error())
		}
		syncHandling = sync
	}

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

	config := &Config{
		Token:         os.Getenv(EnvToken),
		TokenTesting:  os.Getenv(EnvTokenTesting),
		ApiURL:        os.Getenv(EnvApiURL),
		HandleTimeout: timeoutHandle,
		PollTimeout:   timeoutPoll,
		SyncHandling:  syncHandling,
		DownloadType:  downloadType,
		OnError:       nil,
		OnErrorByType: strings.ToLower(os.Getenv(EnvOnError)),
		buildType:     buildTypeEnv,
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
