# tginternalvideo

A pure-Go, zero-dependency, best-effort H.264 first-frame decoder, extracted
from [kittenbark/tg](https://github.com/kittenbark/tg) as its own nested
module. It exists for exactly one reason: let `tg` generate a *real* video
thumbnail for a local upload without requiring `ffmpeg` (or any other binary)
to be installed on the host.

## Scope

This is intentionally not a general-purpose H.264 decoder:

- I-slice (IDR) only, intra prediction only - no inter/motion compensation.
- CABAC entropy coding only (no CAVLC).
- Single slice per picture, progressive frames, 4:2:0 8-bit only.
- Flat/default scaling lists only (no custom quantization matrices).
- No deblocking filter.

Any real-world H.264 sample outside that scope (CAVLC, multi-slice, custom
scaling lists, 4:2:2/4:4:4, etc.) is rejected with an error. There is exactly
one exported entry point:

```go
func DecodeFirstFrame(avcCPayload, sampleData []byte) (image.Image, error)
```

`tg` is responsible for every bit of container parsing (locating the video
track, the avcC box, and sample #1's offset/size in the MP4); this package
only ever sees raw avcC bytes and raw sample bytes in, and an `image.Image`
or an error out. Callers must treat every error as "no thumbnail available"
and fail open - this decoder is a nice-to-have, never a hard dependency.

Decode correctness is checked against real `ffmpeg`-decoded reference frames
in `decode_test.go` (`TestH264DecodeFirstFrameAgainstFFmpeg`), gated on mean
absolute pixel difference rather than exact match, since this decoder skips
the deblocking filter.

## Benchmarks

Measured on an Apple M2, comparing this package against spawning `ffmpeg` as
a subprocess to do the same job (`go test -bench . ./...` in `tg`'s own
`video_thumbnail_bench_test.go`, and `BenchmarkDecodeFirstFrame` in this
package for the decoder in isolation).

Decoding alone (`DecodeFirstFrame`, 720x1280 H.264 frame, no JPEG
resize/encode):

| | time/op | bytes/op | allocs/op |
|---|---|---|---|
| this package | 19.7 ms | 5.8 MB | 21 |

Full pipeline (`probeISOBMFF`'s video-thumbnail path: box-walk, decode,
resize to fit 320px, encode JPEG) vs. spawning `ffmpeg`:

| fixture | this package | `ffmpeg` subprocess |
|---|---|---|
| video.mp4 (720x1280) | **25.1 ms**, 6.2 MB, 108 allocs | 33.9 ms, 16 KB, 56 allocs |
| bigger_2.mp4 (794x1058, cropped) | **31.6 ms**, 5.9 MB, 118 allocs | 45.5 ms, 16 KB, 56 allocs |

This package is consistently ~25-30% faster than the subprocess approach on
both real fixtures, with allocation counts in the same ballpark (the
remaining gap - roughly 2x the subprocess's allocation count, in absolute
terms a rounding error at 108-118 allocations total) is GC/runtime
bookkeeping, not application-level garbage.

That wasn't true of the first working version: an early allocation profile
of this same benchmark showed **48,532 allocations and 20.5 MB per decoded
frame**, and the full pipeline (decode + resize + JPEG-encode) showed
**1,027,816 allocations and 24.7 MB per thumbnail** - roughly tied with the
`ffmpeg` subprocess on wall-clock time, but at a steep GC-pressure cost for
no latency benefit. Two rounds of `pprof -alloc_objects` profiling (one on
this package's own decode path, one on `tg`'s surrounding resize/encode
code) found exactly where that was coming from and fixed it:

1. **Per-macroblock prediction/residual buffers.** CABAC intra prediction
   runs once per macroblock - thousands of times per frame - and several of
   its functions (`substituteRefSamples`, `predictVerticalHorizontalDCPlane`,
   `predictIntra4x4`/`predictIntra8x8`/`predictIntraChroma`,
   `decodeResidualBlockCABAC`) returned a freshly `make()`'d `[]int`/`[]int32`
   slice every call, for a size that was always known and small (<=256
   elements). Switching those return types to fixed-size arrays (`[16]int`,
   `[64]int`, `[256]int`, `[64]int32`) lets the compiler keep them on the
   stack instead of the heap. `substituteRefSamples` alone was 40% of every
   allocation this package made; fixing all of them took decode from 48,532
   allocations down to 21.
2. **`image.Image`/`color.Color` interface boxing**, found in `tg`'s
   surrounding `resizeToFit`/`writeJPEGThumbnail` once this package's own
   allocations stopped dominating the profile: `src.At(x, y)` and
   `dst.Set(x, y, color.RGBA{...})` both box a small concrete value into an
   interface on every call, and `image/jpeg`'s encoder only has an
   allocation-free fast path for `*image.YCbCr`. Reading/writing the
   `*image.RGBA` `Pix` buffer directly, and converting to `*image.YCbCr`
   before encoding, closed the rest of the gap.

Decoder output is bit-for-bit identical before and after every one of these
changes (checked against the same `ffmpeg`-reference test on every step) -
this was purely about how the same pixels get computed and packaged, not
what they are.

## Why a separate module

`tg` is zero-dependency by design, and this decoder - CABAC entropy coding,
intra prediction, inverse transforms - is a categorically bigger and more
fragile piece of code than anything else in that repository. Keeping it in
its own module with a single, narrow exported function keeps that complexity
contained and independently testable/benchmarkable, without `tg` itself
depending on anything beyond the Go standard library plus this one sibling
package (wired in locally via a `replace` directive in `tg`'s `go.mod` - no
separate publish step is needed for development in this monorepo).
