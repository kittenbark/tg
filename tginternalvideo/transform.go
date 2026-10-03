package tginternalvideo

// normAdjust4x4 is v4x4 (§8.5.9): normAdjust4x4[QP%6][category], where
// category groups a 4x4 position (i,j) by (i%2,j%2): 0 = both even,
// 1 = both odd, 2 = mixed. Combined with the flat (default) weightScale=16
// this decoder always uses (no custom scaling lists, per v1 scope), this
// gives LevelScale4x4 = 16 * normAdjust4x4.
var normAdjust4x4 = [6][3]int{
	{10, 16, 13},
	{11, 18, 14},
	{13, 20, 16},
	{14, 23, 18},
	{16, 25, 20},
	{18, 29, 23},
}

func normAdjust4x4Category(i, j int) int {
	switch {
	case i%2 == 0 && j%2 == 0:
		return 0
	case i%2 == 1 && j%2 == 1:
		return 1
	default:
		return 2
	}
}

func levelScale4x4(qp, i, j int) int {
	return 16 * normAdjust4x4[qp%6][normAdjust4x4Category(i, j)]
}

// normAdjust8x8 is v8x8 (§8.5.9) for the 8x8 transform, with position (i,j)
// (mod 4 in each dimension) mapped to one of 6 categories.
var normAdjust8x8 = [6][6]int{
	{20, 18, 32, 19, 25, 24},
	{22, 19, 35, 21, 28, 26},
	{26, 23, 42, 24, 33, 31},
	{28, 25, 45, 26, 35, 33},
	{32, 28, 51, 30, 40, 38},
	{36, 32, 58, 34, 46, 43},
}

var normAdjust8x8CategoryTable = [4][4]int{
	{0, 3, 4, 3},
	{3, 1, 5, 1},
	{4, 5, 2, 5},
	{3, 1, 5, 1},
}

func levelScale8x8(qp, i, j int) int {
	return 16 * normAdjust8x8[qp%6][normAdjust8x8CategoryTable[i%4][j%4]]
}

// dequantize4x4AC dequantizes a 4x4 (or Intra16x16 AC / ChromaAC, which
// reuse the same per-position scaling) block of coefficients in raster
// order (§8.5.12.1); coeff[0] is ignored for Intra16x16AC/ChromaAC callers,
// which supply it separately via the Hadamard-transformed DC.
func dequantize4x4(coeff []int32, qp int) [16]int32 {
	var out [16]int32
	shift := qp / 6
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			idx := i*4 + j
			ls := levelScale4x4(qp, i, j)
			c := int64(coeff[idx])
			if shift >= 4 {
				out[idx] = int32(c * int64(ls) << uint(shift-4))
			} else {
				out[idx] = int32((c*int64(ls) + int64(1<<uint(3-shift))) >> uint(4-shift))
			}
		}
	}
	return out
}

func dequantize8x8(coeff []int32, qp int) [64]int32 {
	var out [64]int32
	shift := qp / 6
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			idx := i*8 + j
			ls := levelScale8x8(qp, i, j)
			c := int64(coeff[idx])
			if shift >= 6 {
				out[idx] = int32(c * int64(ls) << uint(shift-6))
			} else {
				out[idx] = int32((c*int64(ls) + int64(1<<uint(5-shift))) >> uint(6-shift))
			}
		}
	}
	return out
}

// hadamard4x4Inverse applies the plain (unweighted) 4x4 Hadamard transform
// used for Intra16x16 luma DC coefficients (§8.5.10) - row pass then column
// pass of the same symmetric butterfly (distinct from idct4x4's butterfly,
// which has asymmetric >>1 terms the DC transform does not).
func hadamard4x4Inverse(d [16]int32) [16]int32 {
	had1D := func(d0, d1, d2, d3 int32) (int32, int32, int32, int32) {
		e0, e1, e2, e3 := d0+d2, d0-d2, d1-d3, d1+d3
		return e0 + e3, e1 + e2, e1 - e2, e0 - e3
	}
	var tmp [16]int32
	for i := 0; i < 4; i++ {
		tmp[i*4+0], tmp[i*4+1], tmp[i*4+2], tmp[i*4+3] = had1D(d[i*4+0], d[i*4+1], d[i*4+2], d[i*4+3])
	}
	var out [16]int32
	for j := 0; j < 4; j++ {
		out[0*4+j], out[1*4+j], out[2*4+j], out[3*4+j] = had1D(tmp[0*4+j], tmp[1*4+j], tmp[2*4+j], tmp[3*4+j])
	}
	return out
}

// dequantizeLumaDC scales a Hadamard-transformed Intra16x16 luma DC array
// (§8.5.10, continuing the scaling after hadamard4x4Inverse).
func dequantizeLumaDC(f [16]int32, qp int) [16]int32 {
	ls := levelScale4x4(qp, 0, 0)
	shift := qp / 6
	var out [16]int32
	for i, v := range f {
		c := int64(v)
		if shift >= 6 {
			out[i] = int32(c * int64(ls) << uint(shift-6))
		} else {
			out[i] = int32((c*int64(ls) + int64(1<<uint(5-shift))) >> uint(6-shift))
		}
	}
	return out
}

// hadamard2x2Inverse applies the 2x2 Hadamard transform used for ChromaDC
// coefficients in 4:2:0 (§8.5.11.1). d is raster order (TL,TR,BL,BR).
func hadamard2x2Inverse(d [4]int32) [4]int32 {
	return [4]int32{
		d[0] + d[1] + d[2] + d[3],
		d[0] - d[1] + d[2] - d[3],
		d[0] + d[1] - d[2] - d[3],
		d[0] - d[1] - d[2] + d[3],
	}
}

// dequantizeChromaDC scales a Hadamard-transformed 4:2:0 ChromaDC array
// (§8.5.11.2).
func dequantizeChromaDC(f [4]int32, qp int) [4]int32 {
	ls := levelScale4x4(qp, 0, 0)
	shift := qp / 6
	var out [4]int32
	for i, v := range f {
		out[i] = int32((int64(v) * int64(ls) << uint(shift)) >> 5)
	}
	return out
}

// idct4x4 is the inverse 4x4 core transform (§8.5.12.2): row pass, column
// pass of the same asymmetric butterfly, then (+32)>>6.
func idct4x4(d [16]int32) [16]int32 {
	core1D := func(d0, d1, d2, d3 int32) (int32, int32, int32, int32) {
		e0 := d0 + d2
		e1 := d0 - d2
		e2 := (d1 >> 1) - d3
		e3 := d1 + (d3 >> 1)
		return e0 + e3, e1 + e2, e1 - e2, e0 - e3
	}
	var tmp [16]int32
	for i := 0; i < 4; i++ {
		tmp[i*4+0], tmp[i*4+1], tmp[i*4+2], tmp[i*4+3] = core1D(d[i*4+0], d[i*4+1], d[i*4+2], d[i*4+3])
	}
	var col [16]int32
	for j := 0; j < 4; j++ {
		col[0*4+j], col[1*4+j], col[2*4+j], col[3*4+j] = core1D(tmp[0*4+j], tmp[1*4+j], tmp[2*4+j], tmp[3*4+j])
	}
	var out [16]int32
	for i, v := range col {
		out[i] = (v + 32) >> 6
	}
	return out
}

// idct8x8 is the inverse 8x8 core transform (§8.5.13.2): row pass, column
// pass of the 8-point butterfly, then (+32)>>6.
func idct8x8(d [64]int32) [64]int32 {
	core1D := func(d0, d1, d2, d3, d4, d5, d6, d7 int32) [8]int32 {
		a0 := d0 + d4
		a2 := d0 - d4
		a4 := (d2 >> 1) - d6
		a6 := d2 + (d6 >> 1)

		b0 := a0 + a6
		b2 := a2 + a4
		b4 := a2 - a4
		b6 := a0 - a6

		a1 := -d3 + d5 - d7 - (d7 >> 1)
		a3 := d1 + d7 - d3 - (d3 >> 1)
		a5 := -d1 + d7 + d5 + (d5 >> 1)
		a7 := d3 + d5 + d1 + (d1 >> 1)

		b1 := a1 + (a7 >> 2)
		b7 := a7 - (a1 >> 2)
		b3 := a3 + (a5 >> 2)
		b5 := (a3 >> 2) - a5

		var f [8]int32
		f[0] = b0 + b7
		f[7] = b0 - b7
		f[1] = b2 + b5
		f[6] = b2 - b5
		f[2] = b4 + b3
		f[5] = b4 - b3
		f[3] = b6 + b1
		f[4] = b6 - b1
		return f
	}

	var tmp [64]int32
	for i := 0; i < 8; i++ {
		row := core1D(d[i*8+0], d[i*8+1], d[i*8+2], d[i*8+3], d[i*8+4], d[i*8+5], d[i*8+6], d[i*8+7])
		for j := 0; j < 8; j++ {
			tmp[i*8+j] = row[j]
		}
	}
	var col [64]int32
	for j := 0; j < 8; j++ {
		c := core1D(tmp[0*8+j], tmp[1*8+j], tmp[2*8+j], tmp[3*8+j], tmp[4*8+j], tmp[5*8+j], tmp[6*8+j], tmp[7*8+j])
		for i := 0; i < 8; i++ {
			col[i*8+j] = c[i]
		}
	}
	var out [64]int32
	for i, v := range col {
		out[i] = (v + 32) >> 6
	}
	return out
}
