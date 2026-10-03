package tg

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// BenchmarkVideoThumbnail compares generating a video thumbnail with this
// package's own pure-Go H.264 decoder (probeISOBMFF's video-frame-decode
// path) against spawning a real ffmpeg subprocess to do the same job -
// the two realistic options for "get a real first-frame thumbnail out of a
// local video file" that autoFillMediaMetadata chooses between.
//
// The ffmpeg variant isn't tuned for an exact pixel-identical match (it asks
// for a similarly-sized, similarly-encoded JPEG via -vf scale/-q:v, not byte-
// identical output) - this benchmark measures wall-clock cost, not output
// fidelity (see TestH264DecodeFirstFrameAgainstFFmpeg for that).
func BenchmarkVideoThumbnail(b *testing.B) {
	cases := []string{
		"tgtesting/testdata/video.mp4",
		"tgtesting/testdata/bigger_2.mp4",
	}
	for _, path := range cases {
		if _, err := os.Stat(path); err != nil {
			b.Logf("skipping %s: %v", path, err)
			continue
		}
		name := filepath.Base(path)

		b.Run(name+"/ours", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				meta, err := probeISOBMFF(path, true)
				if err != nil {
					b.Fatalf("probeISOBMFF: %v", err)
				}
				if meta.tempThumbnailPath == "" {
					b.Fatal("expected a real decoded thumbnail, got none")
				}
				_ = os.Remove(meta.tempThumbnailPath)
			}
		})

		b.Run(name+"/ffmpeg_subprocess", func(b *testing.B) {
			if _, err := exec.LookPath("ffmpeg"); err != nil {
				b.Skip("ffmpeg not installed")
			}
			out := filepath.Join(b.TempDir(), "thumb.jpg")
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command("ffmpeg",
					"-y", "-v", "error",
					"-i", path,
					"-vf", "scale='min(320,iw)':'min(320,ih)':force_original_aspect_ratio=decrease",
					"-frames:v", "1",
					"-q:v", "5",
					out,
				)
				if output, err := cmd.CombinedOutput(); err != nil {
					b.Fatalf("ffmpeg: %v\n%s", err, output)
				}
			}
		})
	}
}
