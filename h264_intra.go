package tg

// planeView is a read/write accessor over one reconstructed picture plane
// (Y, Cb, or Cr), used by intra prediction and reconstruction. Coordinates
// are absolute plane coordinates; width/height bound availability checks
// (this decoder only ever has picture-edge unavailability to worry about -
// no slice boundaries, since v1 scope requires a single slice per picture).
type planeView struct {
	data          []uint8
	stride        int
	width, height int
}

func (p *planeView) at(x, y int) int {
	if x < 0 || y < 0 || x >= p.width || y >= p.height {
		return 0
	}
	return int(p.data[y*p.stride+x])
}

func (p *planeView) set(x, y, v int) {
	if x < 0 || y < 0 || x >= p.width || y >= p.height {
		return
	}
	p.data[y*p.stride+x] = uint8(clip3(0, 255, v))
}

func (p *planeView) available(x, y int) bool {
	return x >= 0 && y >= 0 && x < p.width && y < p.height
}

// addResidualBlock adds a decoded residual block (row-major, size x size,
// already inverse-transformed) to the plane at (x0,y0), clipping to [0,255].
func (p *planeView) addResidualBlock(x0, y0, size int, residual []int32) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			p.set(x0+x, y0+y, p.at(x0+x, y0+y)+int(residual[y*size+x]))
		}
	}
}

func (p *planeView) writeBlock(x0, y0, size int, pred []int) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			p.set(x0+x, y0+y, pred[y*size+x])
		}
	}
}

// predictIntra4x4 implements the 9 Intra_4x4 prediction modes (§8.3.1.2).
// haveLeft/haveTop/haveTopRight describe neighbor availability (picture
// edges only, in this decoder); unavailable-but-selected modes fall back to
// a neutral value rather than reading garbage, since a conformant encoder
// should never select them in that situation.
func predictIntra4x4(p *planeView, x0, y0, mode int, haveLeft, haveTop, haveTopLeft, haveTopRight bool) []int {
	top := func(i int) int { // p[i,-1], i in -1..7
		if i == -1 {
			if haveTopLeft {
				return p.at(x0-1, y0-1)
			}
			return 128
		}
		if i <= 3 {
			if haveTop {
				return p.at(x0+i, y0-1)
			}
			return 128
		}
		if haveTopRight {
			return p.at(x0+i, y0-1)
		}
		if haveTop {
			return p.at(x0+3, y0-1) // replicate last available top sample
		}
		return 128
	}
	left := func(i int) int { // p[-1,i], i in -1..3
		if i == -1 {
			if haveTopLeft {
				return p.at(x0-1, y0-1)
			}
			return 128
		}
		if haveLeft {
			return p.at(x0-1, y0+i)
		}
		return 128
	}

	out := make([]int, 16)
	set := func(x, y, v int) { out[y*4+x] = v }

	switch mode {
	case 0: // Vertical
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				set(x, y, top(x))
			}
		}
	case 1: // Horizontal
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				set(x, y, left(y))
			}
		}
	case 2: // DC
		var sum, n int
		if haveTop {
			sum += top(0) + top(1) + top(2) + top(3)
			n += 4
		}
		if haveLeft {
			sum += left(0) + left(1) + left(2) + left(3)
			n += 4
		}
		dc := 128
		if n > 0 {
			dc = (sum + n/2) / n
		}
		for i := range out {
			out[i] = dc
		}
	case 3: // Diagonal_Down_Left
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				if x == 3 && y == 3 {
					set(x, y, (top(6)+3*top(7)+2)>>2)
				} else {
					set(x, y, (top(x+y)+2*top(x+y+1)+top(x+y+2)+2)>>2)
				}
			}
		}
	case 4: // Diagonal_Down_Right
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				switch {
				case x > y:
					set(x, y, (top(x-y-2)+2*top(x-y-1)+top(x-y)+2)>>2)
				case x < y:
					set(x, y, (left(y-x-2)+2*left(y-x-1)+left(y-x)+2)>>2)
				default:
					set(x, y, (top(0)+2*top(-1)+left(0)+2)>>2)
				}
			}
		}
	case 5: // Vertical_Right
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				zVR := 2*x - y
				switch {
				case zVR >= 0 && zVR%2 == 0:
					set(x, y, (top(x-(y>>1)-1)+top(x-(y>>1))+1)>>1)
				case zVR >= 0:
					set(x, y, (top(x-(y>>1)-2)+2*top(x-(y>>1)-1)+top(x-(y>>1))+2)>>2)
				case zVR == -1:
					set(x, y, (left(0)+2*top(-1)+top(0)+2)>>2)
				default:
					set(x, y, (left(y-1)+2*left(y-2)+left(y-3)+2)>>2)
				}
			}
		}
	case 6: // Horizontal_Down
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				zHD := 2*y - x
				switch {
				case zHD >= 0 && zHD%2 == 0:
					set(x, y, (left(y-(x>>1)-1)+left(y-(x>>1))+1)>>1)
				case zHD >= 0:
					set(x, y, (left(y-(x>>1)-2)+2*left(y-(x>>1)-1)+left(y-(x>>1))+2)>>2)
				case zHD == -1:
					set(x, y, (left(0)+2*top(-1)+top(0)+2)>>2)
				default:
					set(x, y, (top(x-1)+2*top(x-2)+top(x-3)+2)>>2)
				}
			}
		}
	case 7: // Vertical_Left
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				if y%2 == 0 {
					set(x, y, (top(x+(y>>1))+top(x+(y>>1)+1)+1)>>1)
				} else {
					set(x, y, (top(x+(y>>1))+2*top(x+(y>>1)+1)+top(x+(y>>1)+2)+2)>>2)
				}
			}
		}
	case 8: // Horizontal_Up
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				zHU := x + 2*y
				switch {
				case zHU < 5 && zHU%2 == 0:
					set(x, y, (left(y+(x>>1))+left(y+(x>>1)+1)+1)>>1)
				case zHU < 5:
					set(x, y, (left(y+(x>>1))+2*left(y+(x>>1)+1)+left(y+(x>>1)+2)+2)>>2)
				case zHU == 5:
					set(x, y, (left(2)+3*left(3)+2)>>2)
				default:
					set(x, y, left(3))
				}
			}
		}
	}
	return out
}

// filterRefSamples8x8 applies the Intra_8x8 low-pass reference-sample filter
// (§8.3.2.2.1) to raw neighbor samples before prediction.
type refSamples8x8 struct {
	top                                          [17]int // filtered p'[x,-1], x = -1..15 mapped to top[0..16] (top[0]=p'[-1,-1])
	left                                         [8]int  // filtered p'[-1,y], y = 0..7
	haveLeft, haveTop, haveTopLeft, haveTopRight bool
}

func buildRefSamples8x8(p *planeView, x0, y0 int, haveLeft, haveTop, haveTopLeft, haveTopRight bool) refSamples8x8 {
	raw := func(x, y int) int {
		if x == -1 && y == -1 {
			if haveTopLeft {
				return p.at(x0-1, y0-1)
			}
			if haveTop {
				return p.at(x0, y0-1)
			}
			if haveLeft {
				return p.at(x0-1, y0)
			}
			return 128
		}
		if y == -1 {
			if x <= 7 {
				if haveTop {
					return p.at(x0+x, y0-1)
				}
				return 128
			}
			if haveTopRight {
				return p.at(x0+x, y0-1)
			}
			if haveTop {
				return p.at(x0+7, y0-1)
			}
			return 128
		}
		if haveLeft {
			return p.at(x0-1, y0+y)
		}
		return 128
	}

	var r refSamples8x8
	r.haveLeft, r.haveTop, r.haveTopLeft, r.haveTopRight = haveLeft, haveTop, haveTopLeft, haveTopRight

	// Filtered top-left.
	topLeft := raw(-1, -1)
	if haveTopLeft && haveTop && haveLeft {
		topLeft = (raw(-1, 0) + 2*raw(-1, -1) + raw(0, -1) + 2) >> 2
	}
	r.top[0] = topLeft

	// Filtered top row, x=0..15 (top[1..16]).
	for x := 0; x <= 15; x++ {
		var v int
		switch {
		case x == 0:
			v = (raw(-1, -1) + 2*raw(0, -1) + raw(1, -1) + 2) >> 2
		case x == 15:
			v = (raw(14, -1) + 3*raw(15, -1) + 2) >> 2
		default:
			v = (raw(x-1, -1) + 2*raw(x, -1) + raw(x+1, -1) + 2) >> 2
		}
		r.top[x+1] = v
	}

	// Filtered left column, y=0..7.
	for y := 0; y <= 7; y++ {
		var v int
		switch {
		case y == 0:
			v = (raw(-1, -1) + 2*raw(-1, 0) + raw(-1, 1) + 2) >> 2
		case y == 7:
			v = (raw(-1, 6) + 3*raw(-1, 7) + 2) >> 2
		default:
			v = (raw(-1, y-1) + 2*raw(-1, y) + raw(-1, y+1) + 2) >> 2
		}
		r.left[y] = v
	}
	return r
}

// predictIntra8x8 implements the 9 Intra_8x8 modes (§8.3.2.2) using the
// filtered reference samples from buildRefSamples8x8. Mode numbering and
// relative formulas mirror predictIntra4x4, scaled to an 8x8 block.
func predictIntra8x8(r refSamples8x8, mode int) []int {
	top := func(i int) int { return r.top[i+1] } // i in -1..14
	left := func(i int) int {
		if i == -1 {
			return r.top[0]
		}
		return r.left[i]
	}

	out := make([]int, 64)
	set := func(x, y, v int) { out[y*8+x] = v }

	switch mode {
	case 0: // Vertical
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				set(x, y, top(x))
			}
		}
	case 1: // Horizontal
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				set(x, y, left(y))
			}
		}
	case 2: // DC
		var sum, n int
		if r.haveTop {
			for x := 0; x < 8; x++ {
				sum += top(x)
			}
			n += 8
		}
		if r.haveLeft {
			for y := 0; y < 8; y++ {
				sum += left(y)
			}
			n += 8
		}
		dc := 128
		if n > 0 {
			dc = (sum + n/2) / n
		}
		for i := range out {
			out[i] = dc
		}
	case 3: // Diagonal_Down_Left
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if x == 7 && y == 7 {
					set(x, y, (top(13)+3*top(14)+2)>>2)
				} else {
					set(x, y, (top(x+y)+2*top(x+y+1)+top(x+y+2)+2)>>2)
				}
			}
		}
	case 4: // Diagonal_Down_Right
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				switch {
				case x > y:
					set(x, y, (top(x-y-2)+2*top(x-y-1)+top(x-y)+2)>>2)
				case x < y:
					set(x, y, (left(y-x-2)+2*left(y-x-1)+left(y-x)+2)>>2)
				default:
					set(x, y, (top(0)+2*top(-1)+left(0)+2)>>2)
				}
			}
		}
	case 5: // Vertical_Right
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				zVR := 2*x - y
				switch {
				case zVR >= 0 && zVR%2 == 0:
					set(x, y, (top(x-(y>>1)-1)+top(x-(y>>1))+1)>>1)
				case zVR >= 0:
					set(x, y, (top(x-(y>>1)-2)+2*top(x-(y>>1)-1)+top(x-(y>>1))+2)>>2)
				case zVR == -1:
					set(x, y, (left(0)+2*top(-1)+top(0)+2)>>2)
				default:
					set(x, y, (left(y-1)+2*left(y-2)+left(y-3)+2)>>2)
				}
			}
		}
	case 6: // Horizontal_Down
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				zHD := 2*y - x
				switch {
				case zHD >= 0 && zHD%2 == 0:
					set(x, y, (left(y-(x>>1)-1)+left(y-(x>>1))+1)>>1)
				case zHD >= 0:
					set(x, y, (left(y-(x>>1)-2)+2*left(y-(x>>1)-1)+left(y-(x>>1))+2)>>2)
				case zHD == -1:
					set(x, y, (left(0)+2*top(-1)+top(0)+2)>>2)
				default:
					set(x, y, (top(x-1)+2*top(x-2)+top(x-3)+2)>>2)
				}
			}
		}
	case 7: // Vertical_Left
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if y%2 == 0 {
					set(x, y, (top(x+(y>>1))+top(x+(y>>1)+1)+1)>>1)
				} else {
					set(x, y, (top(x+(y>>1))+2*top(x+(y>>1)+1)+top(x+(y>>1)+2)+2)>>2)
				}
			}
		}
	case 8: // Horizontal_Up
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				zHU := x + 2*y
				switch {
				case zHU < 13 && zHU%2 == 0:
					set(x, y, (left(y+(x>>1))+left(y+(x>>1)+1)+1)>>1)
				case zHU < 13:
					set(x, y, (left(y+(x>>1))+2*left(y+(x>>1)+1)+left(y+(x>>1)+2)+2)>>2)
				case zHU == 13:
					set(x, y, (left(6)+3*left(7)+2)>>2)
				default:
					set(x, y, left(7))
				}
			}
		}
	}
	return out
}

// predictIntra16x16 implements the 4 Intra_16x16 luma modes (§8.3.3).
func predictIntra16x16(p *planeView, x0, y0, mode int, haveLeft, haveTop, haveTopLeft bool) []int {
	return predictVerticalHorizontalDCPlane(p, x0, y0, 16, mode, haveLeft, haveTop, haveTopLeft, 5)
}

// predictIntraChroma implements the 4 chroma intra modes (§8.3.4) for one
// 8x8 chroma block. Chroma's mode numbering (Table 8-4: 0=DC, 1=Horizontal,
// 2=Vertical, 3=Plane) is NOT the same as Intra_16x16's (Table 8-3: 0=
// Vertical, 1=Horizontal, 2=DC, 3=Plane), which predictVerticalHorizontalDCPlane
// is written against - mode 2 (Vertical for chroma) is translated to that
// function's mode 0 below. Mode 0 (DC) uses chroma's own per-quadrant
// averaging rule (§8.3.4.1), not the simple whole-block DC used by luma, so
// it's handled separately rather than translated to that function's mode 2.
// A first attempt at this decoder used the Intra_16x16 numbering directly
// (routing mode 0 into "Vertical" and treating mode 2 as "DC"), which read
// out-of-bounds/garbage samples for the very first macroblock (no top
// neighbor) and cascaded a uniform wrong chroma cast across the whole frame.
func predictIntraChroma(p *planeView, x0, y0, mode int, haveLeft, haveTop, haveTopLeft bool) []int {
	if mode != 0 {
		translated := mode
		if mode == 2 {
			translated = 0
		}
		return predictVerticalHorizontalDCPlane(p, x0, y0, 8, translated, haveLeft, haveTop, haveTopLeft, 17)
	}

	out := make([]int, 64)
	quad := func(qx, qy int, useTop, useLeft bool) {
		var sum, n int
		if useTop {
			for i := 0; i < 4; i++ {
				sum += p.at(x0+qx+i, y0-1)
			}
			n += 4
		}
		if useLeft {
			for i := 0; i < 4; i++ {
				sum += p.at(x0-1, y0+qy+i)
			}
			n += 4
		}
		dc := 128
		if n > 0 {
			dc = (sum + n/2) / n
		}
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				out[(qy+y)*8+(qx+x)] = dc
			}
		}
	}
	quad(0, 0, haveTop, haveLeft)
	quad(4, 0, haveTop, !haveTop && haveLeft)
	quad(0, 4, !haveLeft && haveTop, haveLeft)
	quad(4, 4, haveTop, haveLeft)
	return out
}

// substituteRefSamples implements the reference sample substitution process
// (§8.3.1.2.1-ish - the general "unavailable sample" rule used across the
// Intra prediction clauses) for a plain top-row/left-column/corner neighbor
// set (no top-right needed: Vertical/Horizontal/DC/Plane only ever read
// these). Per spec, samples are considered in scan order left-column
// bottom-to-top, then the corner, then top-row left-to-right; an unavailable
// sample takes the value of the previous sample in that order, or 128 if
// it's the very first one. A first attempt at this decoder read directly
// from the plane buffer with no substitution at all, so any unavailable
// neighbor silently read back 0 (the zero-initialized plane) instead of a
// real substituted value - for a macroblock with a real left neighbor but
// no top one (true for an entire picture's top row beyond its first
// column), that turned Vertical/Plane prediction black instead of
// propagating the real left-neighbor color sideways.
func substituteRefSamples(p *planeView, x0, y0, size int, haveLeft, haveTop, haveTopLeft bool) (top, left []int, corner int) {
	top = make([]int, size)
	left = make([]int, size)
	for i := 0; i < size; i++ {
		if haveTop {
			top[i] = p.at(x0+i, y0-1)
		}
		if haveLeft {
			left[i] = p.at(x0-1, y0+i)
		}
	}
	if haveTopLeft {
		corner = p.at(x0-1, y0-1)
	}

	avail := make([]bool, 2*size+1)
	vals := make([]int, 2*size+1)
	for i := 0; i < size; i++ {
		avail[i] = haveLeft
		vals[i] = left[size-1-i]
	}
	avail[size] = haveTopLeft
	vals[size] = corner
	for i := 0; i < size; i++ {
		avail[size+1+i] = haveTop
		vals[size+1+i] = top[i]
	}

	prev := 128
	for i := range vals {
		if avail[i] {
			prev = vals[i]
		} else {
			vals[i] = prev
		}
	}

	for i := 0; i < size; i++ {
		left[size-1-i] = vals[i]
	}
	corner = vals[size]
	for i := 0; i < size; i++ {
		top[i] = vals[size+1+i]
	}
	return top, left, corner
}

// predictVerticalHorizontalDCPlane implements the vertical/horizontal/DC/
// plane modes shared by Intra_16x16 (size=16, planeCoeff=5) and chroma's
// non-DC modes (size=8, planeCoeff=17) - §8.3.3 and §8.3.4.2-4.
func predictVerticalHorizontalDCPlane(p *planeView, x0, y0, size, mode int, haveLeft, haveTop, haveTopLeft bool, planeCoeff int) []int {
	out := make([]int, size*size)
	set := func(x, y, v int) { out[y*size+x] = v }

	if mode == 0 || mode == 1 || mode == 3 {
		top, left, corner := substituteRefSamples(p, x0, y0, size, haveLeft, haveTop, haveTopLeft)
		switch mode {
		case 0: // Vertical
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					set(x, y, top[x])
				}
			}
		case 1: // Horizontal
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					set(x, y, left[y])
				}
			}
		case 3: // Plane
			half := size/2 - 1
			h, v := 0, 0
			topAt := func(i int) int {
				if i < 0 {
					return corner
				}
				return top[i]
			}
			leftAt := func(i int) int {
				if i < 0 {
					return corner
				}
				return left[i]
			}
			for i := 0; i <= half; i++ {
				weight := i + 1
				h += weight * (topAt(half+1+i) - topAt(half-1-i))
				v += weight * (leftAt(half+1+i) - leftAt(half-1-i))
			}
			a := 16 * (left[size-1] + top[size-1])
			// Luma (§8.3.3.4, size=16): b=(5H+32)>>6. Chroma (§8.3.4.4,
			// size=8): b=(17H+16)>>5 - a different rounding constant and
			// shift, not just a different multiplier, hence the size-based
			// branch here.
			var b, c int
			if size == 16 {
				b = (planeCoeff*h + 32) >> 6
				c = (planeCoeff*v + 32) >> 6
			} else {
				b = (planeCoeff*h + 16) >> 5
				c = (planeCoeff*v + 16) >> 5
			}
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					set(x, y, clip3(0, 255, (a+b*(x-half-1)+c*(y-half-1)+16)>>5))
				}
			}
		}
		return out
	}

	switch mode {
	case 2: // DC (whole-block average; chroma's per-quadrant DC is handled separately)
		var sum, n int
		if haveTop {
			for x := 0; x < size; x++ {
				sum += p.at(x0+x, y0-1)
			}
			n += size
		}
		if haveLeft {
			for y := 0; y < size; y++ {
				sum += p.at(x0-1, y0+y)
			}
			n += size
		}
		dc := 128
		if n > 0 {
			dc = (sum + n/2) / n
		}
		for i := range out {
			out[i] = dc
		}
	}
	return out
}
