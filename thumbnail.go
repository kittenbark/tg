package tg

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	_ "image/png" // register the PNG decoder for image.Decode.
	"os"
)

const (
	// thumbnailMaxDim is Telegram's hard limit: a thumbnail's width and
	// height must not exceed 320 (see the Thumbnail field doc comments on
	// InputMediaVideo/Animation/Audio/Document in api_types.go).
	thumbnailMaxDim = 320
	// thumbnailMaxBytes is Telegram's hard limit: a thumbnail must be less
	// than 200 kB.
	thumbnailMaxBytes     = 200_000
	thumbnailJPEGQuality  = 85
	thumbnailRetryQuality = 60
)

var errThumbnailTooLarge = errors.New("tg: generated thumbnail exceeds telegram's 200kb limit")

// probeStillImage decodes a JPEG/PNG file and returns its dimensions plus a
// generated thumbnail. Returns nil on any failure.
func probeStillImage(path string) *mediaMeta {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}

	bounds := img.Bounds()
	meta := &mediaMeta{width: int64(bounds.Dx()), height: int64(bounds.Dy())}
	if path, err := writeJPEGThumbnail(img); err == nil {
		meta.tempThumbnailPath = path
	}
	return meta
}

// probeGIF decodes a GIF file and returns its dimensions, total animation
// duration (sum of per-frame delays), and a thumbnail generated from its
// first frame. Returns nil on any failure.
func probeGIF(path string) *mediaMeta {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	decoded, err := gif.DecodeAll(file)
	if err != nil || len(decoded.Image) == 0 {
		return nil
	}

	var totalCentiseconds int
	for _, delay := range decoded.Delay {
		totalCentiseconds += delay
	}

	meta := &mediaMeta{
		width:           int64(decoded.Config.Width),
		height:          int64(decoded.Config.Height),
		durationSeconds: int64(totalCentiseconds) / 100,
	}
	if path, err := writeJPEGThumbnail(decoded.Image[0]); err == nil {
		meta.tempThumbnailPath = path
	}
	return meta
}

// writeJPEGThumbnail downsamples img to fit within thumbnailMaxDim on both
// axes (preserving aspect ratio), encodes it as JPEG, and writes it to a new
// temp file - the caller owns removing it. It verifies the encoded result
// respects Telegram's 200kB cap (retrying once at a lower quality), since
// unlike width/height/duration, the Thumbnail field is validated
// server-side: attaching one that violates the limits would turn this
// best-effort feature into a hard failure of the real Send call.
func writeJPEGThumbnail(img image.Image) (string, error) {
	thumb := resizeToFit(img, thumbnailMaxDim)

	tmp, err := os.CreateTemp("", "tg-thumbnail-*.jpg")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	defer func() { _ = tmp.Close() }()

	for _, quality := range []int{thumbnailJPEGQuality, thumbnailRetryQuality} {
		if _, err := tmp.Seek(0, 0); err != nil {
			_ = os.Remove(name)
			return "", err
		}
		if err := jpeg.Encode(tmp, thumb, &jpeg.Options{Quality: quality}); err != nil {
			_ = os.Remove(name)
			return "", err
		}
		info, err := tmp.Stat()
		if err != nil {
			_ = os.Remove(name)
			return "", err
		}
		if info.Size() <= thumbnailMaxBytes {
			_ = tmp.Truncate(info.Size())
			return name, nil
		}
	}
	_ = os.Remove(name)
	return "", errThumbnailTooLarge
}

// fitDimensions returns target dimensions that fit within maxDim on both
// axes, preserving aspect ratio. Returns the input unchanged if it already
// fits.
func fitDimensions(width, height, maxDim int) (int, int) {
	if width <= maxDim && height <= maxDim {
		return width, height
	}
	if width >= height {
		return maxDim, max(1, height*maxDim/width)
	}
	return max(1, width*maxDim/height), maxDim
}

// resizeToFit downscales src to fit within maxDim via box-averaging: each
// destination pixel is the average of its corresponding source rectangle.
// Chosen over nearest-neighbor because, for a pure downscale (large source
// -> <=320px thumbnail), box-averaging avoids the aliasing/moire
// nearest-neighbor produces - a visible quality difference for a
// user-facing thumbnail - while still only visiting every source pixel
// once overall, so the extra cost is negligible next to the JPEG decode
// already paid for. Returns src unchanged if it already fits.
func resizeToFit(src image.Image, maxDim int) image.Image {
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := fitDimensions(srcW, srcH, maxDim)
	if dstW == srcW && dstH == srcH {
		return src
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	xRatio := float64(srcW) / float64(dstW)
	yRatio := float64(srcH) / float64(dstH)
	for dy := 0; dy < dstH; dy++ {
		y0 := int(float64(dy) * yRatio)
		y1 := int(float64(dy+1) * yRatio)
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for dx := 0; dx < dstW; dx++ {
			x0 := int(float64(dx) * xRatio)
			x1 := int(float64(dx+1) * xRatio)
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var rSum, gSum, bSum, count uint64
			for sy := y0; sy < y1 && sy < srcH; sy++ {
				for sx := x0; sx < x1 && sx < srcW; sx++ {
					r, g, b, _ := src.At(bounds.Min.X+sx, bounds.Min.Y+sy).RGBA()
					rSum += uint64(r >> 8)
					gSum += uint64(g >> 8)
					bSum += uint64(b >> 8)
					count++
				}
			}
			if count == 0 {
				count = 1
			}
			dst.Set(dx, dy, color.RGBA{
				R: uint8(rSum / count),
				G: uint8(gSum / count),
				B: uint8(bSum / count),
				A: 255,
			})
		}
	}
	return dst
}
