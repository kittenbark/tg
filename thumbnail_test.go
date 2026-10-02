package tg

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestFitDimensions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                  string
		width, height         int
		maxDim                int
		wantWidth, wantHeight int
	}{
		{"already fits", 100, 200, 320, 100, 200},
		{"exactly fits", 320, 320, 320, 320, 320},
		{"wide downscale", 1280, 720, 320, 320, 180},
		{"tall downscale", 720, 1280, 320, 180, 320},
		{"square downscale", 1000, 1000, 320, 320, 320},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotW, gotH := fitDimensions(tc.width, tc.height, tc.maxDim)
			if gotW != tc.wantWidth || gotH != tc.wantHeight {
				t.Fatalf("fitDimensions(%d,%d,%d) = (%d,%d), want (%d,%d)",
					tc.width, tc.height, tc.maxDim, gotW, gotH, tc.wantWidth, tc.wantHeight)
			}
			if gotW > tc.maxDim || gotH > tc.maxDim {
				t.Fatalf("result (%d,%d) exceeds maxDim %d", gotW, gotH, tc.maxDim)
			}
		})
	}
}

// solidImage builds a width x height image evenly split into a left half of
// one color and a right half of another, so a box-average resize's output
// can be checked against a known expectation.
func solidImage(width, height int, left, right color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x < width/2 {
				img.Set(x, y, left)
			} else {
				img.Set(x, y, right)
			}
		}
	}
	return img
}

func TestResizeToFitNoopWhenAlreadyWithinBounds(t *testing.T) {
	t.Parallel()

	src := solidImage(100, 50, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255})
	dst := resizeToFit(src, thumbnailMaxDim)
	bounds := dst.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 50 {
		t.Fatalf("expected no resize, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeToFitDownscalesAndAverages(t *testing.T) {
	t.Parallel()

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	src := solidImage(640, 640, red, blue)

	dst := resizeToFit(src, thumbnailMaxDim)
	bounds := dst.Bounds()
	if bounds.Dx() != thumbnailMaxDim || bounds.Dy() != thumbnailMaxDim {
		t.Fatalf("got %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), thumbnailMaxDim, thumbnailMaxDim)
	}

	// A pixel deep in the left half should stay purely red; deep in the
	// right half, purely blue - box-averaging only blurs pixels straddling
	// the boundary between the two halves.
	r, g, b, _ := dst.At(10, 10).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Fatalf("left-half pixel = (%d,%d,%d), want pure red", r>>8, g>>8, b>>8)
	}
	r, g, b, _ = dst.At(thumbnailMaxDim-10, 10).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 {
		t.Fatalf("right-half pixel = (%d,%d,%d), want pure blue", r>>8, g>>8, b>>8)
	}
}

func TestWriteJPEGThumbnailRespectsTelegramLimits(t *testing.T) {
	t.Parallel()

	src := solidImage(2000, 1000, color.RGBA{10, 200, 30, 255}, color.RGBA{250, 10, 10, 255})
	path, err := writeJPEGThumbnail(src)
	if err != nil {
		t.Fatalf("writeJPEGThumbnail: %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() > thumbnailMaxBytes {
		t.Fatalf("thumbnail is %d bytes, want <= %d", info.Size(), thumbnailMaxBytes)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer file.Close()
	decoded, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("decoded thumbnail is not a valid JPEG: %v", err)
	}
	bounds := decoded.Bounds()
	if bounds.Dx() > thumbnailMaxDim || bounds.Dy() > thumbnailMaxDim {
		t.Fatalf("thumbnail is %dx%d, want both dims <= %d", bounds.Dx(), bounds.Dy(), thumbnailMaxDim)
	}
}

func TestProbeStillImageJPEG(t *testing.T) {
	t.Parallel()

	src := solidImage(400, 300, color.RGBA{1, 2, 3, 255}, color.RGBA{4, 5, 6, 255})
	path := filepath.Join(t.TempDir(), "photo.jpg")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := jpeg.Encode(file, src, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode: %v", err)
	}
	file.Close()

	meta := probeStillImage(path)
	if meta == nil {
		t.Fatal("expected a non-nil mediaMeta")
	}
	if meta.width != 400 || meta.height != 300 {
		t.Fatalf("got width=%d height=%d, want 400x300", meta.width, meta.height)
	}
	if meta.tempThumbnailPath == "" {
		t.Fatal("expected a generated thumbnail path")
	}
	defer os.Remove(meta.tempThumbnailPath)
}

func TestProbeStillImageRejectsGarbage(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "not-an-image.jpg")
	if err := os.WriteFile(path, []byte("definitely not an image"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if meta := probeStillImage(path); meta != nil {
		t.Fatalf("expected nil for undecodable content, got %+v", meta)
	}
}

func TestProbeGIF(t *testing.T) {
	t.Parallel()

	palette := []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}}
	frame1 := image.NewPaletted(image.Rect(0, 0, 100, 80), palette)
	frame2 := image.NewPaletted(image.Rect(0, 0, 100, 80), palette)
	anim := &gif.GIF{
		Image: []*image.Paletted{frame1, frame2},
		Delay: []int{20, 30}, // centiseconds: 0.5s total.
	}

	path := filepath.Join(t.TempDir(), "anim.gif")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := gif.EncodeAll(file, anim); err != nil {
		t.Fatalf("encode: %v", err)
	}
	file.Close()

	meta := probeGIF(path)
	if meta == nil {
		t.Fatal("expected a non-nil mediaMeta")
	}
	if meta.width != 100 || meta.height != 80 {
		t.Fatalf("got width=%d height=%d, want 100x80", meta.width, meta.height)
	}
	if meta.durationSeconds != 0 {
		t.Fatalf("got durationSeconds=%d, want 0 (50 centiseconds rounds down)", meta.durationSeconds)
	}
	if meta.tempThumbnailPath == "" {
		t.Fatal("expected a generated thumbnail path")
	}
	defer os.Remove(meta.tempThumbnailPath)
}

func TestProbeLocalMediaDispatch(t *testing.T) {
	t.Parallel()

	jpegPath := filepath.Join(t.TempDir(), "photo.jpg")
	file, err := os.Create(jpegPath)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := jpeg.Encode(file, solidImage(50, 50, color.RGBA{1, 1, 1, 255}, color.RGBA{2, 2, 2, 255}), nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	file.Close()

	meta := probeLocalMedia(jpegPath)
	if meta == nil || meta.width != 50 || meta.height != 50 {
		t.Fatalf("got %+v, want a 50x50 image probe", meta)
	}
	if meta.tempThumbnailPath != "" {
		defer os.Remove(meta.tempThumbnailPath)
	}

	unknownPath := filepath.Join(t.TempDir(), "random.bin")
	if err := os.WriteFile(unknownPath, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if meta := probeLocalMedia(unknownPath); meta != nil {
		t.Fatalf("expected nil for an unrecognized format, got %+v", meta)
	}
}
