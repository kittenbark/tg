package tg

import "errors"

// errH264Unsupported marks a feature this best-effort decoder deliberately
// doesn't implement (see the v1 scope cuts in media_autofill.go's design
// notes) - callers must treat it exactly like any other decode failure:
// fall back to metadata-only, never surface it.
var errH264Unsupported = errors.New("tg: h264: unsupported feature")

// errH264Malformed marks a bitstream that doesn't parse as valid H.264 -
// same fail-open handling as errH264Unsupported.
var errH264Malformed = errors.New("tg: h264: malformed bitstream")

// h264BitReader is a MSB-first bit reader over a NAL unit's RBSP bytes
// (emulation-prevention bytes already removed by unescapeRBSP).
type h264BitReader struct {
	data []byte
	pos  int // bit position from the start of data
}

func newH264BitReader(data []byte) *h264BitReader {
	return &h264BitReader{data: data}
}

func (r *h264BitReader) bitsLeft() int {
	return len(r.data)*8 - r.pos
}

func (r *h264BitReader) readBit() (uint32, error) {
	if r.pos >= len(r.data)*8 {
		return 0, errH264Malformed
	}
	byteIdx := r.pos / 8
	bitIdx := 7 - uint(r.pos%8)
	bit := (r.data[byteIdx] >> bitIdx) & 1
	r.pos++
	return uint32(bit), nil
}

func (r *h264BitReader) readBits(n int) (uint32, error) {
	var v uint32
	for i := 0; i < n; i++ {
		b, err := r.readBit()
		if err != nil {
			return 0, err
		}
		v = (v << 1) | b
	}
	return v, nil
}

func (r *h264BitReader) readFlag() (bool, error) {
	b, err := r.readBit()
	if err != nil {
		return false, err
	}
	return b == 1, nil
}

// readUE reads an unsigned Exp-Golomb code (ITU-T H.264 §9.1).
func (r *h264BitReader) readUE() (uint32, error) {
	leadingZeroBits := 0
	for {
		b, err := r.readBit()
		if err != nil {
			return 0, err
		}
		if b != 0 {
			break
		}
		leadingZeroBits++
		if leadingZeroBits > 31 {
			return 0, errH264Malformed
		}
	}
	if leadingZeroBits == 0 {
		return 0, nil
	}
	rest, err := r.readBits(leadingZeroBits)
	if err != nil {
		return 0, err
	}
	return (uint32(1) << uint(leadingZeroBits)) - 1 + rest, nil
}

// readSE reads a signed Exp-Golomb code (ITU-T H.264 §9.1.1).
func (r *h264BitReader) readSE() (int32, error) {
	ue, err := r.readUE()
	if err != nil {
		return 0, err
	}
	if ue%2 == 0 {
		return -int32(ue / 2), nil
	}
	return int32((ue + 1) / 2), nil
}

func (r *h264BitReader) byteAligned() bool {
	return r.pos%8 == 0
}

func (r *h264BitReader) align() {
	r.pos = (r.pos + 7) / 8 * 8
}

// moreRBSPData reports whether more non-trailing data remains (ITU-T H.264
// §7.2's more_rbsp_data()): the RBSP ends with a single '1' stop bit
// followed by zero padding, so this is true iff the current position is
// strictly before that stop bit.
func (r *h264BitReader) moreRBSPData() bool {
	total := len(r.data) * 8
	if r.pos >= total {
		return false
	}
	lastOneBitPos := -1
	for i := total - 1; i >= r.pos; i-- {
		byteIdx := i / 8
		bitIdx := 7 - uint(i%8)
		if (r.data[byteIdx]>>bitIdx)&1 == 1 {
			lastOneBitPos = i
			break
		}
	}
	if lastOneBitPos < 0 {
		return false
	}
	return r.pos < lastOneBitPos
}

// unescapeRBSP strips emulation-prevention bytes from a NAL unit's payload
// (ITU-T H.264 §7.3.1 / §7.4.1.1): any "00 00 03" sequence has the 0x03
// removed before the bits are interpreted as RBSP. This applies regardless
// of whether the NAL was originally delimited by Annex-B start codes or
// AVCC length prefixes - it's intrinsic to how the NAL payload is encoded.
func unescapeRBSP(nal []byte) []byte {
	out := make([]byte, 0, len(nal))
	zeroRun := 0
	for i := 0; i < len(nal); i++ {
		b := nal[i]
		if zeroRun >= 2 && b == 0x03 {
			zeroRun = 0
			continue
		}
		out = append(out, b)
		if b == 0 {
			zeroRun++
		} else {
			zeroRun = 0
		}
	}
	return out
}

// splitNALUnits demuxes AVCC length-prefixed NAL units (as stored in an MP4
// sample, per ISO/IEC 14496-15) using lengthSize bytes (1, 2, or 4, from the
// avcC box) per length prefix. Malformed trailing data is dropped.
func splitNALUnits(data []byte, lengthSize int) [][]byte {
	var nals [][]byte
	for len(data) >= lengthSize {
		var length int
		for i := 0; i < lengthSize; i++ {
			length = length<<8 | int(data[i])
		}
		data = data[lengthSize:]
		if length < 0 || length > len(data) {
			return nals
		}
		nals = append(nals, data[:length])
		data = data[length:]
	}
	return nals
}

// nalUnitType returns a NAL unit's nal_unit_type (the low 5 bits of its
// first byte) and its RBSP payload (header byte stripped, emulation
// prevention removed).
func nalUnitType(nal []byte) (naluType int, rbsp []byte, err error) {
	if len(nal) < 1 {
		return 0, nil, errH264Malformed
	}
	return int(nal[0] & 0x1f), unescapeRBSP(nal[1:]), nil
}
