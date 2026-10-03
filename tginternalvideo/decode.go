package tginternalvideo

import (
	"image"
	"image/color"
)

// DecodeFirstFrame best-effort decodes an H.264 video sample's first IDR
// slice into an image: avcCPayload is the raw avcC (AVCDecoderConfigurationRecord)
// box payload (ISO/IEC 14496-15 §5.2.4.1) from the track's sample entry, and
// sampleData is the sample's raw AVCC length-prefixed NAL data. SPS/PPS come
// from the avcC config, the sample's NAL units are demuxed to find the IDR
// slice, and decodeIDRPictureI does the actual CABAC/intra/transform work.
//
// Any unsupported feature or malformed input returns an error - this covers
// everything out of this decoder's intentionally narrow scope (I-slice-only,
// CABAC-only, no custom scaling lists, single slice per picture, 4:2:0 8-bit
// only, no deblocking filter) as well as genuinely malformed input. Callers
// must treat every error as "no thumbnail available" and fail open to
// metadata-only, never surface it as a real error.
func DecodeFirstFrame(avcCPayload, sampleData []byte) (image.Image, error) {
	cfg, err := parseAVCDecoderConfigurationRecord(avcCPayload)
	if err != nil {
		return nil, err
	}
	if len(cfg.sps) == 0 || len(cfg.pps) == 0 {
		return nil, errH264Malformed
	}

	spsType, spsRBSP, err := nalUnitType(cfg.sps[0])
	if err != nil || spsType != 7 {
		return nil, errH264Malformed
	}
	sps, err := parseSPS(spsRBSP)
	if err != nil {
		return nil, err
	}

	ppsType, ppsRBSP, err := nalUnitType(cfg.pps[0])
	if err != nil || ppsType != 8 {
		return nil, errH264Malformed
	}
	pps, err := parsePPS(ppsRBSP)
	if err != nil {
		return nil, err
	}

	var idrRBSP []byte
	for _, nal := range splitNALUnits(sampleData, cfg.lengthSize) {
		naluType, rbsp, err := nalUnitType(nal)
		if err != nil {
			continue
		}
		if naluType == 5 { // IDR slice
			idrRBSP = rbsp
			break
		}
	}
	if idrRBSP == nil {
		return nil, errH264Unsupported // sample #1 has no IDR slice (unexpected for a sync sample)
	}

	y, cb, cr, err := decodeIDRPictureI(sps, pps, idrRBSP)
	if err != nil {
		return nil, err
	}

	width, height := sps.displaySize()
	if width <= 0 || height <= 0 || width > y.width || height > y.height {
		return nil, errH264Malformed
	}
	return yuv420ToRGBA(y, cb, cr, width, height, sps.videoFullRangeFlag), nil
}

// yuv420ToRGBA converts reconstructed 4:2:0 planes to RGB, honoring
// video_full_range_flag (ITU-T H.264's VUI) - this is hand-written rather
// than using image.YCbCr because Go's image/color.YCbCrToRGB always assumes
// limited/studio range, which would visibly wash out a full-range source.
func yuv420ToRGBA(y, cb, cr *planeView, width, height int, fullRange bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			yv := float64(y.at(px, py))
			cbv := float64(cb.at(px/2, py/2)) - 128
			crv := float64(cr.at(px/2, py/2)) - 128

			var r, g, b float64
			if fullRange {
				r = yv + 1.402*crv
				g = yv - 0.344136*cbv - 0.714136*crv
				b = yv + 1.772*cbv
			} else {
				yv = 1.164 * (yv - 16)
				r = yv + 1.596*crv
				g = yv - 0.392*cbv - 0.813*crv
				b = yv + 2.017*cbv
			}
			img.SetRGBA(px, py, color.RGBA{
				R: uint8(clip3(0, 255, int(r+0.5))),
				G: uint8(clip3(0, 255, int(g+0.5))),
				B: uint8(clip3(0, 255, int(b+0.5))),
				A: 255,
			})
		}
	}
	return img
}
