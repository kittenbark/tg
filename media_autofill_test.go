package tg

import (
	"context"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

// These mirror (a trimmed-down version of) the shape of the generated
// anonymous Request structs in api_methods.go - autoFillMediaMetadata must
// work against them purely by reflection + field name, with no dependency
// on the real generated types.
type fakeSendVideoRequest struct {
	ChatId    int64     `json:"chat_id"`
	Video     InputFile `json:"video"`
	Duration  int64     `json:"duration,omitempty"`
	Width     int64     `json:"width,omitempty"`
	Height    int64     `json:"height,omitempty"`
	Cover     InputFile `json:"cover,omitempty"`
	Thumbnail InputFile `json:"thumbnail,omitempty"`
}

type fakeSendDocumentRequest struct {
	ChatId    int64     `json:"chat_id"`
	Document  InputFile `json:"document"`
	Thumbnail InputFile `json:"thumbnail,omitempty"`
}

type fakeSendAnimationRequest struct {
	ChatId    int64     `json:"chat_id"`
	Animation InputFile `json:"animation"`
	Duration  int64     `json:"duration,omitempty"`
	Width     int64     `json:"width,omitempty"`
	Height    int64     `json:"height,omitempty"`
	Thumbnail InputFile `json:"thumbnail,omitempty"`
}

type fakeSendPhotoRequest struct {
	ChatId int64     `json:"chat_id"`
	Photo  InputFile `json:"photo"`
}

type fakeAmbiguousRequest struct {
	ChatId int64     `json:"chat_id"`
	Video  InputFile `json:"video"`
	Extra  InputFile `json:"extra"`
}

func writeJPEGFile(t *testing.T, path string, width, height int) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}
	if err := jpeg.Encode(file, img, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func writeGIFFile(t *testing.T, path string, width, height int) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer file.Close()
	palette := []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}
	frame := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	anim := &gif.GIF{Image: []*image.Paletted{frame, frame}, Delay: []int{100, 100}} // 2 seconds total.
	if err := gif.EncodeAll(file, anim); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func TestAutoFillMediaMetadataVideo(t *testing.T) {
	t.Parallel()

	path := writeSyntheticMP4(t, 720, 1280, 1000, 5000)
	req := &fakeSendVideoRequest{ChatId: 1, Video: FromDisk(path)}

	cleanup := autoFillMediaMetadata(context.Background(), req)
	defer cleanup()

	if req.Width != 720 || req.Height != 1280 || req.Duration != 5 {
		t.Fatalf("got %+v, want Width=720 Height=1280 Duration=5", req)
	}
	if req.Thumbnail != nil {
		t.Fatalf("expected no thumbnail for an MP4 video, got %+v", req.Thumbnail)
	}
}

func TestAutoFillMediaMetadataNeverOverrides(t *testing.T) {
	t.Parallel()

	path := writeSyntheticMP4(t, 720, 1280, 1000, 5000)
	req := &fakeSendVideoRequest{ChatId: 1, Video: FromDisk(path), Width: 999}

	cleanup := autoFillMediaMetadata(context.Background(), req)
	defer cleanup()

	if req.Width != 999 {
		t.Fatalf("got Width=%d, want the caller-supplied 999 preserved", req.Width)
	}
	if req.Height != 1280 || req.Duration != 5 {
		t.Fatalf("got %+v, want Height/Duration still autofilled", req)
	}
}

func TestAutoFillMediaMetadataDocumentImage(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "doc.jpg")
	writeJPEGFile(t, path, 200, 100)
	req := &fakeSendDocumentRequest{ChatId: 1, Document: FromDisk(path)}

	cleanup := autoFillMediaMetadata(context.Background(), req)
	defer cleanup()

	if req.Thumbnail == nil {
		t.Fatal("expected a generated thumbnail for a JPEG document")
	}
}

func TestAutoFillMediaMetadataAnimationGIF(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "anim.gif")
	writeGIFFile(t, path, 64, 48)
	req := &fakeSendAnimationRequest{ChatId: 1, Animation: FromDisk(path)}

	cleanup := autoFillMediaMetadata(context.Background(), req)
	defer cleanup()

	if req.Width != 64 || req.Height != 48 {
		t.Fatalf("got Width=%d Height=%d, want 64x48", req.Width, req.Height)
	}
	if req.Duration != 2 {
		t.Fatalf("got Duration=%d, want 2 (200 centiseconds)", req.Duration)
	}
	if req.Thumbnail == nil {
		t.Fatal("expected a generated thumbnail for a GIF animation")
	}
}

func TestAutoFillMediaMetadataCloudFileIsNoop(t *testing.T) {
	t.Parallel()

	req := &fakeSendVideoRequest{ChatId: 1, Video: FromCloud("some-file-id")}
	cleanup := autoFillMediaMetadata(context.Background(), req)
	defer cleanup()

	if req.Width != 0 || req.Height != 0 || req.Duration != 0 || req.Thumbnail != nil {
		t.Fatalf("expected no changes for a CloudFile input, got %+v", req)
	}
}

func TestAutoFillMediaMetadataDisabledViaContext(t *testing.T) {
	t.Parallel()

	path := writeSyntheticMP4(t, 720, 1280, 1000, 5000)
	req := &fakeSendVideoRequest{ChatId: 1, Video: FromDisk(path)}

	ctx := context.WithValue(context.Background(), ContextMediaAutofill, false)
	cleanup := autoFillMediaMetadata(ctx, req)
	defer cleanup()

	if req.Width != 0 || req.Height != 0 || req.Duration != 0 || req.Thumbnail != nil {
		t.Fatalf("expected no changes when ContextMediaAutofill is false, got %+v", req)
	}
}

func TestAutoFillMediaMetadataPhotoRequestIsNoop(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "photo.jpg")
	writeJPEGFile(t, path, 100, 100)
	req := &fakeSendPhotoRequest{ChatId: 1, Photo: FromDisk(path)}

	// Must not panic even though this shape has none of the fillable fields.
	cleanup := autoFillMediaMetadata(context.Background(), req)
	cleanup()
}

func TestAutoFillMediaMetadataAmbiguousPrimaryFieldIsNoop(t *testing.T) {
	t.Parallel()

	path := writeSyntheticMP4(t, 720, 1280, 1000, 5000)
	req := &fakeAmbiguousRequest{ChatId: 1, Video: FromDisk(path), Extra: FromDisk(path)}

	// Two qualifying InputFile fields -> findPrimaryMediaField must refuse to
	// guess, so this struct's lack of Width/Height/etc fields is exercised
	// as a no-op (can't observe field mutation here, just that it doesn't panic).
	cleanup := autoFillMediaMetadata(context.Background(), req)
	cleanup()
}
