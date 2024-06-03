package main

import (
	"log"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

// vipsImage is an image managed by the vips library
// Attempts to use this library indicates that there are
// concurrency issues under heavy use - when loading and processing
// multiple images in parallel, the library will sometimes fail with
// a SIGSEGV, or a thread will hang.
type vipsImage struct {
	Image
	img *vips.ImageRef
}

func vipsInit() {
	// Make vips less noisy.
	vips.LoggingSettings(nil, vips.LogLevelError)
	vips.Startup(nil)
}

// NewVipsImage returns an image loaded and managed by the
// cgo bindings to libvips.
func NewVipsImage(src string) (Image, error) {
	vimg, err := vips.NewImageFromFile(src)
	if err != nil {
		return nil, err
	}
	return &vipsImage{img: vimg}, nil
}

func (v *vipsImage) Width() int {
	return v.img.Width()
}

func (v *vipsImage) Height() int {
	return v.img.Height()
}

func (v *vipsImage) Rotate(deg RotateDegrees) error {
	switch deg {
	case Rotate90:
		v.img.Rotate(vips.Angle90)
	case Rotate180:
		v.img.Rotate(vips.Angle180)
	case Rotate270:
		v.img.Rotate(vips.Angle270)
	}
	return nil
}

func (v *vipsImage) Write(destFile string, mtime time.Time, w, h, q int) {
	vimg, err := v.img.Copy()
	if err != nil {
		log.Fatalf("%s: copy: %v", destFile, err)
	}
	// scale it down to fit within the width & height
	xr := float64(w) / float64(vimg.Width())
	yr := float64(h) / float64(vimg.Height())
	if xr < 1 || yr < 1 {
		// Image is larger than requested size, so scale it down
		var err error
		if xr < yr {
			err = vimg.Resize(xr, vips.KernelAuto)
		} else {
			err = vimg.Resize(yr, vips.KernelAuto)
		}
		if err != nil {
			log.Fatalf("%s: resize: %v", destFile, err)
		}
	}
	jp := vips.NewJpegExportParams()
	jp.Quality = q
	b, _, err := vimg.ExportJpeg(jp)
	if err != nil {
		log.Fatalf("%s: export: %v", destFile, err)
	}
	if err := cp(b, destFile, mtime); err != nil {
		log.Fatalf("%s: write: %v", destFile, err)
	}
}
