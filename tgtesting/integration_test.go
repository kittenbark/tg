package tgtesting

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kittenbark/tg"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

var (
	chat, _  = strconv.ParseInt(os.Getenv(tg.EnvTestingChat), 10, 64)
	group, _ = strconv.ParseInt(os.Getenv(tg.EnvTestingGroupChat), 10, 64)
	photo    = OutsideFile("./testdata/photo.jpg", "https://github.com/kittenbark/testdata/raw/dafdaa3ec6f42ecacd5d04e8c0ccd39ba9e70f28/photo.jpg")
	video    = OutsideFile("./testdata/video.mp4", "https://github.com/kittenbark/testdata/raw/dafdaa3ec6f42ecacd5d04e8c0ccd39ba9e70f28/video.mp4")
	dir      = "./testdata"
)

// If this test flaps, create an issue.
func TestIntegrationShort(t *testing.T) {
	t.Setenv(tg.EnvTimeoutHandle, "1.5")
	if chat == 0 {
		t.Skip("no test chat found")
	}
	bot := tg.NewFromEnv().
		Help("/start", "this is start fr (private chat)").
		HelpScoped(&tg.BotCommandScopeAllGroupChats{},
			"/start", "this is start fr (group chat)",
		)

	t.Run("setMessageReaction", func(t *testing.T) {
		t.Parallel()

		sent, err := tg.SendMessage(bot.Context(), chat, ":heart: common")
		require.NoError(t, err)

		ok, err := tg.SetMessageReaction(bot.Context(), sent.Chat.Id, sent.MessageId, tg.CommonReaction("❤"))
		require.NoError(t, err)
		require.Equal(t, true, ok)
	})

	t.Run("setMessageReaction#default", func(t *testing.T) {
		t.Parallel()

		sent, err := tg.SendMessage(bot.Context(), chat, ":fire: default")
		require.NoError(t, err)

		ok, err := tg.SetMessageReaction(bot.Context(), sent.Chat.Id, sent.MessageId, &tg.OptSetMessageReaction{
			Reaction: []tg.ReactionType{&tg.ReactionTypeEmoji{Emoji: "🔥"}},
		})
		require.NoError(t, err)
		require.Equal(t, true, ok)
	})

	t.Run("ParseMode", func(t *testing.T) {
		t.Parallel()

		_, err := tg.SendMessage(
			bot.Context(),
			chat,
			fmt.Sprintf("```python\n%s\n```", tg.EscapeParseMode(tg.ParseModeMarkdownV2, "print('hello world')")),
			&tg.OptSendMessage{ParseMode: tg.ParseModeMarkdownV2},
		)
		require.NoError(t, err)

		html := "<b>bold</b>, <strong>bold</strong>\n<i>italic</i>, <em>italic</em>\n<u>underline</u>, <ins>underline</ins>\n<s>strikethrough</s>, <strike>strikethrough</strike>, <del>strikethrough</del>\n<span class=\"tg-spoiler\">spoiler</span>, <tg-spoiler>spoiler</tg-spoiler>\n<b>bold <i>italic bold <s>italic bold strikethrough <span class=\"tg-spoiler\">italic bold strikethrough spoiler</span></s> <u>underline italic bold</u></i> bold</b>\n<a href=\"http://www.example.com/\">inline URL</a>\n<a href=\"tg://user?id=123456789\">inline mention of a user</a>\n<tg-emoji emoji-id=\"5368324170671202286\">👍</tg-emoji>\n<code>inline fixed-width code</code>\n<pre>pre-formatted fixed-width code block</pre>\n<pre><code class=\"language-python\">pre-formatted fixed-width code block written in the Python programming language</code></pre>\n<blockquote>Block quotation started\\nBlock quotation continued\\nThe last line of the block quotation</blockquote>\n<blockquote expandable>Expandable block quotation started\\nExpandable block quotation continued\\nExpandable block quotation continued\\nHidden by default part of the block quotation started\\nExpandable block quotation continued\\nThe last line of the block quotation</blockquote>"
		_, err = tg.SendMessage(bot.Context(), chat, html, &tg.OptSendMessage{
			ParseMode: tg.ParseModeHTML,
		})
		require.NoError(t, err)
		_, err = tg.SendMessage(bot.Context(), chat, tg.EscapeParseMode(tg.ParseModeHTML, html), &tg.OptSendMessage{
			ParseMode: tg.ParseModeHTML,
		})
		require.NoError(t, err)

		markdown := "*bold text*\n_italic text_\n[inline URL](http://www.example.com/)\n[inline mention of a user](tg://user?id=123456789)\n`inline fixed-width code`\n```\npre-formatted fixed-width code block\n```\n```python\npre-formatted fixed-width code block written in the Python programming language\n```"
		_, err = tg.SendMessage(bot.Context(), chat, markdown, &tg.OptSendMessage{
			ParseMode: tg.ParseModeMarkdown,
		})
		require.NoError(t, err)
		_, err = tg.SendMessage(bot.Context(), chat, tg.EscapeParseMode(tg.ParseModeMarkdown, markdown), &tg.OptSendMessage{
			ParseMode: tg.ParseModeMarkdown,
		})
		require.NoError(t, err)

		markdownV2 := "*bold \\*text*\n_italic \\*text_\n__underline__\n~strikethrough~\n||spoiler||\n*bold _italic bold ~italic bold strikethrough ||italic bold strikethrough spoiler||~ __underline italic bold___ bold*\n[inline URL](http://www.example.com/)\n[inline mention of a user](tg://user?id=123456789)\n![👍](tg://emoji?id=5368324170671202286)\n`inline fixed-width code`\n```\npre-formatted fixed-width code block\n```\n```python\npre-formatted fixed-width code block written in the Python programming language\n```\n>Block quotation started\n>Block quotation continued\n>Block quotation continued\n>Block quotation continued\n>The last line of the block quotation\n**>The expandable block quotation started right after the previous block quotation\n>It is separated from the previous block quotation by an empty bold entity\n>Expandable block quotation continued\n>Hidden by default part of the expandable block quotation started\n>Expandable block quotation continued\n>The last line of the expandable block quotation with the expandability mark||"
		_, err = tg.SendMessage(bot.Context(), chat, markdownV2, &tg.OptSendMessage{
			ParseMode: tg.ParseModeMarkdownV2,
		})
		require.NoError(t, err)
		_, err = tg.SendMessage(bot.Context(), chat, tg.EscapeParseMode(tg.ParseModeMarkdownV2, markdownV2), &tg.OptSendMessage{
			ParseMode: tg.ParseModeMarkdownV2,
		})
		require.NoError(t, err)
	})
}

func TestIntegrationLong(t *testing.T) {
	if photo == "" || video == "" {
		t.Skip("no test photo/video found")
	}
	if chat == 0 {
		t.Skip("no test chat found")
	}

	t.Setenv(tg.EnvTimeoutHandle, "120")
	bot := tg.NewFromEnv()

	t.Run("sendMediaGroup", func(t *testing.T) {
		t.Parallel()
		ctx := bot.Context()

		thisIsKot := "this is kot"

		album, err := tg.SendMediaGroup(ctx, chat, tg.Album{
			&tg.Photo{Media: tg.FromDisk(photo)},
			&tg.Video{Media: tg.FromDisk(video), Caption: thisIsKot, Thumbnail: tg.FromDisk(photo)},
		})
		require.NoError(t, err)

		require.Equal(t, 2, len(album))

		require.Equal(t, "", album[0].Caption)
		require.NotNil(t, album[0].Photo)

		require.Equal(t, thisIsKot, album[1].Caption)
		require.NotNil(t, album[1].Video)
	})

	t.Run("sendVideo", func(t *testing.T) {
		t.Parallel()
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			t.Skip("ffmpeg not installed")
		}
		ctx := bot.Context()

		sent, err := tg.SendVideo(ctx, chat, tg.FromDisk(video), &tg.OptSendVideo{Caption: "this is mp4"})
		require.NoError(t, err)
		require.Equal(t, "this is mp4", sent.Caption)

		target := path.Join(dir, "h264_"+path.Base(video))
		if out, err := ffmpegConvert(ctx, video, target); err != nil {
			slog.Error("convert error", "err", err, "output", string(out))
		}
		defer os.Remove(target)

		downloaded, err := sent.Video.DownloadTemp(ctx, dir)
		require.NoError(t, err)
		defer os.Remove(downloaded)

		converted := "new_" + path.Base(target)
		defer os.Remove(converted)
		out, err := ffmpegConvert(ctx, downloaded, converted)
		require.NoError(t, err, "stdout: "+string(out))

		_, err = tg.SendVideo(ctx, chat, tg.FromDisk(converted))
		require.NoError(t, err)
	})
}

func TestDDOS(t *testing.T) {
	bot := tg.NewFromEnv().Scheduler()

	start := time.Now()
	i := atomic.Int64{}
	for {
		go func() {
			i.Add(1)
			_, err := tg.SendMessage(bot.Context(), group, strconv.Itoa(int(i.Load())))
			require.NoError(t, err)
		}()
		go func() {
			i.Add(1)
			_, err := tg.SendMessage(bot.Context(), chat, strconv.Itoa(int(i.Load())))
			require.NoError(t, err)
		}()
		if time.Since(start).Seconds() > 10 {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func TestIntegrationHandleAlbum(t *testing.T) {
	bot := tg.NewFromEnv()
	ctx, _ := bot.ContextWithCancel()
	pic, err := tg.SendPhoto(ctx, chat, tg.FromDisk(photo))
	require.NoError(t, err)
	vid, err := tg.SendVideo(ctx, chat, tg.FromDisk(video))
	require.NoError(t, err)

	counter := &atomic.Int64{}
	time.AfterFunc(time.Second*6, bot.Stop)
	bot.
		Plugin(tg.PluginLogger(slog.LevelDebug)).
		OnError(tg.OnErrorPanic).
		Filter(tg.OnPrivateMessage).
		Default(tg.HandleAlbum(func(ctx context.Context, updates []*tg.Update) error {
			defer counter.Add(1)
			album := tg.Album{}
			for _, upd := range updates {
				msg := upd.Message
				if msg.Photo != nil {
					album = append(album, &tg.Photo{Media: tg.FromCloud(msg.Photo.FileId())})
				}
				if msg.Video != nil {
					album = append(album, &tg.Video{Media: tg.FromCloud(msg.Video.FileId)})
				}
			}

			_, err := tg.SendMediaGroup(ctx, updates[0].Message.Chat.Id, album)
			return err
		})).
		Start(
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, From: &tg.User{Id: chat}, Photo: pic.Photo}},
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, From: &tg.User{Id: chat}, Video: vid.Video}},
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, From: &tg.User{Id: chat}, Video: vid.Video, MediaGroupId: "1"}},
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, From: &tg.User{Id: chat}, Video: vid.Video, MediaGroupId: "1"}},
		)

	require.Equal(t, int64(3), counter.Load())
}

func TestIntegrationHandleCallback(t *testing.T) {
	t.Skip("integration testing could be done only by hand, afaik")

	type Info struct {
		Id     int64  `json:"i"`
		ChatId int64  `json:"c"`
		Str    string `json:"s"`
	}

	tg.NewFromEnv().
		OnError(tg.OnErrorPanic).
		Plugin(tg.PluginLogger(slog.LevelDebug)).
		Branch(tg.OnCallbackWithData[Info](), func(ctx context.Context, upd *tg.Update) error {
			value, err := tg.CallbackData[Info](upd)
			if err != nil {
				return err
			}
			callback := upd.CallbackQuery
			_, err = tg.SendMessage(ctx, callback.From.Id, fmt.Sprintf("%#v", value))
			_, err = tg.AnswerCallbackQuery(ctx, callback.Id, &tg.OptAnswerCallbackQuery{Text: "love you"})

			require.Equal(t, chat, value.ChatId)
			require.Equal(t, "msg_text", value.Str)

			return err
		}).
		Branch(tg.OnCallback, func(ctx context.Context, upd *tg.Update) error {
			println(upd.CallbackQuery.Data)
			return nil
		}).
		Filter(tg.OnPrivateMessage).
		Branch(tg.OnText, func(ctx context.Context, upd *tg.Update) error {
			msg := upd.Message
			data, _ := json.Marshal(&Info{Id: msg.MessageId, ChatId: msg.Chat.Id, Str: msg.Text})
			println(string(data))
			_, err := tg.SendMessage(ctx, msg.Chat.Id, msg.Text, &tg.OptSendMessage{
				ReplyMarkup: &tg.InlineKeyboardMarkup{
					InlineKeyboard: [][]*tg.InlineKeyboardButton{{{
						Text:         "button",
						CallbackData: string(data),
					}}},
				},
			})
			return err
		}).
		Start()
}

func TestIntegrationOnNewGroup(t *testing.T) {
	bot := tg.NewFromEnv()
	time.AfterFunc(time.Second*6, bot.Stop)

	counter := &atomic.Int64{}
	bot.
		OnError(tg.OnErrorPanic).
		Plugin(tg.PluginLogger(slog.LevelDebug)).
		Branch(tg.OnAddedToGroup, func(ctx context.Context, upd *tg.Update) error {
			defer counter.Add(1)
			msg := upd.Message

			_, _ = tg.SendMessage(ctx, chat, "new chat")
			_, err := tg.SendMessage(ctx, msg.Chat.Id, "hello chat")
			return err
		}).
		Command("/start", tg.CommonTextReply("+++")).
		Start(
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, GroupChatCreated: true}},
			&tg.Update{UpdateId: -1, Message: &tg.Message{Chat: &tg.Chat{Id: chat}, Text: "/start"}},
		)

	require.Equal(t, int64(1), counter.Load())
}

func TestContext(t *testing.T) {
	localUrl := "http://localhost:8080"
	t.Setenv(tg.EnvApiURL, localUrl)

	bot := tg.NewFromEnv()
	require.Equal(t, localUrl, bot.Context().Value(tg.ContextApiUrl).(string))
}

func ffmpegConvert(ctx context.Context, source, target string) ([]byte, error) {
	return exec.CommandContext(ctx, "ffmpeg",
		"-y", "-i", source, "-vf", "format=gray", "-c:v", "libx264", "-preset", "ultrafast", "-crf", "23", "-c:a", "copy", target).
		CombinedOutput()
}
