package tg

import (
	"bytes"
	"encoding/binary"
	"image"
	_ "image/jpeg" // registers the JPEG decoder for image.Decode in TestProbeISOBMFFVideoFrameDecode.
	"os"
	"path/filepath"
	"testing"
)

func buildBox(boxType string, payload []byte) []byte {
	box := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(box[0:4], uint32(8+len(payload)))
	copy(box[4:8], boxType)
	copy(box[8:], payload)
	return box
}

func buildHdlr(handlerType string) []byte {
	payload := make([]byte, 24)
	copy(payload[8:12], handlerType)
	return payload
}

func buildMvhd(version byte, timescale uint32, duration uint64) []byte {
	if version == 1 {
		payload := make([]byte, 32)
		payload[0] = 1
		binary.BigEndian.PutUint32(payload[20:24], timescale)
		binary.BigEndian.PutUint64(payload[24:32], duration)
		return payload
	}
	payload := make([]byte, 20)
	payload[0] = 0
	binary.BigEndian.PutUint32(payload[12:16], timescale)
	binary.BigEndian.PutUint32(payload[16:20], uint32(duration))
	return payload
}

func buildTkhd(version byte, width, height uint32) []byte {
	offset := 76
	if version == 1 {
		offset = 88
	}
	payload := make([]byte, offset+8)
	payload[0] = version
	binary.BigEndian.PutUint32(payload[offset:offset+4], width<<16)
	binary.BigEndian.PutUint32(payload[offset+4:offset+8], height<<16)
	return payload
}

func TestReadBoxes(t *testing.T) {
	t.Parallel()

	buf := append(buildBox("ftyp", []byte("isom0001")), buildBox("free", nil)...)
	boxes := readBoxes(buf)
	if len(boxes) != 2 {
		t.Fatalf("got %d boxes, want 2", len(boxes))
	}
	if boxes[0].boxType != "ftyp" || !bytes.Equal(boxes[0].payload, []byte("isom0001")) {
		t.Fatalf("unexpected first box: %+v", boxes[0])
	}
	if boxes[1].boxType != "free" || len(boxes[1].payload) != 0 {
		t.Fatalf("unexpected second box: %+v", boxes[1])
	}
}

func TestReadBoxesExtendedSize(t *testing.T) {
	t.Parallel()

	payload := []byte("payload-data")
	box := make([]byte, 16+len(payload))
	binary.BigEndian.PutUint32(box[0:4], 1) // size==1 marks a 64-bit extended size.
	copy(box[4:8], "mdat")
	binary.BigEndian.PutUint64(box[8:16], uint64(16+len(payload)))
	copy(box[16:], payload)

	boxes := readBoxes(box)
	if len(boxes) != 1 {
		t.Fatalf("got %d boxes, want 1", len(boxes))
	}
	if boxes[0].boxType != "mdat" || !bytes.Equal(boxes[0].payload, payload) {
		t.Fatalf("unexpected box: %+v", boxes[0])
	}
}

func TestReadBoxesTruncatedIsDropped(t *testing.T) {
	t.Parallel()

	buf := buildBox("ftyp", []byte("isom"))
	buf = append(buf, 0, 0, 0) // a truncated trailing box header.
	boxes := readBoxes(buf)
	if len(boxes) != 1 {
		t.Fatalf("got %d boxes, want 1 (truncated trailer should be dropped)", len(boxes))
	}
}

func TestParseMvhd(t *testing.T) {
	t.Parallel()

	t.Run("version0", func(t *testing.T) {
		duration, ok := parseMvhd(buildMvhd(0, 1000, 5000))
		if !ok || duration != 5 {
			t.Fatalf("got (%d, %v), want (5, true)", duration, ok)
		}
	})

	t.Run("version1", func(t *testing.T) {
		duration, ok := parseMvhd(buildMvhd(1, 2000, 10000))
		if !ok || duration != 5 {
			t.Fatalf("got (%d, %v), want (5, true)", duration, ok)
		}
	})

	t.Run("zero timescale is rejected", func(t *testing.T) {
		if _, ok := parseMvhd(buildMvhd(0, 0, 5000)); ok {
			t.Fatal("expected ok=false for zero timescale")
		}
	})

	t.Run("too short is rejected", func(t *testing.T) {
		if _, ok := parseMvhd([]byte{0, 0, 0}); ok {
			t.Fatal("expected ok=false for a too-short payload")
		}
	})
}

func TestParseTkhd(t *testing.T) {
	t.Parallel()

	t.Run("version0", func(t *testing.T) {
		width, height, ok := parseTkhd(buildTkhd(0, 720, 1280))
		if !ok || width != 720 || height != 1280 {
			t.Fatalf("got (%d, %d, %v), want (720, 1280, true)", width, height, ok)
		}
	})

	t.Run("version1", func(t *testing.T) {
		width, height, ok := parseTkhd(buildTkhd(1, 1080, 1920))
		if !ok || width != 1080 || height != 1920 {
			t.Fatalf("got (%d, %d, %v), want (1080, 1920, true)", width, height, ok)
		}
	})

	t.Run("too short is rejected", func(t *testing.T) {
		if _, _, ok := parseTkhd([]byte{0, 0, 0}); ok {
			t.Fatal("expected ok=false for a too-short payload")
		}
	})
}

func TestIsVideoTrak(t *testing.T) {
	t.Parallel()

	videoTrak := append(buildBox("mdia", buildBox("hdlr", buildHdlr("vide"))), buildBox("tkhd", buildTkhd(0, 100, 100))...)
	if !isVideoTrak(videoTrak) {
		t.Fatal("expected a 'vide' handler trak to be recognized as the video track")
	}

	audioTrak := append(buildBox("mdia", buildBox("hdlr", buildHdlr("soun"))), buildBox("tkhd", buildTkhd(0, 100, 100))...)
	if isVideoTrak(audioTrak) {
		t.Fatal("expected a 'soun' handler trak to NOT be recognized as the video track")
	}
}

// writeSyntheticMP4 assembles a minimal, hand-built ISO-BMFF file (ftyp +
// mdat + moov{mvhd, trak-video, trak-audio}) exercising the same top-level
// Seek/ReadAt scan and nested box traversal probeISOBMFF uses against real
// files, without depending on an actual video fixture.
func writeSyntheticMP4(t *testing.T, width, height uint32, timescale uint32, duration uint64) string {
	t.Helper()

	mvhd := buildBox("mvhd", buildMvhd(0, timescale, duration))

	videoMdia := buildBox("mdia", buildBox("hdlr", buildHdlr("vide")))
	videoTkhd := buildBox("tkhd", buildTkhd(0, width, height))
	videoTrak := buildBox("trak", append(append([]byte{}, videoMdia...), videoTkhd...))

	audioMdia := buildBox("mdia", buildBox("hdlr", buildHdlr("soun")))
	audioTkhd := buildBox("tkhd", buildTkhd(0, 1, 1)) // Must be ignored in favor of the video track.
	audioTrak := buildBox("trak", append(append([]byte{}, audioMdia...), audioTkhd...))

	moovPayload := append(append(append([]byte{}, mvhd...), videoTrak...), audioTrak...)
	moov := buildBox("moov", moovPayload)
	ftyp := buildBox("ftyp", []byte("isom0001"))
	// mdat simulates the real (large, irrelevant) media payload that must
	// never be read into memory by the top-level scan.
	mdat := buildBox("mdat", make([]byte, 1024))

	file := append(append(append([]byte{}, ftyp...), mdat...), moov...)

	path := filepath.Join(t.TempDir(), "synthetic.mp4")
	if err := os.WriteFile(path, file, 0o644); err != nil {
		t.Fatalf("writing synthetic mp4: %v", err)
	}
	return path
}

func TestProbeISOBMFFSynthetic(t *testing.T) {
	t.Parallel()

	path := writeSyntheticMP4(t, 720, 1280, 1000, 5000)
	meta, err := probeISOBMFF(path, true)
	if err != nil {
		t.Fatalf("probeISOBMFF: %v", err)
	}
	if meta.width != 720 || meta.height != 1280 || meta.durationSeconds != 5 {
		t.Fatalf("got %+v, want width=720 height=1280 durationSeconds=5", meta)
	}
	if meta.tempThumbnailPath != "" {
		t.Fatalf("synthetic mp4 has no real H.264 sample table, decode must fail open with no thumbnail, got %q", meta.tempThumbnailPath)
	}
}

func TestProbeISOBMFFMoovTrailing(t *testing.T) {
	t.Parallel()

	// Same shape as writeSyntheticMP4 but with moov placed AFTER a large
	// mdat, as produced by some non-fast-start encoders.
	mvhd := buildBox("mvhd", buildMvhd(0, 1000, 3000))
	mdia := buildBox("mdia", buildBox("hdlr", buildHdlr("vide")))
	tkhd := buildBox("tkhd", buildTkhd(0, 320, 240))
	trak := buildBox("trak", append(append([]byte{}, mdia...), tkhd...))
	moov := buildBox("moov", append(append([]byte{}, mvhd...), trak...))
	ftyp := buildBox("ftyp", []byte("isom0001"))
	mdat := buildBox("mdat", make([]byte, 4096))

	file := append(append(append([]byte{}, ftyp...), mdat...), moov...)
	path := filepath.Join(t.TempDir(), "trailing-moov.mp4")
	if err := os.WriteFile(path, file, 0o644); err != nil {
		t.Fatalf("writing synthetic mp4: %v", err)
	}

	meta, err := probeISOBMFF(path, true)
	if err != nil {
		t.Fatalf("probeISOBMFF: %v", err)
	}
	if meta.width != 320 || meta.height != 240 || meta.durationSeconds != 3 {
		t.Fatalf("got %+v, want width=320 height=240 durationSeconds=3", meta)
	}
}

func TestProbeISOBMFFNoMoov(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "no-moov.mp4")
	if err := os.WriteFile(path, buildBox("ftyp", []byte("isom0001")), 0o644); err != nil {
		t.Fatalf("writing synthetic mp4: %v", err)
	}

	if _, err := probeISOBMFF(path, true); err == nil {
		t.Fatal("expected an error when no moov box is present")
	}
}

func TestProbeISOBMFFRealFixtures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		path                  string
		wantWidth, wantHeight int64
		wantDurationSeconds   int64
	}{
		{"tgtesting/testdata/video.mp4", 720, 1280, 12},
		{"tgtesting/testdata/bigger.mp4", 1080, 1920, 15},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			if _, err := os.Stat(tc.path); err != nil {
				t.Skipf("fixture not found: %v", err)
			}

			meta, err := probeISOBMFF(tc.path, false)
			if err != nil {
				t.Fatalf("probeISOBMFF: %v", err)
			}
			if meta.width != tc.wantWidth || meta.height != tc.wantHeight {
				t.Fatalf("got width=%d height=%d, want width=%d height=%d", meta.width, meta.height, tc.wantWidth, tc.wantHeight)
			}
			if diff := meta.durationSeconds - tc.wantDurationSeconds; diff < -2 || diff > 2 {
				t.Fatalf("duration %d too far from expected %d", meta.durationSeconds, tc.wantDurationSeconds)
			}
			if meta.tempThumbnailPath != "" {
				t.Fatalf("probeISOBMFF with videoFrameDecode=false must never produce a thumbnail, got %q", meta.tempThumbnailPath)
			}
		})
	}
}

// TestProbeISOBMFFVideoFrameDecode exercises the opt-in H.264 first-frame
// thumbnail path end to end against real fixtures: a real H.264 track must
// produce a real, valid, within-limits thumbnail, and an HEVC track (out of
// this decoder's v1 scope) must fail open to metadata-only, exactly as it
// already does with the feature disabled.
func TestProbeISOBMFFVideoFrameDecode(t *testing.T) {
	t.Parallel()

	t.Run("video.mp4 (H.264)", func(t *testing.T) {
		t.Parallel()
		path := "tgtesting/testdata/video.mp4"
		if _, err := os.Stat(path); err != nil {
			t.Skipf("fixture not found: %v", err)
		}

		meta, err := probeISOBMFF(path, true)
		if err != nil {
			t.Fatalf("probeISOBMFF: %v", err)
		}
		if meta.tempThumbnailPath == "" {
			t.Fatal("expected a real decoded thumbnail for an H.264 fixture, got none")
		}
		defer func() { _ = os.Remove(meta.tempThumbnailPath) }()

		info, err := os.Stat(meta.tempThumbnailPath)
		if err != nil {
			t.Fatalf("stat thumbnail: %v", err)
		}
		if info.Size() == 0 || info.Size() > thumbnailMaxBytes {
			t.Fatalf("thumbnail size %d out of bounds (want >0 and <=%d)", info.Size(), thumbnailMaxBytes)
		}
		f, err := os.Open(meta.tempThumbnailPath)
		if err != nil {
			t.Fatalf("open thumbnail: %v", err)
		}
		defer func() { _ = f.Close() }()
		img, _, err := image.Decode(f)
		if err != nil {
			t.Fatalf("decoded thumbnail is not a valid image: %v", err)
		}
		bounds := img.Bounds()
		if bounds.Dx() > thumbnailMaxDim || bounds.Dy() > thumbnailMaxDim {
			t.Fatalf("thumbnail %dx%d exceeds thumbnailMaxDim=%d", bounds.Dx(), bounds.Dy(), thumbnailMaxDim)
		}
	})

	t.Run("bigger.mp4 (HEVC, out of scope)", func(t *testing.T) {
		t.Parallel()
		path := "tgtesting/testdata/bigger.mp4"
		if _, err := os.Stat(path); err != nil {
			t.Skipf("fixture not found: %v", err)
		}

		meta, err := probeISOBMFF(path, true)
		if err != nil {
			t.Fatalf("probeISOBMFF: %v", err)
		}
		if meta.tempThumbnailPath != "" {
			_ = os.Remove(meta.tempThumbnailPath)
			t.Fatalf("expected no thumbnail for an HEVC track (out of v1 scope), got %q", meta.tempThumbnailPath)
		}
	})
}

func TestIsISOBMFF(t *testing.T) {
	t.Parallel()

	if !isISOBMFF(buildBox("ftyp", []byte("isom0001"))) {
		t.Fatal("expected a valid ftyp box to be recognized as ISO-BMFF")
	}
	if isISOBMFF([]byte("not an mp4 at all")) {
		t.Fatal("expected non-ISO-BMFF content to be rejected")
	}
	if isISOBMFF([]byte("short")) {
		t.Fatal("expected a too-short header to be rejected")
	}
}

func TestBuildBoxTestHelperRoundTrip(t *testing.T) {
	t.Parallel()
	// Sanity-checks the test helper itself against readBoxes.
	box := buildBox("test", []byte{1, 2, 3})
	boxes := readBoxes(box)
	if len(boxes) != 1 || !bytes.Equal(boxes[0].payload, []byte{1, 2, 3}) {
		t.Fatalf("got %+v, want one box with payload [1 2 3]", boxes)
	}
}
