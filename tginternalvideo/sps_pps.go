package tginternalvideo

// h264SPS holds the subset of a parsed Sequence Parameter Set (ITU-T H.264
// §7.3.2.1.1) this best-effort decoder needs. Fields it reads only to stay
// correctly bit-aligned en route to a field it does need (e.g. the
// pic-order-cnt bookkeeping) aren't retained.
type h264SPS struct {
	profileIdc              int
	chromaFormatIdc         int
	separateColourPlaneFlag bool
	bitDepthLumaMinus8      int
	bitDepthChromaMinus8    int

	log2MaxFrameNumMinus4       int
	picOrderCntType             int
	log2MaxPicOrderCntLsbMinus4 int
	deltaPicOrderAlwaysZeroFlag bool

	picWidthInMbsMinus1       int
	picHeightInMapUnitsMinus1 int
	frameMbsOnlyFlag          bool
	direct8x8InferenceFlag    bool

	frameCroppingFlag                        bool
	cropLeft, cropRight, cropTop, cropBottom int

	videoFullRangeFlag bool // defaults false (limited/studio range) per spec
}

// picWidthInMbs is PicWidthInMbs (§7.4.2.1.1).
func (sps *h264SPS) picWidthInMbs() int { return sps.picWidthInMbsMinus1 + 1 }

// frameHeightInMbs is FrameHeightInMbs (§7.4.2.1.1); always PicHeightInMapUnits
// itself here since frameMbsOnlyFlag==1 is required (interlaced is unsupported).
func (sps *h264SPS) frameHeightInMbs() int { return sps.picHeightInMapUnitsMinus1 + 1 }

// displaySize returns the cropped (display) width/height in luma samples
// (§7.4.2.1.1's cropping rectangle, 4:2:0-only as required by v1 scope).
func (sps *h264SPS) displaySize() (width, height int) {
	width = sps.picWidthInMbs() * 16
	height = sps.frameHeightInMbs() * 16
	if !sps.frameCroppingFlag {
		return width, height
	}
	const subWidthC, subHeightC = 2, 2 // 4:2:0 only (v1 scope)
	width -= subWidthC * (sps.cropLeft + sps.cropRight)
	height -= subHeightC * (sps.cropTop + sps.cropBottom)
	return width, height
}

func isHighProfileFamily(profileIdc int) bool {
	switch profileIdc {
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		return true
	}
	return false
}

// parseSPS parses a Sequence Parameter Set NAL's RBSP. Returns
// errH264Unsupported for anything outside the v1 scope cuts (interlaced,
// non-4:2:0, non-8-bit, custom scaling lists) and errH264Malformed for a
// bitstream that doesn't parse at all - callers must fail open either way.
func parseSPS(rbsp []byte) (*h264SPS, error) {
	r := newH264BitReader(rbsp)
	sps := &h264SPS{chromaFormatIdc: 1} // default when not signaled below

	profileIdc, err := r.readBits(8)
	if err != nil {
		return nil, err
	}
	sps.profileIdc = int(profileIdc)

	if _, err := r.readBits(8); err != nil { // constraint_set0..5_flag + reserved_zero_2bits
		return nil, err
	}
	if _, err := r.readBits(8); err != nil { // level_idc
		return nil, err
	}
	if _, err := r.readUE(); err != nil { // seq_parameter_set_id
		return nil, err
	}

	if isHighProfileFamily(sps.profileIdc) {
		chromaFormatIdc, err := r.readUE()
		if err != nil {
			return nil, err
		}
		sps.chromaFormatIdc = int(chromaFormatIdc)
		if sps.chromaFormatIdc == 3 {
			sep, err := r.readFlag()
			if err != nil {
				return nil, err
			}
			sps.separateColourPlaneFlag = sep
		}
		bdLuma, err := r.readUE()
		if err != nil {
			return nil, err
		}
		sps.bitDepthLumaMinus8 = int(bdLuma)
		bdChroma, err := r.readUE()
		if err != nil {
			return nil, err
		}
		sps.bitDepthChromaMinus8 = int(bdChroma)
		if _, err := r.readFlag(); err != nil { // qpprime_y_zero_transform_bypass_flag
			return nil, err
		}
		scalingPresent, err := r.readFlag()
		if err != nil {
			return nil, err
		}
		if scalingPresent {
			// Custom scaling lists are out of v1 scope - and parsing past
			// this point would require decoding the variable-length
			// scaling_list() structures we don't implement, so bail now
			// rather than desync.
			return nil, errH264Unsupported
		}
	}
	if sps.chromaFormatIdc != 1 {
		return nil, errH264Unsupported // v1 scope: 4:2:0 only
	}
	if sps.bitDepthLumaMinus8 != 0 || sps.bitDepthChromaMinus8 != 0 {
		return nil, errH264Unsupported // v1 scope: 8-bit only
	}

	log2MaxFrameNumMinus4, err := r.readUE()
	if err != nil {
		return nil, err
	}
	sps.log2MaxFrameNumMinus4 = int(log2MaxFrameNumMinus4)

	picOrderCntType, err := r.readUE()
	if err != nil {
		return nil, err
	}
	sps.picOrderCntType = int(picOrderCntType)

	switch sps.picOrderCntType {
	case 0:
		v, err := r.readUE()
		if err != nil {
			return nil, err
		}
		sps.log2MaxPicOrderCntLsbMinus4 = int(v)
	case 1:
		deltaAlwaysZero, err := r.readFlag()
		if err != nil {
			return nil, err
		}
		sps.deltaPicOrderAlwaysZeroFlag = deltaAlwaysZero
		if _, err := r.readSE(); err != nil { // offset_for_non_ref_pic
			return nil, err
		}
		if _, err := r.readSE(); err != nil { // offset_for_top_to_bottom_field
			return nil, err
		}
		numRefFrames, err := r.readUE()
		if err != nil {
			return nil, err
		}
		for i := uint32(0); i < numRefFrames; i++ {
			if _, err := r.readSE(); err != nil { // offset_for_ref_frame[i]
				return nil, err
			}
		}
	}

	if _, err := r.readUE(); err != nil { // max_num_ref_frames
		return nil, err
	}
	if _, err := r.readFlag(); err != nil { // gaps_in_frame_num_value_allowed_flag
		return nil, err
	}

	picWidthInMbsMinus1, err := r.readUE()
	if err != nil {
		return nil, err
	}
	sps.picWidthInMbsMinus1 = int(picWidthInMbsMinus1)

	picHeightInMapUnitsMinus1, err := r.readUE()
	if err != nil {
		return nil, err
	}
	sps.picHeightInMapUnitsMinus1 = int(picHeightInMapUnitsMinus1)

	frameMbsOnly, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	sps.frameMbsOnlyFlag = frameMbsOnly
	if !sps.frameMbsOnlyFlag {
		// Interlaced/MBAFF - macroblock addressing and field decoding are
		// substantially more complex, out of v1 scope.
		return nil, errH264Unsupported
	}

	direct8x8, err := r.readFlag()
	if err != nil { // direct_8x8_inference_flag
		return nil, err
	}
	sps.direct8x8InferenceFlag = direct8x8

	frameCropping, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	sps.frameCroppingFlag = frameCropping
	if sps.frameCroppingFlag {
		left, err := r.readUE()
		if err != nil {
			return nil, err
		}
		right, err := r.readUE()
		if err != nil {
			return nil, err
		}
		top, err := r.readUE()
		if err != nil {
			return nil, err
		}
		bottom, err := r.readUE()
		if err != nil {
			return nil, err
		}
		sps.cropLeft, sps.cropRight = int(left), int(right)
		sps.cropTop, sps.cropBottom = int(top), int(bottom)
	}

	vuiPresent, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	if vuiPresent {
		if err := parseVUIForColorRange(r, sps); err != nil {
			return nil, err
		}
	}

	return sps, nil
}

// parseVUIForColorRange walks just enough of vui_parameters() (Annex E.1.1)
// to reach video_full_range_flag, which is the only VUI field this decoder
// needs (for correct YCbCr->RGB conversion) - everything else in and after
// it is read-and-discarded purely to stay bit-aligned, or simply never
// reached once we have what we need.
func parseVUIForColorRange(r *h264BitReader, sps *h264SPS) error {
	aspectRatioPresent, err := r.readFlag()
	if err != nil {
		return err
	}
	if aspectRatioPresent {
		aspectRatioIdc, err := r.readBits(8)
		if err != nil {
			return err
		}
		const extendedSAR = 255
		if aspectRatioIdc == extendedSAR {
			if _, err := r.readBits(16); err != nil { // sar_width
				return err
			}
			if _, err := r.readBits(16); err != nil { // sar_height
				return err
			}
		}
	}

	overscanPresent, err := r.readFlag()
	if err != nil {
		return err
	}
	if overscanPresent {
		if _, err := r.readFlag(); err != nil { // overscan_appropriate_flag
			return err
		}
	}

	videoSignalPresent, err := r.readFlag()
	if err != nil {
		return err
	}
	if videoSignalPresent {
		if _, err := r.readBits(3); err != nil { // video_format
			return err
		}
		fullRange, err := r.readFlag()
		if err != nil {
			return err
		}
		sps.videoFullRangeFlag = fullRange
	}
	return nil
}

// h264PPS holds the subset of a parsed Picture Parameter Set (ITU-T H.264
// §7.3.2.2) this decoder needs.
type h264PPS struct {
	spsId                                 int
	entropyCodingModeCABAC                bool
	bottomFieldPicOrderInFramePresentFlag bool
	picInitQpMinus26                      int
	chromaQpIndexOffset                   int
	secondChromaQpIndexOffset             int
	deblockingFilterControlPresentFlag    bool
	constrainedIntraPredFlag              bool
	redundantPicCntPresentFlag            bool
	transform8x8ModeFlag                  bool
}

// parsePPS parses a Picture Parameter Set NAL's RBSP. Same fail-open error
// contract as parseSPS.
func parsePPS(rbsp []byte) (*h264PPS, error) {
	r := newH264BitReader(rbsp)
	pps := &h264PPS{}

	if _, err := r.readUE(); err != nil { // pic_parameter_set_id
		return nil, err
	}
	spsId, err := r.readUE()
	if err != nil {
		return nil, err
	}
	pps.spsId = int(spsId)

	entropyMode, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	pps.entropyCodingModeCABAC = entropyMode
	if !pps.entropyCodingModeCABAC {
		return nil, errH264Unsupported // v1 scope: CABAC only, not CAVLC
	}

	bottomFieldPicOrderPresent, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	pps.bottomFieldPicOrderInFramePresentFlag = bottomFieldPicOrderPresent

	numSliceGroupsMinus1, err := r.readUE()
	if err != nil {
		return nil, err
	}
	if numSliceGroupsMinus1 > 0 {
		return nil, errH264Unsupported // FMO / slice groups, out of scope
	}

	if _, err := r.readUE(); err != nil { // num_ref_idx_l0_default_active_minus1
		return nil, err
	}
	if _, err := r.readUE(); err != nil { // num_ref_idx_l1_default_active_minus1
		return nil, err
	}
	if _, err := r.readFlag(); err != nil { // weighted_pred_flag
		return nil, err
	}
	if _, err := r.readBits(2); err != nil { // weighted_bipred_idc
		return nil, err
	}

	picInitQpMinus26, err := r.readSE()
	if err != nil {
		return nil, err
	}
	pps.picInitQpMinus26 = int(picInitQpMinus26)

	if _, err := r.readSE(); err != nil { // pic_init_qs_minus26
		return nil, err
	}

	chromaQpIndexOffset, err := r.readSE()
	if err != nil {
		return nil, err
	}
	pps.chromaQpIndexOffset = int(chromaQpIndexOffset)
	pps.secondChromaQpIndexOffset = pps.chromaQpIndexOffset // default, per spec, unless overridden below

	deblockCtrl, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	pps.deblockingFilterControlPresentFlag = deblockCtrl

	constrainedIntra, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	pps.constrainedIntraPredFlag = constrainedIntra

	redundant, err := r.readFlag()
	if err != nil {
		return nil, err
	}
	pps.redundantPicCntPresentFlag = redundant

	if r.moreRBSPData() {
		transform8x8, err := r.readFlag()
		if err != nil {
			return nil, err
		}
		pps.transform8x8ModeFlag = transform8x8

		scalingPresent, err := r.readFlag()
		if err != nil {
			return nil, err
		}
		if scalingPresent {
			return nil, errH264Unsupported // v1 scope cut, same as SPS
		}

		secondChromaQp, err := r.readSE()
		if err != nil {
			return nil, err
		}
		pps.secondChromaQpIndexOffset = int(secondChromaQp)
	}

	return pps, nil
}
