package tg

import (
	"context"
	"net/http"
	"os"
	"reflect"
)

const (
	fieldThumbnail = "Thumbnail"
	fieldCover     = "Cover"
	fieldWidth     = "Width"
	fieldHeight    = "Height"
	fieldDuration  = "Duration"
	fieldLength    = "Length" // SendVideoNote: single square side, instead of Width+Height.
)

// mediaMeta is the result of best-effort local-file inspection: whatever a
// prober could determine without any external dependency. Zero values mean
// "unknown", not "zero" - autoFillMediaMetadata only ever writes positive
// values into the outgoing request.
type mediaMeta struct {
	width, height     int64
	durationSeconds   int64
	tempThumbnailPath string // "" if no thumbnail was generated.
}

// autoFillMediaMetadata inspects request (a pointer to one of the generated
// Send*'s anonymous Request structs) via reflection. If it has exactly one
// InputFile-typed field other than Thumbnail/Cover, and that field holds a
// *LocalFile, it best-effort fills any zero-valued Thumbnail/Width/Height/
// Duration/Length sibling fields it finds by name - never overriding a value
// the caller already set.
//
// This is a best-effort enhancement: any failure (unreadable file, decode
// error, unsupported format) is a silent no-op, never an error surfaced to
// the caller - the real Send* call must behave identically either way.
//
// The returned cleanup func removes any temp thumbnail file created and must
// be deferred by the caller, but only invoked once the whole request
// (including any retries) is done: the multipart body is streamed by a
// background goroutine that reads the temp file while the HTTP call is in
// flight.
func autoFillMediaMetadata[Request any](ctx context.Context, request *Request) (cleanup func()) {
	cleanup = func() {}
	if request == nil || !getOrDefault(ctx, ContextMediaAutofill, true) {
		return cleanup
	}

	structVal := reflect.ValueOf(request).Elem()
	if structVal.Kind() != reflect.Struct {
		return cleanup
	}

	// Skip all file I/O if this Request shape couldn't possibly use any of
	// the fields we'd fill (e.g. SendPhoto, or the []InputMedia shape used
	// by SendMediaGroup).
	if !structVal.FieldByName(fieldThumbnail).IsValid() &&
		!structVal.FieldByName(fieldWidth).IsValid() &&
		!structVal.FieldByName(fieldHeight).IsValid() &&
		!structVal.FieldByName(fieldDuration).IsValid() &&
		!structVal.FieldByName(fieldLength).IsValid() {
		return cleanup
	}

	primary, ok := findPrimaryMediaField(structVal)
	if !ok {
		return cleanup
	}
	local, ok := primary.Interface().(*LocalFile)
	if !ok || local == nil {
		return cleanup // CloudFile, or unset - nothing local to inspect.
	}

	meta := probeLocalMedia(local.Path)
	if meta == nil {
		return cleanup
	}
	if meta.tempThumbnailPath != "" {
		cleanup = func() { _ = os.Remove(meta.tempThumbnailPath) }
	}

	setZeroInt64Field(structVal, fieldWidth, meta.width)
	setZeroInt64Field(structVal, fieldHeight, meta.height)
	setZeroInt64Field(structVal, fieldDuration, meta.durationSeconds)
	if side := min(meta.width, meta.height); side > 0 {
		setZeroInt64Field(structVal, fieldLength, side)
	}
	if meta.tempThumbnailPath != "" {
		setThumbnailField(structVal, meta.tempThumbnailPath)
	}
	return cleanup
}

// findPrimaryMediaField returns the single InputFile-typed field that isn't
// named Thumbnail/Cover. This is what lets the pass work generically across
// every generated Send*'s anonymous Request shape without a type switch per
// method - and what makes it a safe no-op (ok=false) for shapes with zero or
// more than one such field, rather than ever guessing wrong.
func findPrimaryMediaField(structVal reflect.Value) (reflect.Value, bool) {
	inputFileType := reflect.TypeFor[InputFile]()
	structType := structVal.Type()

	var found reflect.Value
	count := 0
	for i := 0; i < structType.NumField(); i++ {
		name := structType.Field(i).Name
		if name == fieldThumbnail || name == fieldCover {
			continue
		}
		field := structVal.Field(i)
		if field.Type() == inputFileType && !field.IsNil() {
			found, count = field, count+1
		}
	}
	return found, count == 1
}

func setZeroInt64Field(structVal reflect.Value, name string, value int64) {
	if value <= 0 {
		return
	}
	field := structVal.FieldByName(name)
	if !field.IsValid() || field.Kind() != reflect.Int64 || !field.CanSet() || !field.IsZero() {
		return
	}
	field.SetInt(value)
}

func setThumbnailField(structVal reflect.Value, path string) {
	field := structVal.FieldByName(fieldThumbnail)
	if !field.IsValid() || field.Type() != reflect.TypeFor[InputFile]() || !field.CanSet() || !field.IsNil() {
		return
	}
	field.Set(reflect.ValueOf(InputFile(FromDisk(path, "thumbnail.jpg"))))
}

// probeLocalMedia sniffs path's content and dispatches to the matching
// prober. Returns nil on any failure - never errors to the caller.
func probeLocalMedia(path string) *mediaMeta {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	header := make([]byte, 512)
	n, _ := file.Read(header)
	_ = file.Close()
	if n < 0 {
		n = 0
	}
	header = header[:n]

	if isISOBMFF(header) {
		meta, err := probeISOBMFF(path)
		if err != nil {
			return nil
		}
		return meta
	}

	switch http.DetectContentType(header) {
	case "image/gif":
		return probeGIF(path)
	case "image/jpeg", "image/png":
		return probeStillImage(path)
	default:
		return nil
	}
}

// isISOBMFF reports whether header looks like an MP4/MOV (ISO Base Media
// File Format) container: a box-size (4 bytes) followed by the "ftyp"
// fourcc at offset 4.
func isISOBMFF(header []byte) bool {
	return len(header) >= 8 && string(header[4:8]) == "ftyp"
}
