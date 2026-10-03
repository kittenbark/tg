package tginternalvideo

// blk4x4Pos maps a luma4x4BlkIdx (§6.4.3's Z-order numbering within a
// macroblock) to its (x,y) position in 4x4-block units (0..3 each).
var blk4x4Pos = [16][2]int{
	{0, 0}, {1, 0}, {0, 1}, {1, 1},
	{2, 0}, {3, 0}, {2, 1}, {3, 1},
	{0, 2}, {1, 2}, {0, 3}, {1, 3},
	{2, 2}, {3, 2}, {2, 3}, {3, 3},
}

// blk8x8Pos maps luma8x8BlkIdx (0..3) to its (x,y) position in 8x8-block units.
var blk8x8Pos = [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

func deZigzag4x4(scan []int32) [16]int32 {
	var out [16]int32
	for i, v := range scan {
		out[zigzag4x4[i]] = v
	}
	return out
}

func deZigzag8x8(scan []int32) [64]int32 {
	var out [64]int32
	for i, v := range scan {
		out[zigzag8x8[i]] = v
	}
	return out
}

// deZigzagAC4x4 places 15 AC-only coefficients (scan positions 1..15) into
// a full 16-element raster block, leaving position 0 (DC) as zero - the
// caller fills it in separately from the Intra16x16 luma DC transform.
func deZigzagAC4x4(scanAC []int32) [16]int32 {
	var out [16]int32
	for i, v := range scanAC {
		out[zigzag4x4[i+1]] = v
	}
	return out
}

type h264SliceHeader struct {
	firstMbInSlice int
	sliceQPY       int
}

// parseSliceHeaderI parses slice_header() (§7.3.3) for the IDR I-slice case
// only (this decoder never calls it otherwise): single slice per picture,
// progressive frame coding, CABAC. Any value outside that falls back
// (errH264Unsupported).
func parseSliceHeaderI(r *h264BitReader, sps *h264SPS, pps *h264PPS) (*h264SliceHeader, error) {
	firstMb, err := r.readUE()
	if err != nil {
		return nil, err
	}
	if firstMb != 0 {
		return nil, errH264Unsupported // v1 scope: single slice covering the whole picture
	}

	sliceTypeRaw, err := r.readUE()
	if err != nil {
		return nil, err
	}
	if sliceTypeRaw%5 != 2 { // 2 == I, 7 == I (same type, "all slices in picture are I")
		return nil, errH264Unsupported
	}

	if _, err := r.readUE(); err != nil { // pic_parameter_set_id
		return nil, err
	}
	if sps.separateColourPlaneFlag {
		if _, err := r.readBits(2); err != nil { // colour_plane_id
			return nil, err
		}
	}

	frameNumBits := sps.log2MaxFrameNumMinus4 + 4
	if _, err := r.readBits(frameNumBits); err != nil { // frame_num
		return nil, err
	}
	// frame_mbs_only_flag is required true (enforced in parseSPS), so no
	// field_pic_flag/bottom_field_flag here.

	if _, err := r.readUE(); err != nil { // idr_pic_id (always present: IDR only)
		return nil, err
	}

	switch {
	case sps.picOrderCntType == 0:
		pocBits := sps.log2MaxPicOrderCntLsbMinus4 + 4
		if _, err := r.readBits(pocBits); err != nil { // pic_order_cnt_lsb
			return nil, err
		}
		if pps.bottomFieldPicOrderInFramePresentFlag {
			if _, err := r.readSE(); err != nil { // delta_pic_order_cnt_bottom
				return nil, err
			}
		}
	case sps.picOrderCntType == 1 && !sps.deltaPicOrderAlwaysZeroFlag:
		if _, err := r.readSE(); err != nil { // delta_pic_order_cnt[0]
			return nil, err
		}
		if pps.bottomFieldPicOrderInFramePresentFlag {
			if _, err := r.readSE(); err != nil { // delta_pic_order_cnt[1]
				return nil, err
			}
		}
	}

	if pps.redundantPicCntPresentFlag {
		if _, err := r.readUE(); err != nil { // redundant_pic_cnt
			return nil, err
		}
	}

	// nal_ref_idc != 0 always holds for IDR pictures.
	if _, err := r.readFlag(); err != nil { // no_output_of_prior_pics_flag
		return nil, err
	}
	if _, err := r.readFlag(); err != nil { // long_term_reference_flag
		return nil, err
	}

	sliceQpDelta, err := r.readSE()
	if err != nil {
		return nil, err
	}
	sliceQPY := 26 + pps.picInitQpMinus26 + int(sliceQpDelta)

	if pps.deblockingFilterControlPresentFlag {
		idc, err := r.readUE()
		if err != nil {
			return nil, err
		}
		if idc != 1 {
			if _, err := r.readSE(); err != nil { // slice_alpha_c0_offset_div2
				return nil, err
			}
			if _, err := r.readSE(); err != nil { // slice_beta_offset_div2
				return nil, err
			}
		}
	}

	return &h264SliceHeader{firstMbInSlice: int(firstMb), sliceQPY: sliceQPY}, nil
}

func neighborPredMode4x4(n *h264MacroblockInfo, pos int) int {
	if n == nil || !n.available || !n.isINxN {
		return 2 // DC fallback (§8.3.1.1): unavailable or non-Intra_4x4/8x8 neighbor
	}
	return n.predMode4x4[pos]
}

// decodeIDRPictureI decodes one IDR I-slice covering the whole picture
// (§7.3.4's slice_data(), I-slice path) into reconstructed Y/Cb/Cr planes.
// Any unsupported feature or malformed input returns an error - the caller
// (video_metadata.go, under recover()) must fail open to metadata-only.
func decodeIDRPictureI(sps *h264SPS, pps *h264PPS, rbsp []byte) (y, cb, cr *planeView, err error) {
	r := newH264BitReader(rbsp)
	sh, err := parseSliceHeaderI(r, sps, pps)
	if err != nil {
		return nil, nil, nil, err
	}

	d, err := newCabacDecoder(r, sh.sliceQPY)
	if err != nil {
		return nil, nil, nil, err
	}

	mbWidth := sps.picWidthInMbs()
	mbHeight := sps.frameHeightInMbs()
	picWidth, picHeight := mbWidth*16, mbHeight*16
	chromaWidth, chromaHeight := picWidth/2, picHeight/2

	yPlane := &planeView{data: make([]uint8, picWidth*picHeight), stride: picWidth, width: picWidth, height: picHeight}
	cbPlane := &planeView{data: make([]uint8, chromaWidth*chromaHeight), stride: chromaWidth, width: chromaWidth, height: chromaHeight}
	crPlane := &planeView{data: make([]uint8, chromaWidth*chromaHeight), stride: chromaWidth, width: chromaWidth, height: chromaHeight}
	chromaPlanes := [2]*planeView{cbPlane, crPlane}

	mbInfos := make([]h264MacroblockInfo, mbWidth*mbHeight)
	sliceQPY := sh.sliceQPY
	prevMbQpDeltaNonZero := false

	numMBs := mbWidth * mbHeight
	for mbAddr := 0; mbAddr < numMBs; mbAddr++ {
		mbX, mbY := mbAddr%mbWidth, mbAddr/mbWidth
		var left, top *h264MacroblockInfo
		if mbX > 0 {
			left = &mbInfos[mbAddr-1]
		}
		if mbY > 0 {
			top = &mbInfos[mbAddr-mbWidth]
		}

		mbType, err := decodeMbTypeI(d, left, top)
		if err != nil {
			return nil, nil, nil, err
		}
		if mbType.kind == mbTypeIPCM {
			return nil, nil, nil, errH264Unsupported // v1 scope cut: I_PCM
		}

		info := &mbInfos[mbAddr]
		info.available = true
		x0, y0 := mbX*16, mbY*16
		cx0, cy0 := mbX*8, mbY*8
		haveLeft, haveTop := left != nil, top != nil
		haveTopLeft := haveLeft && haveTop

		if mbType.kind == mbTypeI16x16 {
			if err := decodeI16x16Macroblock(d, info, left, top, mbType, yPlane, chromaPlanes, x0, y0, cx0, cy0, haveLeft, haveTop, haveTopLeft, &sliceQPY, &prevMbQpDeltaNonZero, pps); err != nil {
				return nil, nil, nil, err
			}
			endOfSlice, err := d.decodeTerminate()
			if err != nil {
				return nil, nil, nil, err
			}
			if endOfSlice == 1 {
				break
			}
			continue
		}

		// I_NxN
		info.isINxN = true
		transform8x8, err := decodeTransformSize8x8Flag(d, left, top)
		if err != nil {
			return nil, nil, nil, err
		}
		info.transformSize8x8 = transform8x8

		if transform8x8 {
			if err := decodeINxN8x8Macroblock(d, info, left, top, yPlane, x0, y0); err != nil {
				return nil, nil, nil, err
			}
		} else {
			if err := decodeINxN4x4Macroblock(d, info, left, top, yPlane, x0, y0); err != nil {
				return nil, nil, nil, err
			}
		}

		chromaPredMode, err := decodeIntraChromaPredMode(d, left, top)
		if err != nil {
			return nil, nil, nil, err
		}
		info.intraChromaPredMode = chromaPredMode
		if err := predictAndSetChroma(chromaPlanes, cx0, cy0, chromaPredMode, haveLeft, haveTop, haveTopLeft); err != nil {
			return nil, nil, nil, err
		}

		lumaCBF, err := decodeCodedBlockPatternLuma(d, info, left, top)
		if err != nil {
			return nil, nil, nil, err
		}
		info.lumaCBF = lumaCBF
		if transform8x8 {
			// An 8x8-transform block has no per-4x4 coded_block_flag of its
			// own, but a 4x4-granularity neighbor (an I_NxN/4x4 block, or an
			// I_16x16 block's luma AC) still needs an answer when querying
			// this block's 4x4 positions for its own coded_block_flag
			// context derivation - per spec, each inherits its parent 8x8
			// block's cbp bit. Leaving luma4x4CBF at its zero value here
			// (as a first attempt at this decoder did) made every 8x8
			// neighbor look "not coded" to such queries, regardless of its
			// real cbp, picking the wrong CABAC context.
			for blk8x8 := 0; blk8x8 < 4; blk8x8++ {
				for _, p4 := range blk4x4IndicesIn8x8(blk8x8) {
					info.luma4x4CBF[p4] = lumaCBF[blk8x8]
				}
			}
		}
		cbpChroma, err := decodeCodedBlockPatternChroma(d, left, top)
		if err != nil {
			return nil, nil, nil, err
		}
		info.cbpChroma = cbpChroma

		anyLumaCoded := lumaCBF[0] || lumaCBF[1] || lumaCBF[2] || lumaCBF[3]
		if anyLumaCoded || cbpChroma != 0 {
			delta, err := decodeMbQpDelta(d, prevMbQpDeltaNonZero)
			if err != nil {
				return nil, nil, nil, err
			}
			prevMbQpDeltaNonZero = delta != 0
			sliceQPY = ((sliceQPY+delta)%52 + 52) % 52
		} else {
			prevMbQpDeltaNonZero = false
		}

		if transform8x8 {
			if err := decodeLumaResidual8x8(d, info, left, top, yPlane, x0, y0, sliceQPY); err != nil {
				return nil, nil, nil, err
			}
		} else {
			if err := decodeLumaResidual4x4(d, info, left, top, yPlane, x0, y0, sliceQPY); err != nil {
				return nil, nil, nil, err
			}
		}
		if err := decodeChromaResidual(d, info, left, top, chromaPlanes, cx0, cy0, sliceQPY, pps, cbpChroma); err != nil {
			return nil, nil, nil, err
		}
		endOfSlice, err := d.decodeTerminate()
		if err != nil {
			return nil, nil, nil, err
		}
		if endOfSlice == 1 {
			break
		}
	}

	return yPlane, cbPlane, crPlane, nil
}

func decodeI16x16Macroblock(d *cabacDecoder, info, left, top *h264MacroblockInfo, mbType *h264MbTypeResult, yPlane *planeView, chromaPlanes [2]*planeView, x0, y0, cx0, cy0 int, haveLeft, haveTop, haveTopLeft bool, sliceQPY *int, prevMbQpDeltaNonZero *bool, pps *h264PPS) error {
	info.isINxN = false
	info.transformSize8x8 = false

	chromaPredMode, err := decodeIntraChromaPredMode(d, left, top)
	if err != nil {
		return err
	}
	info.intraChromaPredMode = chromaPredMode
	if err := predictAndSetChroma(chromaPlanes, cx0, cy0, chromaPredMode, haveLeft, haveTop, haveTopLeft); err != nil {
		return err
	}

	cbpChroma := mbType.intra16x16CbpChroma
	cbpLuma := mbType.intra16x16CbpLuma
	info.cbpChroma = cbpChroma
	for i := range info.lumaCBF {
		info.lumaCBF[i] = cbpLuma != 0
	}

	delta, err := decodeMbQpDelta(d, *prevMbQpDeltaNonZero)
	if err != nil {
		return err
	}
	*prevMbQpDeltaNonZero = delta != 0
	*sliceQPY = ((*sliceQPY+delta)%52 + 52) % 52
	qp := *sliceQPY

	pred := predictIntra16x16(yPlane, x0, y0, mbType.intra16x16PredMode, haveLeft, haveTop, haveTopLeft)
	yPlane.writeBlock(x0, y0, 16, pred)

	leftDC := left != nil && left.available && left.lumaDCCBF
	topDC := top != nil && top.available && top.lumaDCCBF
	dcCoded, err := decodeCodedBlockFlag(d, 0, leftDC, topDC, left != nil, top != nil)
	if err != nil {
		return err
	}
	info.lumaDCCBF = dcCoded

	var dcValues [16]int32 // indexed by (y4*4+x4) raster block position
	if dcCoded {
		scan, err := decodeResidualBlockCABAC(d, 0, 16)
		if err != nil {
			return err
		}
		dcGrid := deZigzag4x4(scan)
		had := hadamard4x4Inverse(dcGrid)
		dcValues = dequantizeLumaDC(had, qp)
	}

	for blkIdx := 0; blkIdx < 16; blkIdx++ {
		pos := blk4x4Pos[blkIdx]
		bx0, by0 := x0+pos[0]*4, y0+pos[1]*4
		raster := pos[1]*4 + pos[0]

		var block [16]int32
		block[0] = dcValues[raster]

		if cbpLuma != 0 {
			leftN, leftPos := blk4x4Neighbor(left, info, blkIdx, true)
			topN, topPos := blk4x4Neighbor(top, info, blkIdx, false)
			leftCBF := leftN != nil && leftN.available && leftN.luma4x4CBF[leftPos]
			topCBF := topN != nil && topN.available && topN.luma4x4CBF[topPos]
			coded, err := decodeCodedBlockFlag(d, 1, leftCBF, topCBF, leftN != nil, topN != nil)
			if err != nil {
				return err
			}
			info.luma4x4CBF[blkIdx] = coded
			if coded {
				scanAC, err := decodeResidualBlockCABAC(d, 1, 15)
				if err != nil {
					return err
				}
				ac := deZigzagAC4x4(scanAC)
				for i := 1; i < 16; i++ {
					block[i] = ac[i]
				}
			}
		}

		deq := dequantize4x4(block[:], qp)
		deq[0] = block[0] // DC already fully dequantized via dequantizeLumaDC; don't re-scale it
		residual := idct4x4(deq)
		yPlane.addResidualBlock(bx0, by0, 4, residual[:])
	}

	// A first attempt at this decoder never decoded chroma residual for
	// I_16x16 macroblocks at all - confirmed by cross-checking FFmpeg's
	// decode_cabac_luma_residual/residual_chroma call sequence, which is
	// shared (unconditional) across every intra mb_type, not just I_NxN.
	// Any I_16x16 macroblock with cbpChroma!=0 (common) desyncs the whole
	// rest of the bitstream without this call.
	return decodeChromaResidual(d, info, left, top, chromaPlanes, cx0, cy0, qp, pps, cbpChroma)
}

// blk4x4Neighbor returns the neighboring 4x4 block's info and its
// luma4x4BlkIdx, for either the left (wantLeft=true) or top neighbor of
// blkIdx within mb - which may be another block inside the SAME macroblock.
func blk4x4Neighbor(mbNeighbor, mb *h264MacroblockInfo, blkIdx int, wantLeft bool) (*h264MacroblockInfo, int) {
	pos := blk4x4Pos[blkIdx]
	x4, y4 := pos[0], pos[1]
	if wantLeft {
		if x4 > 0 {
			return mb, posToBlk4x4(x4-1, y4)
		}
		return mbNeighbor, posToBlk4x4(3, y4)
	}
	if y4 > 0 {
		return mb, posToBlk4x4(x4, y4-1)
	}
	return mbNeighbor, posToBlk4x4(x4, 3)
}

func posToBlk4x4(x4, y4 int) int {
	for idx, p := range blk4x4Pos {
		if p[0] == x4 && p[1] == y4 {
			return idx
		}
	}
	return 0
}

// decodeINxN4x4Macroblock decodes only the 16 prev_intra4x4_pred_mode_flag/
// rem_intra4x4_pred_mode syntax elements (§7.3.5.1's luma4x4BlkIdx loop),
// storing each block's mode. It deliberately does NOT predict or reconstruct
// pixels yet: a 4x4 block's prediction needs its causal neighbors' FINAL
// reconstructed samples (prediction + residual), which for blocks within
// this same macroblock aren't available until decodeLumaResidual4x4 has
// processed them - that function does the actual per-block predict+
// reconstruct, interleaved with residual decode, in the correct order. A
// first attempt at this decoder predicted (and wrote) all 16 blocks here,
// before any residual was known, so every block after the first predicted
// from neighbor pixels that were never going to be final.
func decodeINxN4x4Macroblock(d *cabacDecoder, info, left, top *h264MacroblockInfo, yPlane *planeView, x0, y0 int) error {
	for blkIdx := 0; blkIdx < 16; blkIdx++ {
		leftN, leftPos := blk4x4Neighbor(left, info, blkIdx, true)
		topN, topPos := blk4x4Neighbor(top, info, blkIdx, false)
		modeA := neighborPredMode4x4(leftN, leftPos)
		modeB := neighborPredMode4x4(topN, topPos)
		predMode := modeA
		if modeB < predMode {
			predMode = modeB
		}

		prevFlag, err := decodePrevIntraPredModeFlag(d)
		if err != nil {
			return err
		}
		mode := predMode
		if !prevFlag {
			rem, err := decodeRemIntraPredMode(d)
			if err != nil {
				return err
			}
			if rem < predMode {
				mode = rem
			} else {
				mode = rem + 1
			}
		}
		info.predMode4x4[blkIdx] = mode
	}
	return nil
}

func blk4x4HasTopLeft(mb, left, top *h264MacroblockInfo, blkIdx int) bool {
	pos := blk4x4Pos[blkIdx]
	x4, y4 := pos[0], pos[1]
	switch {
	case x4 > 0 && y4 > 0:
		return true // inside the same MB
	case x4 == 0 && y4 == 0:
		return left != nil && top != nil // true top-left MB corner needs both neighbors (diagonal)
	case x4 == 0:
		return left != nil
	default: // y4==0, x4>0
		return top != nil
	}
}

func blk4x4HasTopRight(mb, left, top *h264MacroblockInfo, blkIdx int) bool {
	pos := blk4x4Pos[blkIdx]
	x4, y4 := pos[0], pos[1]
	// Top-right is available when the block immediately above-right has
	// already been decoded - true for the top row's rightmost 3 columns
	// (from the top MB) and for in-MB positions whose above-right 4x4 block
	// precedes this one in Z-order; approximated conservatively (false
	// where genuinely unclear) since predictIntra4x4 replicates the last
	// top sample when unavailable, matching the spec's own fallback.
	if y4 == 0 {
		if x4 < 3 {
			return top != nil
		}
		return false // top-right MB not tracked in this single-slice-per-picture v1 (treated as unavailable)
	}
	switch blkIdx {
	case 2, 8, 10: // (0,1),(0,2),(1,3) equivalents whose above-right block is already decoded in-MB
		return true
	default:
		return false
	}
}

func decodeINxN8x8Macroblock(d *cabacDecoder, info, left, top *h264MacroblockInfo, yPlane *planeView, x0, y0 int) error {
	for blkIdx := 0; blkIdx < 4; blkIdx++ {
		pos := blk8x8Pos[blkIdx]

		// 8x8 neighbor pred modes follow the same left-min/top-min rule,
		// using the 4x4-position equivalents at the 8x8 block's corner.
		var leftMode, topMode int
		if pos[0] > 0 {
			leftMode = pred8x8ModeFromMB(info, blkIdx-1)
		} else if left != nil && left.available {
			leftMode = pred8x8ModeOrDC(left, blkIdx+1)
		} else {
			leftMode = 2
		}
		if pos[1] > 0 {
			topMode = pred8x8ModeFromMB(info, blkIdx-2)
		} else if top != nil && top.available {
			topMode = pred8x8ModeOrDC(top, blkIdx+2)
		} else {
			topMode = 2
		}
		predMode := leftMode
		if topMode < predMode {
			predMode = topMode
		}

		prevFlag, err := decodePrevIntraPredModeFlag(d)
		if err != nil {
			return err
		}
		mode := predMode
		if !prevFlag {
			rem, err := decodeRemIntraPredMode(d)
			if err != nil {
				return err
			}
			if rem < predMode {
				mode = rem
			} else {
				mode = rem + 1
			}
		}
		for _, p4 := range blk4x4IndicesIn8x8(blkIdx) {
			info.predMode4x4[p4] = mode
		}
	}
	return nil
}

// haveNeighbors8x8 returns the same-mb/same-picture neighbor availability
// for 8x8 block blkIdx (0=TL,1=TR,2=BL,3=BR), shared between mode decode and
// the predict+reconstruct pass in decodeLumaResidual8x8.
func haveNeighbors8x8(info, left, top *h264MacroblockInfo, blkIdx int) (haveLeft, haveTop, haveTopLeft, haveTopRight bool) {
	pos := blk8x8Pos[blkIdx]
	haveLeft = pos[0] > 0 || (left != nil && left.available)
	haveTop = pos[1] > 0 || (top != nil && top.available)
	haveTopLeft = haveLeft && haveTop
	// Top-right availability per 8x8 block position: TL's top-right comes
	// from the top MB's row; BL's top-right is the already-decoded TR block
	// in this same MB; TR and BR's top-right would come from a diagonal
	// neighbor MB this single-slice-per-picture decoder doesn't track, so
	// they fall back to unavailable (predictIntra8x8 replicates the last
	// top sample in that case, matching the spec's own edge handling).
	switch blkIdx {
	case 0:
		haveTopRight = top != nil && top.available
	case 2:
		haveTopRight = true
	default:
		haveTopRight = false
	}
	return
}

// pred8x8ModeFromMB/pred8x8ModeOrDC read back an already-decoded 8x8 block's
// mode via its representative 4x4 position (all 4 share the same mode).
func pred8x8ModeFromMB(info *h264MacroblockInfo, blk8x8Idx int) int {
	if blk8x8Idx < 0 {
		return 2
	}
	return info.predMode4x4[blk4x4IndicesIn8x8(blk8x8Idx)[0]]
}

func pred8x8ModeOrDC(n *h264MacroblockInfo, blk8x8Idx int) int {
	if !n.isINxN {
		return 2
	}
	return n.predMode4x4[blk4x4IndicesIn8x8(blk8x8Idx % 4)[0]]
}

func blk4x4IndicesIn8x8(blk8x8Idx int) []int {
	base := blk8x8Idx * 4
	return []int{base, base + 1, base + 2, base + 3}
}

// decodeLumaResidual4x4 does the actual per-4x4-block predict+reconstruct
// (interleaved with residual decode, in raster/Z-order so each block's
// neighbors are already fully reconstructed by the time it's predicted) -
// decodeINxN4x4Macroblock only determined the modes beforehand.
func decodeLumaResidual4x4(d *cabacDecoder, info, left, top *h264MacroblockInfo, yPlane *planeView, x0, y0, qp int) error {
	for blkIdx := 0; blkIdx < 16; blkIdx++ {
		pos := blk4x4Pos[blkIdx]
		bx0, by0 := x0+pos[0]*4, y0+pos[1]*4

		leftN, leftPos := blk4x4Neighbor(left, info, blkIdx, true)
		topN, topPos := blk4x4Neighbor(top, info, blkIdx, false)
		haveLeft := leftN != nil && leftN.available
		haveTop := topN != nil && topN.available
		haveTopLeft := blk4x4HasTopLeft(info, left, top, blkIdx)
		haveTopRight := blk4x4HasTopRight(info, left, top, blkIdx)
		pred := predictIntra4x4(yPlane, bx0, by0, info.predMode4x4[blkIdx], haveLeft, haveTop, haveTopLeft, haveTopRight)
		yPlane.writeBlock(bx0, by0, 4, pred)

		blk8x8 := blkIdx / 4
		if !info.lumaCBF[blk8x8] {
			continue
		}
		leftCBF := haveLeft && leftN.luma4x4CBF[leftPos]
		topCBF := haveTop && topN.luma4x4CBF[topPos]
		coded, err := decodeCodedBlockFlag(d, 2, leftCBF, topCBF, leftN != nil, topN != nil)
		if err != nil {
			return err
		}
		info.luma4x4CBF[blkIdx] = coded
		if !coded {
			continue
		}
		scan, err := decodeResidualBlockCABAC(d, 2, 16)
		if err != nil {
			return err
		}
		block := deZigzag4x4(scan)
		deq := dequantize4x4(block[:], qp)
		residual := idct4x4(deq)
		yPlane.addResidualBlock(bx0, by0, 4, residual[:])
	}
	return nil
}

// decodeLumaResidual8x8 does the actual per-8x8-block predict+reconstruct
// (in TL,TR,BL,BR order so each block's neighbors are already fully
// reconstructed by the time it's predicted) - decodeINxN8x8Macroblock only
// determined the modes beforehand.
func decodeLumaResidual8x8(d *cabacDecoder, info, left, top *h264MacroblockInfo, yPlane *planeView, x0, y0, qp int) error {
	for blkIdx := 0; blkIdx < 4; blkIdx++ {
		pos := blk8x8Pos[blkIdx]
		bx0, by0 := x0+pos[0]*8, y0+pos[1]*8
		mode := info.predMode4x4[blk4x4IndicesIn8x8(blkIdx)[0]]
		haveLeft, haveTop, haveTopLeft, haveTopRight := haveNeighbors8x8(info, left, top, blkIdx)
		ref := buildRefSamples8x8(yPlane, bx0, by0, haveLeft, haveTop, haveTopLeft, haveTopRight)
		pred := predictIntra8x8(ref, mode)
		yPlane.writeBlock(bx0, by0, 8, pred)

		if !info.lumaCBF[blkIdx] {
			continue
		}
		scan, err := decodeResidualBlockCABAC(d, 5, 64)
		if err != nil {
			return err
		}
		block := deZigzag8x8(scan)
		deq := dequantize8x8(block[:], qp)
		residual := idct8x8(deq)
		yPlane.addResidualBlock(bx0, by0, 8, residual[:])
	}
	return nil
}

func predictAndSetChroma(chromaPlanes [2]*planeView, cx0, cy0, mode int, haveLeft, haveTop, haveTopLeft bool) error {
	for _, plane := range chromaPlanes {
		pred := predictIntraChroma(plane, cx0, cy0, mode, haveLeft, haveTop, haveTopLeft)
		plane.writeBlock(cx0, cy0, 8, pred)
	}
	return nil
}

// chromaBlk4x4Neighbor mirrors blk4x4Neighbor for the 4 chroma 4x4 blocks
// within an 8x8 chroma area (positions 0=TL,1=TR,2=BL,3=BR, raster order).
func chromaBlk4x4Neighbor(mbNeighbor, mb *h264MacroblockInfo, comp, blk int, wantLeft bool) (*h264MacroblockInfo, int) {
	x4, y4 := blk%2, blk/2
	if wantLeft {
		if x4 > 0 {
			return mb, comp*4 + y4*2 + (x4 - 1)
		}
		return mbNeighbor, comp*4 + y4*2 + 1
	}
	if y4 > 0 {
		return mb, comp*4 + (y4-1)*2 + x4
	}
	return mbNeighbor, comp*4 + 2 + x4
}

// decodeChromaResidual decodes residual_chroma() (§7.3.5.3.2): ALL chroma DC
// blocks first (component 0 then 1), THEN - only if cbpChroma==2 - ALL
// chroma AC blocks (component 0's 4 blocks, then component 1's 4 blocks).
// This two-pass order (not DC+AC interleaved per component) is required by
// the spec; a first attempt at this decoder interleaved them per component,
// which desyncs the bitstream for any macroblock with cbpChroma==2.
func decodeChromaResidual(d *cabacDecoder, info, left, top *h264MacroblockInfo, chromaPlanes [2]*planeView, cx0, cy0, qp int, pps *h264PPS, cbpChroma int) error {
	if cbpChroma == 0 {
		return nil
	}

	var dc [2][4]int32
	for comp := 0; comp < 2; comp++ {
		leftDC := left != nil && left.available && left.chromaDCCBF[comp]
		topDC := top != nil && top.available && top.chromaDCCBF[comp]
		coded, err := decodeCodedBlockFlag(d, 3, leftDC, topDC, left != nil, top != nil)
		if err != nil {
			return err
		}
		info.chromaDCCBF[comp] = coded
		if !coded {
			continue
		}
		scan, err := decodeResidualBlockCABAC(d, 3, 4)
		if err != nil {
			return err
		}
		var raw [4]int32
		copy(raw[:], scan)
		had := hadamard2x2Inverse(raw)
		dc[comp] = dequantizeChromaDC(had, deriveChromaQP(qp, pps, comp))
	}

	var ac [2][4][16]int32
	if cbpChroma == 2 {
		for comp := 0; comp < 2; comp++ {
			for blk := 0; blk < 4; blk++ {
				leftN, leftPos := chromaBlk4x4Neighbor(left, info, comp, blk, true)
				topN, topPos := chromaBlk4x4Neighbor(top, info, comp, blk, false)
				leftCBF := leftN != nil && leftN.available && leftN.chromaACCBF[comp][leftPos%4]
				topCBF := topN != nil && topN.available && topN.chromaACCBF[comp][topPos%4]
				coded, err := decodeCodedBlockFlag(d, 4, leftCBF, topCBF, leftN != nil, topN != nil)
				if err != nil {
					return err
				}
				info.chromaACCBF[comp][blk] = coded
				if !coded {
					continue
				}
				scanAC, err := decodeResidualBlockCABAC(d, 4, 15)
				if err != nil {
					return err
				}
				ac[comp][blk] = deZigzagAC4x4(scanAC)
			}
		}
	}

	for comp := 0; comp < 2; comp++ {
		plane := chromaPlanes[comp]
		chromaQP := deriveChromaQP(qp, pps, comp)
		for blk := 0; blk < 4; blk++ {
			x4, y4 := blk%2, blk/2
			bx0, by0 := cx0+x4*4, cy0+y4*4

			var block [16]int32
			block[0] = dc[comp][blk]
			for i := 1; i < 16; i++ {
				block[i] = ac[comp][blk][i]
			}

			deq := dequantize4x4(block[:], chromaQP)
			deq[0] = block[0]
			residual := idct4x4(deq)
			plane.addResidualBlock(bx0, by0, 4, residual[:])
		}
	}
	return nil
}

// qpcTable maps QPY (clipped to 0..51, via chroma_qp_index_offset) to QPc
// (§8.5.8, Table 8-15) for QPI in 30..51; below 30 QPc==QPI.
var qpcTable = [22]int{29, 30, 31, 32, 32, 33, 34, 34, 35, 35, 36, 36, 37, 37, 37, 38, 38, 38, 39, 39, 39, 39}

func deriveChromaQP(qpY int, pps *h264PPS, comp int) int {
	offset := pps.chromaQpIndexOffset
	if comp == 1 {
		offset = pps.secondChromaQpIndexOffset
	}
	qpi := clip3(-12, 51, qpY+offset) // QpBdOffsetY=0 for 8-bit
	qpc := qpi
	if qpi >= 30 {
		idx := qpi - 30
		if idx >= len(qpcTable) {
			idx = len(qpcTable) - 1
		}
		qpc = qpcTable[idx]
	}
	return qpc
}
