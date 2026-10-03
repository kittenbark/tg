package tginternalvideo

import (
	"encoding/binary"
	"os"
	"testing"
)

// This file duplicates just enough of kittenbark/tg's ISO-BMFF box-walking
// (video_metadata.go) to pull a real avcC payload + first-sample bytes out
// of the real .mp4 fixtures in ../tgtesting/testdata for this package's own
// tests. tg itself owns the real, shipped box-walking implementation; this
// decoder package doesn't depend on tg (that dependency runs the other way:
// tg imports this package), so test-only plumbing lives here instead of
// being shared.

type isoBox struct {
	boxType string
	payload []byte
}

func readBoxes(buf []byte) []isoBox {
	var boxes []isoBox
	for len(buf) >= 8 {
		size := uint64(binary.BigEndian.Uint32(buf[0:4]))
		boxType := string(buf[4:8])
		header := 8
		switch size {
		case 0:
			size = uint64(len(buf))
		case 1:
			if len(buf) < 16 {
				return boxes
			}
			size = binary.BigEndian.Uint64(buf[8:16])
			header = 16
		}
		if size < uint64(header) || uint64(len(buf)) < size {
			return boxes
		}
		boxes = append(boxes, isoBox{boxType: boxType, payload: buf[header:size]})
		buf = buf[size:]
	}
	return boxes
}

func isVideoTrak(trak []byte) bool {
	for _, box := range readBoxes(trak) {
		if box.boxType != "mdia" {
			continue
		}
		for _, mdiaBox := range readBoxes(box.payload) {
			if mdiaBox.boxType == "hdlr" && len(mdiaBox.payload) >= 12 {
				return string(mdiaBox.payload[8:12]) == "vide"
			}
		}
	}
	return false
}

func findStbl(videoTrak []byte) []byte {
	for _, trakBox := range readBoxes(videoTrak) {
		if trakBox.boxType != "mdia" {
			continue
		}
		for _, mdiaBox := range readBoxes(trakBox.payload) {
			if mdiaBox.boxType != "minf" {
				continue
			}
			for _, minfBox := range readBoxes(mdiaBox.payload) {
				if minfBox.boxType == "stbl" {
					return minfBox.payload
				}
			}
		}
	}
	return nil
}

func findVideoSampleEntry(stsdPayload []byte) (fourcc string, entryPayload []byte, ok bool) {
	if len(stsdPayload) < 8 {
		return "", nil, false
	}
	entries := readBoxes(stsdPayload[8:])
	if len(entries) == 0 {
		return "", nil, false
	}
	return entries[0].boxType, entries[0].payload, true
}

const visualSampleEntryFixedSize = 78

func findAvcCBox(sampleEntryPayload []byte) ([]byte, bool) {
	if len(sampleEntryPayload) <= visualSampleEntryFixedSize {
		return nil, false
	}
	for _, box := range readBoxes(sampleEntryPayload[visualSampleEntryFixedSize:]) {
		if box.boxType == "avcC" {
			return box.payload, true
		}
	}
	return nil, false
}

func firstSampleSize(stszPayload []byte) (int64, bool) {
	if len(stszPayload) < 12 {
		return 0, false
	}
	sampleSize := binary.BigEndian.Uint32(stszPayload[4:8])
	if sampleSize != 0 {
		return int64(sampleSize), true
	}
	sampleCount := binary.BigEndian.Uint32(stszPayload[8:12])
	if sampleCount == 0 || len(stszPayload) < 16 {
		return 0, false
	}
	return int64(binary.BigEndian.Uint32(stszPayload[12:16])), true
}

func firstChunkOffset(box isoBox) (int64, bool) {
	p := box.payload
	if len(p) < 12 {
		return 0, false
	}
	if box.boxType == "co64" {
		if len(p) < 16 {
			return 0, false
		}
		return int64(binary.BigEndian.Uint64(p[8:16])), true
	}
	return int64(binary.BigEndian.Uint32(p[8:12])), true
}

func firstSampleIsSync(stssPayload []byte) bool {
	if len(stssPayload) < 8 {
		return true
	}
	entryCount := binary.BigEndian.Uint32(stssPayload[4:8])
	if entryCount == 0 || len(stssPayload) < 12 {
		return true
	}
	return binary.BigEndian.Uint32(stssPayload[8:12]) == 1
}

// locateH264SampleForTest mirrors tg's locateH264Sample, but returns the raw
// avcC payload bytes rather than a parsed config, matching this package's
// own DecodeFirstFrame's (avcCPayload, sampleData) signature.
func locateH264SampleForTest(videoTrak []byte) (avcCPayload []byte, sampleOffset, sampleSize int64, ok bool) {
	stbl := findStbl(videoTrak)
	if stbl == nil {
		return nil, 0, 0, false
	}

	var stsdPayload, stszPayload, stssPayload []byte
	var chunkOffsetBox *isoBox
	for _, box := range readBoxes(stbl) {
		switch box.boxType {
		case "stsd":
			stsdPayload = box.payload
		case "stsz":
			stszPayload = box.payload
		case "stco", "co64":
			b := box
			chunkOffsetBox = &b
		case "stss":
			stssPayload = box.payload
		}
	}
	if stsdPayload == nil || stszPayload == nil || chunkOffsetBox == nil {
		return nil, 0, 0, false
	}
	if stssPayload != nil && !firstSampleIsSync(stssPayload) {
		return nil, 0, 0, false
	}

	fourcc, entryPayload, ok1 := findVideoSampleEntry(stsdPayload)
	if !ok1 || (fourcc != "avc1" && fourcc != "avc3") {
		return nil, 0, 0, false
	}
	avcCPayload, ok2 := findAvcCBox(entryPayload)
	if !ok2 {
		return nil, 0, 0, false
	}

	size, ok3 := firstSampleSize(stszPayload)
	if !ok3 || size <= 0 {
		return nil, 0, 0, false
	}
	offset, ok4 := firstChunkOffset(*chunkOffsetBox)
	if !ok4 {
		return nil, 0, 0, false
	}

	return avcCPayload, offset, size, true
}

// readMoovForTest duplicates probeISOBMFF's top-level box scan just enough
// to hand back the moov payload for test plumbing above.
func readMoovForTest(t *testing.T, path string) []byte {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	data := make([]byte, stat.Size())
	if _, err := file.ReadAt(data, 0); err != nil {
		t.Fatalf("read: %v", err)
	}
	for _, box := range readBoxes(data) {
		if box.boxType == "moov" {
			return box.payload
		}
	}
	t.Fatalf("%s: no moov box", path)
	return nil
}

// findVideoTrakForTest walks a real fixture's top-level boxes down to its
// video trak, for test plumbing above.
func findVideoTrakForTest(t *testing.T, path string) []byte {
	t.Helper()
	moov := readMoovForTest(t, path)
	for _, box := range readBoxes(moov) {
		if box.boxType == "trak" && isVideoTrak(box.payload) {
			return box.payload
		}
	}
	t.Fatalf("%s: no video trak found", path)
	return nil
}
