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
	thumb := toYCbCr(resizeToFit(img, thumbnailMaxDim))

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

// rgbGetter returns a closure reading (r,g,b) as 0-255 values at absolute
// image coordinates. *image.RGBA - what both this package's own thumbnail
// output and, notably, tginternalvideo's decoded video frames are - gets a
// fast path straight into its Pix buffer, skipping image.Image.At()'s
// color.Color interface boxing, which before this was resizeToFit's (and so
// thumbnail generation's) single largest allocation source by a wide margin:
// box-averaging a 720x1280 frame down to a thumbnail visits most of its
// ~920,000 source pixels once each, and at one boxed color.Color per visit
// that's the bulk of a thumbnail's total allocations. Reading Pix directly
// is exactly equivalent to *image.RGBA's own At(x,y).RGBA() (RGBA's doc
// comment: Pix holds already alpha-premultiplied values, and color.RGBA's
// RGBA() method only widens 8-bit to 16-bit, it doesn't premultiply again) -
// safe regardless of alpha. Anything else (e.g. a decoded GIF frame's
// image.Paletted) falls back to the generic, always-correct interface path.
func rgbGetter(src image.Image) func(x, y int) (r, g, b uint8) {
	if rgba, ok := src.(*image.RGBA); ok {
		return func(x, y int) (uint8, uint8, uint8) {
			i := rgba.PixOffset(x, y)
			p := rgba.Pix[i : i+3 : i+3]
			return p[0], p[1], p[2]
		}
	}
	return func(x, y int) (uint8, uint8, uint8) {
		r, g, b, _ := src.At(x, y).RGBA()
		return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
	}
}

// toYCbCr converts img to *image.YCbCr (4:2:0 chroma subsampling, matching
// this package's own video-frame decoder and typical JPEG output) - the one
// concrete type image/jpeg's encoder has a fast per-pixel path for. Encoding
// anything else (including the *image.RGBA resizeToFit itself produces)
// makes the stdlib encoder fall back to boxed image.Image.At() calls for
// every pixel during encoding, which - once resizeToFit's own equivalent
// At() cost was fixed - became thumbnail generation's next largest
// allocation source. Chroma is sampled from each 2x2 block's top-left pixel
// rather than averaged: a thumbnail is small and already lossy, and a single
// sample avoids a second full pass over the image for a difference that
// isn't visible at this size.
func toYCbCr(img image.Image) *image.YCbCr {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewYCbCr(image.Rect(0, 0, w, h), image.YCbCrSubsampleRatio420)
	getRGB := rgbGetter(img)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b := getRGB(bounds.Min.X+x, bounds.Min.Y+y)
			yy, cb, cr := color.RGBToYCbCr(r, g, b)
			dst.Y[dst.YOffset(x, y)] = yy
			if x%2 == 0 && y%2 == 0 {
				ci := dst.COffset(x, y)
				dst.Cb[ci] = cb
				dst.Cr[ci] = cr
			}
		}
	}
	return dst
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

	getRGB := rgbGetter(src)

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
					r, g, b := getRGB(bounds.Min.X+sx, bounds.Min.Y+sy)
					rSum += uint64(r)
					gSum += uint64(g)
					bSum += uint64(b)
					count++
				}
			}
			if count == 0 {
				count = 1
			}
			// dst.Set(dx, dy, color.RGBA{...}) would box that value into a
			// color.Color interface on every call - a single allocation,
			// but incurred once per destination pixel, so (after fixing
			// getRGB above) this became resizeToFit's last remaining
			// allocation source. dst is always our own freshly-made
			// *image.RGBA, so writing its Pix buffer directly is just as
			// correct and allocates nothing.
			i := dst.PixOffset(dx, dy)
			dst.Pix[i+0] = uint8(rSum / count)
			dst.Pix[i+1] = uint8(gSum / count)
			dst.Pix[i+2] = uint8(bSum / count)
			dst.Pix[i+3] = 255
		}
	}
	return dst
}
