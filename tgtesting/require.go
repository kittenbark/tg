package tgtesting

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"reflect"
	"testing"
)

// require he-he.
var require = &bicycle{}

type bicycle struct{}

func (cycle *bicycle) NoError(t *testing.T, err error, msg ...string) {
	if err != nil {
		t.Fatalf("require.NoError failed: err '%v' %#v", err, msg)
	}
}

func (cycle *bicycle) Error(t *testing.T, err error, msg ...string) {
	if err == nil {
		t.Fatalf("require.Error failed: no error %#v", msg)
	}
}

func (cycle *bicycle) Equal(t *testing.T, expected any, actual any, msg ...string) {
	if !reflect.DeepEqual(actual, expected) {
		expectedJSON, _ := json.Marshal(expected)
		actualJSON, _ := json.Marshal(actual)
		t.Fatalf("require.Equal failed: expected '%v' // actual '%v' %v", string(expectedJSON), string(actualJSON), msg)
	}
}

func (cycle *bicycle) True(t *testing.T, actual any, msg ...string) {
	cycle.Equal(t, true, actual, msg...)
}

func (cycle *bicycle) False(t *testing.T, actual any, msg ...string) {
	cycle.Equal(t, false, actual, msg...)
}

func (cycle *bicycle) NotNil(t *testing.T, actual any, msg ...string) {
	if actual == nil || reflect.ValueOf(actual).IsNil() {
		t.Fatalf("require.NotNil failed: expected nil, got '%#v' %#v", actual, msg)
	}
}

func OutsideFile(local string, url string) string {
	if _, err := os.Stat(path.Dir(local)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(local), 0755); err != nil {
			return ""
		}
	}
	if _, err := os.Stat(local); err == nil {
		return local
	}

	slog.Info("downloading...", "path", local, "url", url)
	result, err := os.Create(local)
	if err != nil {
		return ""
	}
	defer result.Close()

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if _, err := io.Copy(result, resp.Body); err != nil {
		return ""
	}

	return local
}
