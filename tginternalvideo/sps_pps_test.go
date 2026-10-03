package tginternalvideo

import (
	"os"
	"testing"
)

// extractH264Config opens path, walks its box structure, and returns the
// parsed avcC/SPS/PPS for its video track - test-only plumbing that
// exercises the exact same path DecodeFirstFrame uses in the real pipeline.
func extractH264Config(t *testing.T, path string) (*avcDecoderConfig, *h264SPS, *h264PPS) {
	t.Helper()

	videoTrak := findVideoTrakForTest(t, path)
	avcCPayload, _, _, ok := locateH264SampleForTest(videoTrak)
	if !ok {
		t.Fatalf("%s: locateH264SampleForTest failed (not H.264, or box structure not recognized)", path)
	}
	cfg, err := parseAVCDecoderConfigurationRecord(avcCPayload)
	if err != nil {
		t.Fatalf("%s: parseAVCDecoderConfigurationRecord: %v", path, err)
	}

	spsType, spsRBSP, err := nalUnitType(cfg.sps[0])
	if err != nil || spsType != 7 {
		t.Fatalf("%s: bad SPS NAL: type=%d err=%v", path, spsType, err)
	}
	sps, err := parseSPS(spsRBSP)
	if err != nil {
		t.Fatalf("%s: parseSPS: %v", path, err)
	}

	ppsType, ppsRBSP, err := nalUnitType(cfg.pps[0])
	if err != nil || ppsType != 8 {
		t.Fatalf("%s: bad PPS NAL: type=%d err=%v", path, ppsType, err)
	}
	pps, err := parsePPS(ppsRBSP)
	if err != nil {
		t.Fatalf("%s: parsePPS: %v", path, err)
	}

	return cfg, sps, pps
}

// TestH264SPSPPSRealFixtures cross-checks parseSPS/parsePPS against ground
// truth captured via `ffmpeg -bsf:v trace_headers` on the real fixtures
// (see the plan notes): both are H.264 High profile, 4:2:0, 8-bit, CABAC,
// with the 8x8 transform enabled and flat (no custom) scaling lists -
// video.mp4 is limited-range with no cropping, bigger_2.mp4 is full-range
// with non-MB-aligned dimensions requiring cropping.
func TestH264SPSPPSRealFixtures(t *testing.T) {
	t.Parallel()

	t.Run("video.mp4", func(t *testing.T) {
		t.Parallel()
		path := "../tgtesting/testdata/video.mp4"
		if _, err := os.Stat(path); err != nil {
			t.Skipf("fixture not found: %v", err)
		}
		_, sps, pps := extractH264Config(t, path)

		if sps.profileIdc != 100 {
			t.Fatalf("profileIdc = %d, want 100 (High)", sps.profileIdc)
		}
		if sps.chromaFormatIdc != 1 {
			t.Fatalf("chromaFormatIdc = %d, want 1", sps.chromaFormatIdc)
		}
		if sps.bitDepthLumaMinus8 != 0 || sps.bitDepthChromaMinus8 != 0 {
			t.Fatalf("bit depth minus8 = (%d,%d), want (0,0)", sps.bitDepthLumaMinus8, sps.bitDepthChromaMinus8)
		}
		if sps.frameCroppingFlag {
			t.Fatal("expected no cropping for video.mp4")
		}
		if sps.videoFullRangeFlag {
			t.Fatal("expected limited range for video.mp4")
		}
		width, height := sps.displaySize()
		if width != 720 || height != 1280 {
			t.Fatalf("displaySize = (%d,%d), want (720,1280)", width, height)
		}
		if !pps.entropyCodingModeCABAC {
			t.Fatal("expected CABAC")
		}
		if !pps.transform8x8ModeFlag {
			t.Fatal("expected transform_8x8_mode_flag=1")
		}
	})

	t.Run("bigger_2.mp4", func(t *testing.T) {
		t.Parallel()
		path := "../tgtesting/testdata/bigger_2.mp4"
		if _, err := os.Stat(path); err != nil {
			t.Skipf("fixture not found: %v", err)
		}
		_, sps, pps := extractH264Config(t, path)

		if sps.profileIdc != 100 {
			t.Fatalf("profileIdc = %d, want 100 (High)", sps.profileIdc)
		}
		if !sps.frameCroppingFlag {
			t.Fatal("expected cropping for bigger_2.mp4")
		}
		if sps.cropRight != 3 || sps.cropBottom != 7 || sps.cropLeft != 0 || sps.cropTop != 0 {
			t.Fatalf("crop = (L%d R%d T%d B%d), want (L0 R3 T0 B7)", sps.cropLeft, sps.cropRight, sps.cropTop, sps.cropBottom)
		}
		if !sps.videoFullRangeFlag {
			t.Fatal("expected full range for bigger_2.mp4")
		}
		width, height := sps.displaySize()
		if width != 794 || height != 1058 {
			t.Fatalf("displaySize = (%d,%d), want (794,1058)", width, height)
		}
		if !pps.entropyCodingModeCABAC {
			t.Fatal("expected CABAC")
		}
		if !pps.transform8x8ModeFlag {
			t.Fatal("expected transform_8x8_mode_flag=1")
		}
	})

	t.Run("bigger.mp4 is HEVC, not H.264", func(t *testing.T) {
		t.Parallel()
		path := "../tgtesting/testdata/bigger.mp4"
		if _, err := os.Stat(path); err != nil {
			t.Skipf("fixture not found: %v", err)
		}
		videoTrak := findVideoTrakForTest(t, path)
		if _, _, _, ok := locateH264SampleForTest(videoTrak); ok {
			t.Fatal("expected locateH264SampleForTest to report ok=false for an HEVC track")
		}
	})
}
