package vips

import (
	"fmt"
	"os"
	"time"

	"github.com/aamcrae/pweb/imager"
	"github.com/davidbyttow/govips/v2/vips"
)

// vipsImage is an image managed by the vips library
// Attempts to use this library indicates that there are
// concurrency issues under heavy use - when loading and processing
// multiple images in parallel, the library will sometimes fail with
// a SIGSEGV, or a thread will hang.
type vipsImage struct {
	img *vips.ImageRef
}

func VipsInit() {
	// Make vips less noisy.
	vips.LoggingSettings(nil, vips.LogLevelError)
	vips.Startup(nil)
}

// NewVipsImage returns an image loaded and managed by the
// cgo bindings to libvips.
func NewVipsImage(src string) (imager.Image, error) {
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

func (v *vipsImage) Rotate(deg int) error {
	switch deg {
	case 90:
		v.img.Rotate(vips.Angle90)
	case 180:
		v.img.Rotate(vips.Angle180)
	case 270:
		v.img.Rotate(vips.Angle270)
	default:
		return fmt.Errorf("unsupported rotation value (%d)", deg)
	}
	return nil
}

func (v *vipsImage) Write(destFile string, mtime time.Time, w, h, q int) error {
	vimg, err := v.img.Copy()
	if err != nil {
		return err
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
			return err
		}
	}
	jp := vips.NewJpegExportParams()
	jp.Quality = q
	b, _, err := vimg.ExportJpeg(jp)
	if err != nil {
		return err
	}
	if err := os.WriteFile(destFile, b, 0644); err != nil {
		return err
	}
	return os.Chtimes(destFile, mtime, mtime)
}
