package tginternalvideo

import (
	"image"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// decodeFirstFrameOfFile runs the real pipeline (box-walk -> locate sample
// -> decode) against a video file on disk, for testing.
func decodeFirstFrameOfFile(t *testing.T, path string) (image.Image, error) {
	t.Helper()
	videoTrak := findVideoTrakForTest(t, path)
	avcCPayload, offset, size, ok := locateH264SampleForTest(videoTrak)
	if !ok {
		t.Fatalf("%s: locateH264SampleForTest failed", path)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer file.Close()
	sample := make([]byte, size)
	if _, err := file.ReadAt(sample, offset); err != nil {
		t.Fatalf("read sample: %v", err)
	}

	return DecodeFirstFrame(avcCPayload, sample)
}

// ffmpegReferenceFrame extracts frame 1 of path via ffmpeg as ground truth.
// Skips the test if ffmpeg isn't installed, matching this repo's existing
// optional-ffmpeg precedent (tgtesting/integration_test.go).
func ffmpegReferenceFrame(t *testing.T, path string) image.Image {
	t.Helper()
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	out := filepath.Join(t.TempDir(), "ref.png")
	cmd := exec.Command("ffmpeg", "-y", "-v", "error", "-i", path, "-frames:v", "1", out)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("ffmpeg failed: %v\n%s", err, output)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatalf("open reference: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decode reference png: %v", err)
	}
	return img
}

// meanAbsDiff computes the mean absolute per-channel pixel difference
// between two same-size images (0-255 scale) - a simple, robust similarity
// metric that tolerates small differences (e.g. from omitting the
// deblocking filter) while still catching a badly wrong decode.
func meanAbsDiff(t *testing.T, got, want image.Image) float64 {
	t.Helper()
	gb, wb := got.Bounds(), want.Bounds()
	if gb.Dx() != wb.Dx() || gb.Dy() != wb.Dy() {
		t.Fatalf("size mismatch: got %dx%d, want %dx%d", gb.Dx(), gb.Dy(), wb.Dx(), wb.Dy())
	}
	var total float64
	var count int
	for y := 0; y < gb.Dy(); y++ {
		for x := 0; x < gb.Dx(); x++ {
			gr, gg, gbv, _ := got.At(gb.Min.X+x, gb.Min.Y+y).RGBA()
			wr, wg, wbv, _ := want.At(wb.Min.X+x, wb.Min.Y+y).RGBA()
			total += math.Abs(float64(gr>>8) - float64(wr>>8))
			total += math.Abs(float64(gg>>8) - float64(wg>>8))
			total += math.Abs(float64(gbv>>8) - float64(wbv>>8))
			count += 3
		}
	}
	return total / float64(count)
}

func saveDebugPNG(t *testing.T, name string, img image.Image) {
	t.Helper()
	path := filepath.Join(os.TempDir(), name)
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_ = png.Encode(f, img)
	t.Logf("wrote debug image to %s", path)
}

func TestH264DecodeFirstFrameAgainstFFmpeg(t *testing.T) {
	cases := []string{
		"../tgtesting/testdata/video.mp4",
		"../tgtesting/testdata/bigger_2.mp4",
	}
	for _, path := range cases {
		t.Run(path, func(t *testing.T) {
			if _, err := os.Stat(path); err != nil {
				t.Skipf("fixture not found: %v", err)
			}

			got, err := decodeFirstFrameOfFile(t, path)
			if err != nil {
				t.Fatalf("DecodeFirstFrame: %v", err)
			}
			saveDebugPNG(t, filepath.Base(path)+".decoded.png", got)

			want := ffmpegReferenceFrame(t, path)
			saveDebugPNG(t, filepath.Base(path)+".reference.png", want)
			diff := meanAbsDiff(t, got, want)
			t.Logf("mean abs per-channel diff vs ffmpeg: %.2f (0-255 scale)", diff)
			// Loose threshold: this decoder skips the deblocking filter
			// entirely (expected small differences at block edges) and its
			// CABAC context tables are empirically validated rather than
			// spec-transcribed (see h264_cabac_tables.go) - a mostly-correct
			// decode should still land well under half the full range.
			if diff > 40 {
				t.Errorf("decoded frame too different from ffmpeg's reference (diff=%.2f) - decoder likely has a correctness bug, see debug PNG", diff)
			}
		})
	}
}
