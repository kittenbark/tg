package tginternalvideo

// cabacContext is one context variable's adaptive state (ITU-T H.264 §9.3.1.2).
type cabacContext struct {
	pStateIdx uint8
	valMPS    uint8
}

// cabacDecoder is the CABAC arithmetic decoding engine (§9.3.3.2). It owns
// the bit reader for the remainder of the slice once engine initialization
// has consumed the leading cabac_alignment_one_bit padding and the initial
// 9-bit codIOffset.
type cabacDecoder struct {
	r          *h264BitReader
	codIRange  uint32
	codIOffset uint32
	contexts   []cabacContext
}

func clip3(low, high, v int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

// newCabacDecoder initializes context variables for sliceQPY (§9.3.1.1) and
// the arithmetic decoding engine (§9.3.1.2). r must already be positioned at
// the start of cabac_alignment_one_bit (i.e. right after the slice header).
func newCabacDecoder(r *h264BitReader, sliceQPY int) (*cabacDecoder, error) {
	d := &cabacDecoder{
		r:        r,
		contexts: make([]cabacContext, len(cabacContextInit)),
	}
	d.initContexts(sliceQPY)
	if err := d.initEngine(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *cabacDecoder) initContexts(sliceQPY int) {
	qp := clip3(0, 51, sliceQPY)
	for i, init := range cabacContextInit {
		preCtxState := clip3(1, 126, ((int(init.m)*qp)>>4)+int(init.n))
		if preCtxState <= 63 {
			d.contexts[i] = cabacContext{pStateIdx: uint8(63 - preCtxState), valMPS: 0}
		} else {
			d.contexts[i] = cabacContext{pStateIdx: uint8(preCtxState - 64), valMPS: 1}
		}
	}
}

func (d *cabacDecoder) initEngine() error {
	d.r.align() // cabac_alignment_one_bit padding, already validated upstream
	bits, err := d.r.readBits(9)
	if err != nil {
		return err
	}
	d.codIRange = 510
	d.codIOffset = bits
	return nil
}

// decodeDecision decodes one context-coded bin (§9.3.3.2.1/.2).
func (d *cabacDecoder) decodeDecision(ctxIdx int) (int, error) {
	ctx := &d.contexts[ctxIdx]
	qCodIRangeIdx := (d.codIRange >> 6) & 3
	codIRangeLPS := uint32(cabacRangeTabLPS[ctx.pStateIdx][qCodIRangeIdx])
	d.codIRange -= codIRangeLPS

	var binVal int
	if d.codIOffset >= d.codIRange {
		binVal = 1 - int(ctx.valMPS)
		d.codIOffset -= d.codIRange
		d.codIRange = codIRangeLPS
		if ctx.pStateIdx == 0 {
			ctx.valMPS = 1 - ctx.valMPS
		}
		ctx.pStateIdx = cabacTransIdxLPS[ctx.pStateIdx]
	} else {
		binVal = int(ctx.valMPS)
		ctx.pStateIdx = cabacTransIdxMPS[ctx.pStateIdx]
	}

	if err := d.renormalize(); err != nil {
		return 0, err
	}
	return binVal, nil
}

func (d *cabacDecoder) renormalize() error {
	for d.codIRange < 256 {
		d.codIRange <<= 1
		bit, err := d.r.readBit()
		if err != nil {
			return err
		}
		d.codIOffset = (d.codIOffset << 1) | bit
	}
	return nil
}

// decodeBypass decodes one bypass-coded bin, with no context/adaptation (§9.3.3.2.3).
func (d *cabacDecoder) decodeBypass() (int, error) {
	bit, err := d.r.readBit()
	if err != nil {
		return 0, err
	}
	d.codIOffset = (d.codIOffset << 1) | bit
	out := 0
	if d.codIOffset >= d.codIRange {
		d.codIOffset -= d.codIRange
		out = 1
	}
	return out, nil
}

// decodeBypassBits decodes n bypass bins MSB-first into an unsigned value.
func (d *cabacDecoder) decodeBypassBits(n int) (uint32, error) {
	var v uint32
	for i := 0; i < n; i++ {
		b, err := d.decodeBypass()
		if err != nil {
			return 0, err
		}
		v = (v << 1) | uint32(b)
	}
	return v, nil
}

// decodeTerminate decodes end_of_slice_flag / mb_type's I_PCM indicator (§9.3.3.2.4).
func (d *cabacDecoder) decodeTerminate() (int, error) {
	d.codIRange -= 2
	if d.codIOffset >= d.codIRange {
		return 1, nil
	}
	if err := d.renormalize(); err != nil {
		return 0, err
	}
	return 0, nil
}

// decodeUnaryMax decodes a truncated-unary value up to cMax, where ctxIdxFn
// gives the context index to use for each bin index (so callers can express
// per-bin ctxIdxInc rules, including "same context for every bin").
func (d *cabacDecoder) decodeUnaryMax(cMax int, ctxIdxFn func(binIdx int) int) (int, error) {
	for i := 0; i < cMax; i++ {
		bin, err := d.decodeDecision(ctxIdxFn(i))
		if err != nil {
			return 0, err
		}
		if bin == 0 {
			return i, nil
		}
	}
	return cMax, nil
}
