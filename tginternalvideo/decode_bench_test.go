package tginternalvideo

import (
	"os"
	"testing"
)

// BenchmarkDecodeFirstFrame isolates DecodeFirstFrame's own cost (CABAC +
// intra + transform + color convert), separate from the surrounding box-walk
// and JPEG-encode/resize steps measured by tg's own BenchmarkVideoThumbnail.
func BenchmarkDecodeFirstFrame(b *testing.B) {
	path := "../tgtesting/testdata/video.mp4"
	if _, err := os.Stat(path); err != nil {
		b.Skipf("fixture not found: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		b.Fatalf("stat: %v", err)
	}
	data := make([]byte, stat.Size())
	if _, err := file.ReadAt(data, 0); err != nil {
		b.Fatalf("read: %v", err)
	}
	var moov []byte
	for _, box := range readBoxes(data) {
		if box.boxType == "moov" {
			moov = box.payload
			break
		}
	}
	if moov == nil {
		b.Fatalf("%s: no moov box", path)
	}
	var videoTrak []byte
	for _, box := range readBoxes(moov) {
		if box.boxType == "trak" && isVideoTrak(box.payload) {
			videoTrak = box.payload
			break
		}
	}
	if videoTrak == nil {
		b.Fatalf("%s: no video trak found", path)
	}

	avcCPayload, offset, size, ok := locateH264SampleForTest(videoTrak)
	if !ok {
		b.Fatalf("locateH264SampleForTest failed")
	}
	sample := make([]byte, size)
	if _, err := file.ReadAt(sample, offset); err != nil {
		b.Fatalf("read sample: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := DecodeFirstFrame(avcCPayload, sample); err != nil {
			b.Fatalf("DecodeFirstFrame: %v", err)
		}
	}
}
