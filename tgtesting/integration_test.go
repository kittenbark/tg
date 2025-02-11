package tgtesting

import (
	"context"
	"github.com/kittenbark/tg"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"
)

var (
	chat, _ = strconv.ParseInt(os.Getenv(tg.EnvTestingChat), 10, 64)
	photo   = OutsideFile("./testdata/photo.jpg", "https://github.com/kittenbark/testdata/raw/dafdaa3ec6f42ecacd5d04e8c0ccd39ba9e70f28/photo.jpg")
	video   = OutsideFile("./testdata/video.mp4", "https://github.com/kittenbark/testdata/raw/dafdaa3ec6f42ecacd5d04e8c0ccd39ba9e70f28/video.mp4")
	dir     = "./testdata"
)

func TestIntegrationVideo(t *testing.T) {
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
			&tg.Video{Media: tg.FromDisk(video), Caption: thisIsKot},
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

		sent, err := tg.SendVideo(ctx, chat, tg.FromDisk(video))
		require.NoError(t, err)

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

func _() {
	tg.NewFromEnv().
		Filter(tg.IsPrivateMessage).
		HandleCommand("/start", tg.CommonTextReply("hello this echo bot is made with @kittenbark_tg")).
		HandleCommand("/help", tg.CommonTextReply("just send a message")).
		Branch(tg.OnMessage, func(ctx context.Context, upd *tg.Update) error {
			msg := upd.Message
			_, err := tg.CopyMessage(ctx, msg.Chat.Id, msg.Chat.Id, msg.MessageId)
			return err
		}).
		Default(tg.CommonTextReply("unsupported message type, /help?")).
		Start()
}

func ffmpegConvert(ctx context.Context, source, target string) ([]byte, error) {
	return exec.CommandContext(ctx, "ffmpeg",
		"-i", source, "-vf", "format=gray", "-c:v", "libx264", "-preset", "ultrafast", "-crf", "23", "-c:a", "copy", target).
		CombinedOutput()
}
