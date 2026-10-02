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
func probeISOBMFF(path string, videoFrameDecode bool) (*mediaMeta, error) {
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
	if videoFrameDecode && videoTrak != nil {
		if thumb := tryDecodeH264Thumbnail(file, videoTrak); thumb != "" {
			meta.tempThumbnailPath = thumb
		}
	}
	return meta, nil
}

// tryDecodeH264Thumbnail best-effort decodes sample #1 of an H.264 video
// track into a real first-frame thumbnail. file must still be open and
// seekable (the caller's defer hasn't fired yet). Returns "" on literally
// anything going wrong - unsupported feature, malformed input, or even an
// internal panic in the experimental decoder - since this entire path is
// opt-in best-effort on top of the metadata-only behavior above, which must
// never be put at risk by it.
func tryDecodeH264Thumbnail(file *os.File, videoTrak []byte) (thumbPath string) {
	defer func() { _ = recover() }()

	cfg, sampleOffset, sampleSize, ok := locateH264Sample(videoTrak)
	if !ok {
		return ""
	}
	sample := make([]byte, sampleSize)
	if _, err := file.ReadAt(sample, sampleOffset); err != nil {
		return ""
	}
	img, err := decodeFirstH264Frame(cfg, sample)
	if err != nil {
		return ""
	}
	path, err := writeJPEGThumbnail(img)
	if err != nil {
		return ""
	}
	return path
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

// findStbl descends trak -> mdia -> minf -> stbl, returning the Sample
// Table Box's payload (where the sample description, sizes, and chunk
// offsets this decoder needs all live), or nil if the structure is missing.
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

// findVideoSampleEntry returns the first (and normally only) sample entry
// inside an stsd box's payload (ISO/IEC 14496-12 §8.5.2): its fourcc (e.g.
// "avc1"/"avc3" for H.264, "hvc1"/"hev1" for HEVC) and its own payload. A
// SampleEntry is itself size+fourcc+payload shaped, so readBoxes parses it
// directly once the 8-byte stsd FullBox header (version/flags+entry_count)
// is skipped.
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

// visualSampleEntryFixedSize is VisualSampleEntry's fixed-field length
// (ISO/IEC 14496-12 §12.1.3) before any child boxes (like avcC) begin:
// SampleEntry's reserved(6)+data_reference_index(2) = 8, then
// pre_defined(2)+reserved(2)+pre_defined[3](12)+width(2)+height(2)+
// horizresolution(4)+vertresolution(4)+reserved(4)+frame_count(2)+
// compressorname(32)+depth(2)+pre_defined(2) = 70, total 78.
const visualSampleEntryFixedSize = 78

// findAvcCBox locates the avcC (AVCDecoderConfigurationRecord) child box
// inside a parsed "avc1"/"avc3" sample entry's payload.
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

// firstSampleSize reads the size of sample #1 from an stsz box (ISO/IEC
// 14496-12 §8.7.3.2): either the box's fixed per-sample size, or the first
// entry of its per-sample size table.
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

// firstChunkOffset reads the first chunk's file offset from an stco (32-bit)
// or co64 (64-bit) box (ISO/IEC 14496-12 §8.7.5). Sample #1 is always the
// first sample of the first chunk, so this is also sample #1's file offset -
// no stsc (sample-to-chunk) parsing is needed just to locate it.
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

// firstSampleIsSync reports whether sample #1 is listed in an stss (Sync
// Sample Box, ISO/IEC 14496-12 §8.6.2) table. A malformed/empty table is
// treated leniently (true) - this check exists only to avoid confidently
// decoding a non-keyframe as if it were one, not to block on every
// stss edge case.
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

// locateH264Sample finds sample #1's H.264 decoder config and file location
// within a video trak, by descending into its Sample Table Box. Returns
// ok=false (never an error) for anything that doesn't fit - including
// simply not being H.264 (e.g. "hvc1"/"hev1" HEVC tracks), which is the
// expected, common case that must leave the existing metadata-only
// behavior untouched.
func locateH264Sample(videoTrak []byte) (cfg *avcDecoderConfig, sampleOffset, sampleSize int64, ok bool) {
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
	avcCfg, err := parseAVCDecoderConfigurationRecord(avcCPayload)
	if err != nil {
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

	return avcCfg, offset, size, true
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
