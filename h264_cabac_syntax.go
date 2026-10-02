package tg

// h264MacroblockInfo holds the subset of a decoded macroblock's state that
// later macroblocks need for CABAC context derivation (§9.3.3.1.1's
// condTermFlagN pattern: "look at my left/top neighbor's equivalent value,
// treating an unavailable neighbor as the base case").
type h264MacroblockInfo struct {
	available           bool
	isINxN              bool // mb_type == I_NxN (vs I_16x16/I_PCM)
	isIPCM              bool
	transformSize8x8    bool
	intraChromaPredMode int
	cbpChroma           int
	lumaCBF             [4]bool    // per 8x8 luma block (shared across its 4x4 sub-blocks, used for cbp context)
	lumaDCCBF           bool       // Intra16x16 luma DC block's coded_block_flag (ctxBlockCat 0)
	luma4x4CBF          [16]bool   // per 4x4 luma block (LumaLevel4x4, ctxBlockCat 2 - only meaningful for I_NxN/4x4 MBs)
	predMode4x4         [16]int    // effective per-4x4-position intra pred mode, for neighbor derivation (always populated, even for 8x8-transform MBs by replication, and for I_16x16 MBs where it's unused by callers since isINxN is false)
	chromaDCCBF         [2]bool    // per chroma component
	chromaACCBF         [2][4]bool // per chroma component, per 4x4 block
}

// zigzag4x4 is Table 8-13's frame scan: zigzag4x4[scanIdx] = raster index
// (row*4+col) within the 4x4 block.
var zigzag4x4 = [16]int{0, 1, 4, 8, 5, 2, 3, 6, 9, 12, 13, 10, 7, 11, 14, 15}

// zigzag8x8 is Table 8-14's frame scan: zigzag8x8[scanIdx] = raster index
// (row*8+col) within the 8x8 block. This is FFmpeg's ff_zigzag_direct[64]
// (libavcodec/mathtables.c) with its row/column TRANSPOSE applied
// (libavcodec/h264_slice.c: "zigzag_scan8x8[i] = TRANSPOSE(ff_zigzag_direct[i])",
// TRANSPOSE(x) = (x>>3)|((x&7)<<3)) - ff_zigzag_direct alone is in a
// column-major convention and does not trace a valid zigzag path over a
// row-major (row*8+col) grid; a first attempt at this decoder used the
// untransposed values directly, corrupting coefficient placement for every
// 8x8 block.
var zigzag8x8 = [64]int{
	0, 8, 1, 2, 9, 16, 24, 17,
	10, 3, 4, 11, 18, 25, 32, 40,
	33, 26, 19, 12, 5, 6, 13, 20,
	27, 34, 41, 48, 56, 49, 42, 35,
	28, 21, 14, 7, 15, 22, 29, 36,
	43, 50, 57, 58, 51, 44, 37, 30,
	23, 31, 38, 45, 52, 59, 60, 53,
	46, 39, 47, 54, 61, 62, 55, 63,
}

func condTermFlag(n *h264MacroblockInfo, cond func(*h264MacroblockInfo) bool) int {
	if n == nil || !n.available {
		return 0
	}
	if cond(n) {
		return 1
	}
	return 0
}

// decodeMbTypeI decodes mb_type in an I slice (Table 9-36/9-39): returns one
// of mbTypeINxN, mbTypeI16x16 (with its 24 sub-variants folded into the
// returned predMode/cbpLuma16/cbpChroma), or mbTypeIPCM.
const (
	mbTypeINxN = iota
	mbTypeI16x16
	mbTypeIPCM
)

type h264MbTypeResult struct {
	kind                int
	intra16x16PredMode  int
	intra16x16CbpChroma int
	intra16x16CbpLuma   int // 0 or 15 (all 4x4 AC blocks coded or none)
}

func decodeMbTypeI(d *cabacDecoder, left, top *h264MacroblockInfo) (*h264MbTypeResult, error) {
	condA := condTermFlag(left, func(n *h264MacroblockInfo) bool { return !n.isINxN })
	condB := condTermFlag(top, func(n *h264MacroblockInfo) bool { return !n.isINxN })
	bin0, err := d.decodeDecision(ctxMbTypeIPrefix + condA + condB)
	if err != nil {
		return nil, err
	}
	if bin0 == 0 {
		return &h264MbTypeResult{kind: mbTypeINxN}, nil
	}

	isPCM, err := d.decodeTerminate()
	if err != nil {
		return nil, err
	}
	if isPCM == 1 {
		return &h264MbTypeResult{kind: mbTypeIPCM}, nil
	}

	// Raw mb_type (1-24), computed via the exact bin-weighting FFmpeg's
	// decode_cabac_intra_mb_type uses, then looked up in iMbTypeInfo -
	// NOT a direct bitwise unpacking of the decoded bins, which a first
	// attempt at this decoder assumed (and got wrong: the pred-mode bit
	// pairs do not map to pred mode 0-3 in bit order).
	cbpLumaBit, err := d.decodeDecision(ctxMbTypeI16x16 + 0)
	if err != nil {
		return nil, err
	}
	mbType := 1 + 12*cbpLumaBit

	cbpChromaBit0, err := d.decodeDecision(ctxMbTypeI16x16 + 1)
	if err != nil {
		return nil, err
	}
	if cbpChromaBit0 == 1 {
		cbpChromaBit1, err := d.decodeDecision(ctxMbTypeI16x16 + 2)
		if err != nil {
			return nil, err
		}
		mbType += 4 + 4*cbpChromaBit1
	}

	predBitA, err := d.decodeDecision(ctxMbTypeI16x16 + 3)
	if err != nil {
		return nil, err
	}
	predBitB, err := d.decodeDecision(ctxMbTypeI16x16 + 4)
	if err != nil {
		return nil, err
	}
	mbType += 2*predBitA + predBitB

	info := iMbTypeInfo[mbType]
	return &h264MbTypeResult{
		kind:                mbTypeI16x16,
		intra16x16PredMode:  info.predMode,
		intra16x16CbpChroma: info.cbpChroma,
		intra16x16CbpLuma:   info.cbpLuma,
	}, nil
}

func decodeTransformSize8x8Flag(d *cabacDecoder, left, top *h264MacroblockInfo) (bool, error) {
	condA := condTermFlag(left, func(n *h264MacroblockInfo) bool { return n.transformSize8x8 })
	condB := condTermFlag(top, func(n *h264MacroblockInfo) bool { return n.transformSize8x8 })
	bin, err := d.decodeDecision(ctxTransformSize8x8Flag + condA + condB)
	if err != nil {
		return false, err
	}
	return bin == 1, nil
}

// decodeIntraChromaPredMode decodes intra_chroma_pred_mode (§9.3.2.4, Table
// 9-18): truncated unary, cMax=3.
func decodeIntraChromaPredMode(d *cabacDecoder, left, top *h264MacroblockInfo) (int, error) {
	condA := condTermFlag(left, func(n *h264MacroblockInfo) bool { return !n.isIPCM && n.intraChromaPredMode != 0 })
	condB := condTermFlag(top, func(n *h264MacroblockInfo) bool { return !n.isIPCM && n.intraChromaPredMode != 0 })
	first := ctxIntraChromaPredMode + condA + condB
	return d.decodeUnaryMax(3, func(binIdx int) int {
		if binIdx == 0 {
			return first
		}
		return ctxIntraChromaPredMode + 3
	})
}

func decodePrevIntraPredModeFlag(d *cabacDecoder) (bool, error) {
	bin, err := d.decodeDecision(ctxPrevIntraPredModeFlag)
	if err != nil {
		return false, err
	}
	return bin == 1, nil
}

// decodeRemIntraPredMode decodes rem_intra4x4/8x8_pred_mode: FL(3) - 3
// fixed-length bins, all using the same single context. The three bins are
// combined LSB-first (bin0 has weight 1, bin1 weight 2, bin2 weight 4) -
// verified against FFmpeg's decode_cabac_mb_intra4x4_pred_mode; a first
// attempt at this decoder combined them MSB-first, which doesn't desync the
// bitstream (still exactly 3 bins either way) but silently chose the wrong
// predicted mode for most I_NxN blocks.
func decodeRemIntraPredMode(d *cabacDecoder) (int, error) {
	v := 0
	for i := 0; i < 3; i++ {
		bin, err := d.decodeDecision(ctxRemIntraPredMode)
		if err != nil {
			return 0, err
		}
		v |= bin << i
	}
	return v, nil
}

// decodeCodedBlockPatternLuma decodes the 4 luma coded_block_pattern bits
// (one per 8x8 luma block, raster order 0=TL,1=TR,2=BL,3=BR), each using its
// own left/top 8x8-block neighbor (which may be inside the same macroblock).
func decodeCodedBlockPatternLuma(d *cabacDecoder, mb, left, top *h264MacroblockInfo) ([4]bool, error) {
	var cbp [4]bool
	leftOf := func(blk int) (*h264MacroblockInfo, int) {
		switch blk {
		case 0:
			return left, 1
		case 1:
			return mb, 0
		case 2:
			return left, 3
		default:
			return mb, 2
		}
	}
	topOf := func(blk int) (*h264MacroblockInfo, int) {
		switch blk {
		case 0:
			return top, 2
		case 1:
			return top, 3
		case 2:
			return mb, 0
		default:
			return mb, 1
		}
	}
	blockCoded := func(n *h264MacroblockInfo, idx int) bool {
		if n == nil || !n.available || n.isIPCM {
			return false
		}
		if n == mb {
			return cbp[idx]
		}
		return n.lumaCBF[idx]
	}
	for blk := 0; blk < 4; blk++ {
		ln, li := leftOf(blk)
		tn, ti := topOf(blk)
		// condTermFlagN is 1 only when the neighbor is available and its
		// referenced bit is NOT coded; an unavailable neighbor contributes 0
		// (verified against a real FFmpeg CABAC bin-trace for mb0 of
		// tgtesting/testdata/video.mp4, which has no left/top neighbor at
		// all and decodes coded_block_pattern's first 3 bins at the base
		// ctxIdx with zero contribution from either side). A prior attempt
		// inverted this (unavailable => 1), reasoning from a misremembered
		// reading of §9.3.3.1.1.4; that regressed real decode output even
		// though it matched the general default used by most OTHER CABAC
		// syntax elements (condTermFlag helper below), which do not apply
		// here - coded_block_pattern's own derivation is special-cased.
		condA := 0
		if ln != nil && ln.available && !ln.isIPCM && !blockCoded(ln, li) {
			condA = 1
		}
		condB := 0
		if tn != nil && tn.available && !tn.isIPCM && !blockCoded(tn, ti) {
			condB = 1
		}
		bin, err := d.decodeDecision(ctxCbpLuma + condA + 2*condB)
		if err != nil {
			return cbp, err
		}
		cbp[blk] = bin == 1
	}
	return cbp, nil
}

// decodeCodedBlockPatternChroma decodes coded_block_pattern's chroma part
// (0, 1, or 2): truncated unary, cMax=2. Both bins combine left AND top
// neighbor conditions (ctxIdxInc 0-3 each) - a first attempt at this
// decoder only considered the left neighbor for bin1, which under-sized
// its context range (2 slots instead of the real 4).
func decodeCodedBlockPatternChroma(d *cabacDecoder, left, top *h264MacroblockInfo) (int, error) {
	cond := func(n *h264MacroblockInfo, want func(int) bool) int {
		if n == nil || !n.available || n.isIPCM {
			return 0
		}
		if want(n.cbpChroma) {
			return 1
		}
		return 0
	}
	bin0ctx := ctxCbpChromaBin0 + cond(left, func(c int) bool { return c != 0 }) + 2*cond(top, func(c int) bool { return c != 0 })
	bin0, err := d.decodeDecision(bin0ctx)
	if err != nil {
		return 0, err
	}
	if bin0 == 0 {
		return 0, nil
	}
	bin1ctx := ctxCbpChromaBin1 + cond(left, func(c int) bool { return c == 2 }) + 2*cond(top, func(c int) bool { return c == 2 })
	bin1, err := d.decodeDecision(bin1ctx)
	if err != nil {
		return 0, err
	}
	return 1 + bin1, nil
}

// decodeMbQpDelta decodes mb_qp_delta (§9.3.2.7): truncated-unary-like with
// no fixed cMax (practically bounded well under 52), special ctxIdxInc.
func decodeMbQpDelta(d *cabacDecoder, prevMbQpDeltaNonZero bool) (int, error) {
	first := ctxMbQpDelta
	if prevMbQpDeltaNonZero {
		first++
	}
	binIdx := 0
	ctxIdxFn := func(_ int) int {
		var c int
		switch {
		case binIdx == 0:
			c = first
		case binIdx == 1:
			c = ctxMbQpDelta + 2
		default:
			c = ctxMbQpDelta + 3
		}
		binIdx++
		return c
	}
	ue, err := d.decodeUnaryMax(64, ctxIdxFn)
	if err != nil {
		return 0, err
	}
	// mb_qp_delta's code value maps unary-coded magnitude to signed delta
	// via the standard "se(v)"-like mapping: 0,1,-1,2,-2,3,-3,...
	if ue == 0 {
		return 0, nil
	}
	if ue%2 == 1 {
		return (ue + 1) / 2, nil
	}
	return -(ue / 2), nil
}

// residualBlock decodes one CABAC-coded residual block (§7.3.5.3.3,
// §9.3.3.1.3): coded_block_flag (skipped for cat5, which has none - the
// caller only calls this when it already knows the block is coded),
// significant_coeff_flag/last_significant_coeff_flag in forward scan order,
// then coeff_abs_level_minus1/coeff_sign_flag in REVERSE order from the
// last significant coefficient back to the first (per spec - this order
// matters for the numDecodAbsLevelGt1/Eq1 running counts). Returns
// coefficients in scan order (not yet de-zigzagged).
func decodeResidualBlockCABAC(d *cabacDecoder, ctxBlockCat, maxNumCoeff int) ([]int32, error) {
	coeffs := make([]int32, maxNumCoeff)

	sigBase := sigCoeffFlagCtxBase(ctxBlockCat)
	lastBase := lastSigCoeffFlagCtxBase(ctxBlockCat)
	is8x8 := ctxBlockCat == 5

	// Forward scan: find every significant position up to (and including)
	// the last one, stopping as soon as last_significant_coeff_flag reports
	// this is the final one; per spec, the final position (maxNumCoeff-1) is
	// always itself significant whether we reach it via an explicit
	// last_significant_coeff_flag==1 or simply by exhausting every earlier
	// position without ever seeing one (§7.3.5.3.3/§9.3.2.3). For 8x8 blocks
	// (cat 5), the scan-position-to-context mapping goes through
	// significantCoeffFlagOffset8x8/lastCoeffFlagOffset8x8 (far fewer
	// distinct contexts than scan positions) rather than directly indexing
	// by scan position.
	var significantPositions []int
	foundLast := false
	for i := 0; i < maxNumCoeff-1; i++ {
		sigCtx, lastCtx := sigBase+i, lastBase+i
		if is8x8 {
			sigCtx = sigBase + int(significantCoeffFlagOffset8x8[i])
			lastCtx = lastBase + int(lastCoeffFlagOffset8x8[i])
		}
		sig, err := d.decodeDecision(sigCtx)
		if err != nil {
			return nil, err
		}
		if sig == 1 {
			significantPositions = append(significantPositions, i)
			last, err := d.decodeDecision(lastCtx)
			if err != nil {
				return nil, err
			}
			if last == 1 {
				foundLast = true
				break
			}
		}
	}
	if !foundLast {
		significantPositions = append(significantPositions, maxNumCoeff-1)
	}

	// Reverse scan: coeff_abs_level_minus1 + coeff_sign_flag, from the last
	// significant position back to the first - context selection follows
	// this processing order, not forward scan order, via the node_ctx state
	// machine in coeffAbsLevel1Ctx/coeffAbsLevelGt1Ctx/coeffAbsLevelTransition
	// (see their doc comment in h264_cabac_tables.go for why this isn't a
	// simple running count).
	absLevelBase := coeffAbsLevelCtxBase(ctxBlockCat)
	nodeCtx := 0
	for i := len(significantPositions) - 1; i >= 0; i-- {
		pos := significantPositions[i]

		gt1Ctx := absLevelBase + int(coeffAbsLevel1Ctx[nodeCtx])
		bin0, err := d.decodeDecision(gt1Ctx)
		if err != nil {
			return nil, err
		}

		var absLevel int
		if bin0 == 0 {
			absLevel = 1
			nodeCtx = int(coeffAbsLevelTransition[0][nodeCtx])
		} else {
			levelCtx := absLevelBase + int(coeffAbsLevelGt1Ctx[nodeCtx])
			nodeCtx = int(coeffAbsLevelTransition[1][nodeCtx])

			coeffAbs := 2
			for coeffAbs < 15 {
				bin, err := d.decodeDecision(levelCtx)
				if err != nil {
					return nil, err
				}
				if bin == 0 {
					break
				}
				coeffAbs++
			}
			if coeffAbs >= 15 {
				// UEG0 (unsigned Exp-Golomb order 0) suffix, bypass-coded:
				// count a run of 1-bits (bypass) terminated by a 0, then
				// read that many more bypass bits as the binary suffix.
				// Final value is 14 + 2^runLength + suffix (NOT
				// 14 + (2^runLength - 1) + suffix, the standard Exp-Golomb
				// formula - this UEGk suffix, continuing from a TU-14
				// prefix, is offset by one from a from-scratch Exp-Golomb
				// code; a first attempt at this decoder used the standard
				// formula here, which is the confirmed, concrete bug this
				// rewrite fixes).
				runLength := 0
				for {
					bin, err := d.decodeBypass()
					if err != nil {
						return nil, err
					}
					if bin == 0 {
						break
					}
					runLength++
					if runLength > 32 {
						return nil, errH264Malformed
					}
				}
				suffix := 1
				for k := 0; k < runLength; k++ {
					bit, err := d.decodeBypass()
					if err != nil {
						return nil, err
					}
					suffix = suffix<<1 | bit
				}
				coeffAbs = 14 + suffix
			}
			absLevel = coeffAbs
		}

		signBit, err := d.decodeBypass()
		if err != nil {
			return nil, err
		}
		if signBit == 1 {
			coeffs[pos] = -int32(absLevel)
		} else {
			coeffs[pos] = int32(absLevel)
		}
	}

	return coeffs, nil
}

// decodeCodedBlockFlag decodes coded_block_flag for ctxBlockCat (0..4; never
// called for cat 5, which has none). An unavailable neighbor contributes 1
// (treated as if coded), not 0 - verified against a real FFmpeg CABAC
// bin-trace for mb0 of tgtesting/testdata/video.mp4 (no left/top mb at all),
// whose first luma 4x4 coded_block_flag decodes at ctxIdx 96 = base(93) + 3,
// i.e. both condTermFlags are 1 despite neither neighbor existing. This is
// the opposite default from coded_block_pattern's own derivation (see
// decodeCodedBlockPatternLuma) - each CABAC syntax element's condTermFlag
// rule for an unavailable neighbor must be checked individually, not assumed
// from another element's convention.
func decodeCodedBlockFlag(d *cabacDecoder, ctxBlockCat int, leftCoded, topCoded, leftAvail, topAvail bool) (bool, error) {
	condA := 1
	if leftAvail && !leftCoded {
		condA = 0
	}
	condB := 1
	if topAvail && !topCoded {
		condB = 0
	}
	base := codedBlockFlagCtxBase(ctxBlockCat)
	bin, err := d.decodeDecision(base + condA + 2*condB)
	if err != nil {
		return false, err
	}
	return bin == 1, nil
}
