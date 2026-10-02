package tg

import (
	"encoding/binary"
	"errors"
	"os"
)

var (
	errNoMoovBox         = errors.New("tg: no moov box found")
	errNoUsableVideoMeta = errors.New("tg: no usable mp4/mov metadata")
)

const (
	// maxTopLevelBoxes bounds the top-level box scan against a pathological
	// or corrupt file (each iteration is a couple of small reads, so this is
	// a cheap, generous ceiling).
	maxTopLevelBoxes = 4096
	// maxMoovSize bounds how much of the file probeISOBMFF will ever read
	// into memory at once (the moov box, which holds only metadata, never
	// the actual audio/video payload - that lives in mdat, which is never
	// read here).
	maxMoovSize = 64 << 20 // 64MB
)

// isoBox is one parsed ISO-BMFF ("box"/"atom") record: its fourcc type and
// payload (header already stripped).
type isoBox struct {
	boxType string
	payload []byte
}

// readBoxes splits an in-memory buffer of consecutive ISO-BMFF boxes into
// (type, payload) pairs, handling the size==1 64-bit extended-size case.
// Truncated/malformed trailing bytes are silently dropped - this parser is
// best-effort only, per autoFillMediaMetadata's fail-open contract.
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

// probeISOBMFF extracts width/height/duration from an MP4/MOV file without
// decoding any video frame. It scans only TOP-LEVEL boxes via Seek/ReadAt -
// crucially, mdat (which holds the actual, potentially huge, audio/video
// payload) is skipped without ever reading its contents, and moov is found
// whether it appears near the start ("fast-start" files) or trails at the
// end after mdat (common with some encoders) - the file is never read in
// full. Returns an error (never a panic) on any malformed/unsupported input.
func probeISOBMFF(path string) (*mediaMeta, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()

	var offset int64
	var moov []byte
	header := make([]byte, 16)
	for i := 0; i < maxTopLevelBoxes && offset+8 <= size; i++ {
		if _, err := file.ReadAt(header[:8], offset); err != nil {
			break
		}
		boxSize := uint64(binary.BigEndian.Uint32(header[0:4]))
		boxType := string(header[4:8])
		headerLen := int64(8)
		if boxSize == 1 {
			if _, err := file.ReadAt(header[8:16], offset+8); err != nil {
				break
			}
			boxSize = binary.BigEndian.Uint64(header[8:16])
			headerLen = 16
		}
		if boxSize == 0 {
			boxSize = uint64(size - offset)
		}
		if boxSize < uint64(headerLen) {
			break
		}

		if boxType == "moov" {
			payloadSize := int64(boxSize) - headerLen
			if payloadSize <= 0 || payloadSize > maxMoovSize {
				break
			}
			moov = make([]byte, payloadSize)
			if _, err := file.ReadAt(moov, offset+headerLen); err != nil {
				return nil, err
			}
			break
		}
		offset += int64(boxSize) // Skip ftyp/mdat/free/wide/... without reading the payload.
	}
	if moov == nil {
		return nil, errNoMoovBox
	}

	meta := &mediaMeta{}
	var videoTrak []byte
	for _, box := range readBoxes(moov) {
		switch box.boxType {
		case "mvhd":
			if duration, ok := parseMvhd(box.payload); ok {
				meta.durationSeconds = duration
			}
		case "trak":
			if videoTrak == nil && isVideoTrak(box.payload) {
				videoTrak = box.payload
			}
		}
	}
	if videoTrak != nil {
		for _, box := range readBoxes(videoTrak) {
			if box.boxType == "tkhd" {
				if width, height, ok := parseTkhd(box.payload); ok {
					meta.width, meta.height = width, height
				}
			}
		}
	}
	if meta.width == 0 && meta.height == 0 && meta.durationSeconds == 0 {
		return nil, errNoUsableVideoMeta
	}
	return meta, nil
}

// isVideoTrak descends trak -> mdia -> hdlr and reports whether its
// handler_type is 'vide'. A file has one trak (and one tkhd) per track, so
// this is needed to pick the video track's tkhd over an audio track's.
func isVideoTrak(trak []byte) bool {
	for _, box := range readBoxes(trak) {
		if box.boxType != "mdia" {
			continue
		}
		for _, mdiaBox := range readBoxes(box.payload) {
			if mdiaBox.boxType == "hdlr" && len(mdiaBox.payload) >= 12 {
				// HandlerBox (ISO/IEC 14496-12 §8.4.3): version/flags(4) +
				// pre_defined(4) + handler_type(4) + ...
				return string(mdiaBox.payload[8:12]) == "vide"
			}
		}
	}
	return false
}

// parseMvhd reads duration/timescale from a MovieHeaderBox (ISO/IEC
// 14496-12 §8.2.2), handling both box versions (version 1 uses 64-bit
// time/duration fields).
func parseMvhd(payload []byte) (durationSeconds int64, ok bool) {
	if len(payload) < 4 {
		return 0, false
	}
	version := payload[0]

	var timescaleOffset, durationOffset, durationLen int
	if version == 1 {
		timescaleOffset, durationOffset, durationLen = 20, 24, 8
	} else {
		timescaleOffset, durationOffset, durationLen = 12, 16, 4
	}
	if len(payload) < durationOffset+durationLen {
		return 0, false
	}

	timescale := binary.BigEndian.Uint32(payload[timescaleOffset : timescaleOffset+4])
	if timescale == 0 {
		return 0, false
	}

	var duration uint64
	if durationLen == 8 {
		duration = binary.BigEndian.Uint64(payload[durationOffset : durationOffset+8])
	} else {
		duration = uint64(binary.BigEndian.Uint32(payload[durationOffset : durationOffset+4]))
	}
	return int64(float64(duration) / float64(timescale)), true
}

// parseTkhd reads width/height from a TrackHeaderBox (ISO/IEC 14496-12
// §8.3.2) as 16.16 fixed-point values. They sit right after: version/
// flags(4) + {creation+modification+track_ID+reserved+duration}({20|32}
// depending on version) + reserved(8)+layer(2)+alternate_group(2)+
// volume(2)+reserved(2)+matrix(36) = {76|88} bytes in.
func parseTkhd(payload []byte) (width, height int64, ok bool) {
	if len(payload) < 4 {
		return 0, 0, false
	}
	offset := 76
	if payload[0] == 1 {
		offset = 88
	}
	if len(payload) < offset+8 {
		return 0, 0, false
	}
	w := binary.BigEndian.Uint32(payload[offset : offset+4])
	h := binary.BigEndian.Uint32(payload[offset+4 : offset+8])
	return int64(w >> 16), int64(h >> 16), true
}
