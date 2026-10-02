package tg

// avcDecoderConfig is a parsed AVCDecoderConfigurationRecord (ISO/IEC
// 14496-15 §5.2.4.1.1) - the avcC box payload stored alongside an H.264
// ("avc1"/"avc3") video sample entry in an MP4's stsd box. It carries the
// NAL length-field size used to delimit samples, plus the SPS/PPS NAL units
// needed before any slice data can be decoded.
type avcDecoderConfig struct {
	lengthSize int
	sps        [][]byte
	pps        [][]byte
}

// parseAVCDecoderConfigurationRecord parses an avcC box's payload. Returns
// errH264Malformed/errH264Unsupported (never panics) on anything that
// doesn't fit the expected shape - callers must fail open.
func parseAVCDecoderConfigurationRecord(data []byte) (*avcDecoderConfig, error) {
	if len(data) < 6 {
		return nil, errH264Malformed
	}
	if data[0] != 1 { // configurationVersion
		return nil, errH264Unsupported
	}

	cfg := &avcDecoderConfig{
		lengthSize: int(data[4]&0x03) + 1,
	}

	offset := 5
	numSPS := int(data[offset] & 0x1f)
	offset++
	for i := 0; i < numSPS; i++ {
		nal, next, err := readAVCCLengthPrefixedNAL(data, offset)
		if err != nil {
			return nil, err
		}
		cfg.sps = append(cfg.sps, nal)
		offset = next
	}

	if offset >= len(data) {
		return nil, errH264Malformed
	}
	numPPS := int(data[offset])
	offset++
	for i := 0; i < numPPS; i++ {
		nal, next, err := readAVCCLengthPrefixedNAL(data, offset)
		if err != nil {
			return nil, err
		}
		cfg.pps = append(cfg.pps, nal)
		offset = next
	}

	if len(cfg.sps) == 0 || len(cfg.pps) == 0 {
		return nil, errH264Malformed
	}
	return cfg, nil
}

// readAVCCLengthPrefixedNAL reads one (uint16 length, NAL bytes) record
// (used for the avcC box's own SPS/PPS list, which is always a 2-byte
// length regardless of the record's lengthSizeMinusOne - that field only
// governs the length prefix used in the actual sample data).
func readAVCCLengthPrefixedNAL(data []byte, offset int) (nal []byte, next int, err error) {
	if offset+2 > len(data) {
		return nil, 0, errH264Malformed
	}
	length := int(data[offset])<<8 | int(data[offset+1])
	offset += 2
	if offset+length > len(data) {
		return nil, 0, errH264Malformed
	}
	return data[offset : offset+length], offset + length, nil
}
