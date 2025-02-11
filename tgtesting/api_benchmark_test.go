/*
>> NOTE: the only important thing in these benchmarks is the CPU profiler.
These tests are run against a local http server, which produces too much overhead.
Though you could check the CPU profiler results, to see who your generic changes are affecting performance.
*/
package tgtesting

import (
	"github.com/kittenbark/tg"
	"net/http"
	"testing"
)

/*
goos: darwin
goarch: arm64
pkg: tg/tgtesting
cpu: Apple M2
BenchmarkSendMessage
BenchmarkSendMessage-8   	   23455	     46669 ns/op
PASS
*/
func BenchmarkSendMessage(b *testing.B) {
	ctx := NewTestingEnvLessStrict(&Config{
		Stubs: []Stub{{Url: "/sendMessage", Result: func(req *http.Request) (status int, body *Response) {
			return http.StatusOK, &Response{
				Ok: true,
				Result: &tg.Message{
					MessageId: 123,
					Chat:      &tg.Chat{Id: 123456},
					Text:      "hello",
				},
			}
		}}},
	})
	for i := 0; i < b.N; i++ {
		_, err := tg.SendMessage(ctx, int64(12345), "hello", &tg.OptSendMessage{MessageThreadId: 16})
		if err != nil {
			b.Fatal(err)
		}
	}
}

/*
goos: darwin
goarch: arm64
pkg: tg/tgtesting
cpu: Apple M2
BenchmarkSendPhotoMultipart
BenchmarkSendPhotoMultipart-8   	   16150	     68793 ns/op
PASS
*/
func BenchmarkSendPhotoMultipart(b *testing.B) {
	ctx := NewTestingEnvLessStrict(&Config{
		Stubs: []Stub{{Url: "/sendPhoto", Result: func(req *http.Request) (status int, body *Response) {
			return http.StatusOK, &Response{
				Ok: true,
				Result: &tg.Message{
					MessageId: 123,
					Chat:      &tg.Chat{Id: 123456},
					Photo: []*tg.PhotoSize{{
						FileId:   "abc123",
						Width:    600,
						Height:   800,
						FileSize: 1024,
					}}},
			}
		}}},
	})

	for i := 0; i < b.N; i++ {
		_, err := tg.SendPhoto(ctx, int64(123), tg.FromCloud("hello"))
		if err != nil {
			b.Fatal(err)
		}
	}
}
