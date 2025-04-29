package tgtesting

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kittenbark/tg"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const ContextTestingMuxServer = "kittenbark_testing_mux_server"

var port = &atomic.Int32{}

func NewTestingContext(t *testing.T, cfg *Config) context.Context {
	cfg = cfg.WithDefaults()

	mux := http.NewServeMux()
	for _, stub := range cfg.Stubs {
		stub.RegisterTesting(t, cfg, mux)
	}
	go func() {
		server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: mux}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Millisecond * 10)

	ctx := context.WithValue(cfg.Context, tg.ContextToken, cfg.Token)
	ctx = context.WithValue(ctx, tg.ContextApiUrl, cfg.UrlWithPort())
	ctx = context.WithValue(ctx, ContextTestingMuxServer, mux)

	return ctx
}

func SetTestingEnv(t TestingEnv, cfg *Config) {
	cfg = cfg.WithDefaults()

	mux := http.NewServeMux()
	for _, stub := range cfg.Stubs {
		stub.RegisterTesting(t, cfg, mux)
	}
	go func() {
		server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: mux}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Millisecond * 10)

	t.Setenv(tg.EnvToken, cfg.Token)
	t.Setenv(tg.EnvApiURL, cfg.UrlWithPort())
}

func NewTestingEnvLessStrict(cfg *Config) context.Context {
	cfg = cfg.WithDefaults()

	mux := http.NewServeMux()
	for _, stub := range cfg.Stubs {
		stub.Register(cfg, mux)
	}
	go func() {
		server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: mux}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Millisecond * 10)

	ctx := context.WithValue(cfg.Context, tg.ContextToken, cfg.Token)
	ctx = context.WithValue(ctx, tg.ContextApiUrl, cfg.UrlWithPort())
	ctx = context.WithValue(ctx, ContextTestingMuxServer, mux)
	return ctx
}

func StubResultOK(code int, result any) func(req *http.Request) (int, *Response) {
	return func(req *http.Request) (int, *Response) {
		return code, &Response{
			Ok:     true,
			Result: result,
		}
	}
}

func StubResultError(code int, description string, params ...any) func(req *http.Request) (int, *Response) {
	return func(req *http.Request) (int, *Response) {
		var respParams map[string]any
		for i := range len(params) / 2 {
			if respParams == nil {
				respParams = map[string]any{}
			}
			respParams[params[i*2].(string)] = params[i*2+1]
		}

		return code, &Response{
			Ok:          false,
			ErrorCode:   code,
			Description: description,
			Parameters:  respParams,
		}
	}
}

type Stub struct {
	Url       string
	Validator func(req *http.Request) bool
	Result    func(req *http.Request) (status int, body *Response)
}

type TestingEnv interface {
	Fatalf(format string, args ...any)
	Setenv(key, value string)
}

func (stub *Stub) RegisterTesting(t TestingEnv, cfg *Config, mux *http.ServeMux) {
	url := fmt.Sprintf("/bot%s/%s", cfg.Token, strings.TrimLeft(stub.Url, "/"))
	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		if stub.Validator != nil && !stub.Validator(req) {
			t.Fatalf("validation of %v failed", req)
		}
		if stub.Result != nil {
			status, body := stub.Result(req)
			w.WriteHeader(status)
			if body != nil {
				data, err := json.Marshal(body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					if _, err := w.Write([]byte(err.Error())); err != nil {
						panic(err)
					}
				}
				if _, err := w.Write(data); err != nil {
					panic(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}

func (stub *Stub) Register(cfg *Config, mux *http.ServeMux) {
	url := fmt.Sprintf("/bot%s/%s", cfg.Token, strings.TrimLeft(stub.Url, "/"))
	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		if stub.Validator != nil && !stub.Validator(req) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if stub.Result != nil {
			status, body := stub.Result(req)
			w.WriteHeader(status)
			if body != nil {
				data, err := json.Marshal(body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					if _, err := w.Write([]byte(err.Error())); err != nil {
						panic(err)
					}
				}
				if _, err := w.Write(data); err != nil {
					panic(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}

type Config struct {
	Context  context.Context
	Stubs    []Stub
	Token    string
	Host     string
	Port     int
	UseHTTPS bool
}

func (cfg *Config) WithDefaults() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Stubs == nil {
		cfg.Stubs = []Stub{}
	}
	if cfg.Token == "" {
		cfg.Token = "123456:ABCDEFGHIJKLMN"
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == 0 {
		cfg.Port = 8080 + int(port.Add(1))
	}
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}
	return cfg
}

func (cfg *Config) UrlWithPort() string {
	protocol := "http://"
	if cfg.UseHTTPS {
		protocol = "https://"
	}
	return fmt.Sprintf("%s%s:%d", protocol, cfg.Host, cfg.Port)
}

type Response struct {
	Ok          bool                   `json:"ok"`
	ErrorCode   int                    `json:"error_code,omitempty"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Result      any                    `json:"result,omitempty"`
}
