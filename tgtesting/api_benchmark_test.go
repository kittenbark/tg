/*
>> NOTE: the only important thing in these benchmarks is the CPU profiler.
These tests are run against a local http server, which produces too much overhead.
Though you could check the CPU profiler results, to see who your generic changes are affecting performance.
*/
package tgtesting

import (
	"context"
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

func BenchmarkUrlParse(b *testing.B) {
	urls :=
		[]string{
			"http://www.youtube.com",
			"http://www.facebook.com",
			"http://www.baidu.com",
			"http://www.yahoo.com",
			"http://www.amazon.com",
			"http://www.wikipedia.org",
			"http://www.qq.com",
			"http://www.google.co.in",
			"http://www.twitter.com",
			"http://www.live.com",
			"http://www.taobao.com",
			"http://www.bing.com",
			"http://www.instagram.com",
			"http://www.weibo.com",
			"http://www.sina.com.cn",
			"http://www.linkedin.com",
			"http://www.yahoo.co.jp",
			"http://www.msn.com",
			"http://www.vk.com",
			"http://www.google.de",
			"http://www.yandex.ru",
			"http://www.hao123.com",
			"http://www.google.co.uk",
			"http://www.reddit.com",
			"http://www.ebay.com",
			"http://www.google.fr",
			"http://www.t.co",
			"http://www.tmall.com",
			"http://www.google.com.br",
			"http://www.360.cn",
			"http://www.sohu.com",
			"http://www.amazon.co.jp",
			"http://www.pinterest.com",
			"http://www.netflix.com",
			"http://www.google.it",
			"http://www.google.ru",
			"http://www.microsoft.com",
			"http://www.google.es",
			"http://www.wordpress.com",
			"http://www.gmw.cn",
			"http://www.tumblr.com",
			"http://www.paypal.com",
			"http://www.blogspot.com",
			"http://www.imgur.com",
			"http://www.stackoverflow.com",
			"http://www.aliexpress.com",
			"http://www.naver.com",
			"http://www.ok.ru",
			"http://www.apple.com",
			"http://www.github.com",
			"http://www.chinadaily.com.cn",
			"http://www.imdb.com",
			"http://www.google.co.kr",
			"http://www.fc2.com",
			"http://www.jd.com",
			"http://www.blogger.com",
			"http://www.163.com",
			"http://www.google.ca",
			"http://www.whatsapp.com",
			"http://www.amazon.in",
			"http://www.office.com",
			"http://www.tianya.cn",
			"http://www.google.co.id",
			"http://www.youku.com",
			"http://www.rakuten.co.jp",
			"http://www.craigslist.org",
			"http://www.amazon.de",
			"http://www.nicovideo.jp",
			"http://www.google.pl",
			"http://www.soso.com",
			"http://www.bilibili.com",
			"http://www.dropbox.com",
			"http://www.xinhuanet.com",
			"http://www.outbrain.com",
			"http://www.pixnet.net",
			"http://www.alibaba.com",
			"http://www.alipay.com",
			"http://www.microsoftonline.com",
			"http://www.booking.com",
			"http://www.googleusercontent.com",
			"http://www.google.com.au",
			"http://www.popads.net",
			"http://www.cntv.cn",
			"http://www.zhihu.com",
			"http://www.amazon.co.uk",
			"http://www.diply.com",
			"http://www.coccoc.com",
			"http://www.cnn.com",
			"http://www.bbc.co.uk",
			"http://www.twitch.tv",
			"http://www.wikia.com",
			"http://www.google.co.th",
			"http://www.go.com",
			"http://www.google.com.ph",
			"http://www.doubleclick.net",
			"http://www.onet.pl",
			"http://www.googleadservices.com",
			"http://www.accuweather.com",
			"http://www.googleweblight.com",
			"http://www.answers.yahoo.com",
		}

	messages := []string{
		"Jesus Christ",
		"Dogs bark loudly at midnight",
		"The train arrived late",
		"She danced through the rain",
		"Mountains stand tall against the sky",
		"Coffee tastes better hot",
		"Books opened new worlds",
		"Flowers bloom in spring",
		"He laughed until tears came",
		"Birds flew south",
		"Running water never grows stale",
		"Time flies",
		"Stars twinkled brightly",
		"Laughter is the best medicine",
		"Silence speaks volumes",
		"Trees swayed gently",
		"Music soothes the soul",
		"Dreams became reality",
		"Waves crashed against rocks",
		"Life finds a way",
	}

	const pattern = `^(https?|ftp):\/\/` + // Protocol (http, https, ftp)
		`(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*` + // Subdomains
		`([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])` + // Domain name
		`(:\d+)?` + // Port (optional)
		`(\/[-a-zA-Z0-9\._~:/?#[\]@!$&'()*+,;=%]*)??` // Path (optional)
	b.Run("regex", func(b *testing.B) {
		handler := tg.OnTextRegexp(pattern)
		ctx := context.Background()
		for b.Loop() {
			for _, url := range urls {
				handler(ctx, &tg.Update{
					Message: &tg.Message{Text: url},
				})
			}
			for _, message := range messages {
				handler(ctx, &tg.Update{
					Message: &tg.Message{Text: message},
				})
			}
		}
	})

	b.Run("url.Parse", func(b *testing.B) {
		handler := tg.OnUrl
		ctx := context.Background()
		for b.Loop() {
			for _, url := range urls {
				handler(ctx, &tg.Update{
					Message: &tg.Message{Text: url},
				})
			}
			for _, message := range messages {
				handler(ctx, &tg.Update{
					Message: &tg.Message{Text: message},
				})
			}
		}
	})
}
