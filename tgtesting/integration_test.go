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

// If this test flaps, create an issue.
func TestIntegrationShort(t *testing.T) {
	t.Setenv(tg.EnvTimeoutHandle, "1")
	if chat == 0 {
		t.Skip("no test chat found")
	}
	bot := tg.NewFromEnv()

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

func ffmpegConvert(ctx context.Context, source, target string) ([]byte, error) {
	return exec.CommandContext(ctx, "ffmpeg",
		"-i", source, "-vf", "format=gray", "-c:v", "libx264", "-preset", "ultrafast", "-crf", "23", "-c:a", "copy", target).
		CombinedOutput()
}
