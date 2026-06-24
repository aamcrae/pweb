package dis

import (
	"image"
	"os"
	"time"

	"github.com/aamcrae/pweb/imager"
	"github.com/disintegration/imaging"
)

type disImage struct {
	img image.Image
}

func NewDisImage(src string) (imager.Image, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return &disImage{img: img}, nil
}

func (d *disImage) Width() int {
	return d.img.Bounds().Max.X
}

func (d *disImage) Height() int {
	return d.img.Bounds().Max.Y
}

func (d *disImage) Rotate(deg imager.RotateDegrees) error {
	switch deg {
	case imager.Rotate90:
		d.img = imaging.Rotate90(d.img)
	case imager.Rotate180:
		d.img = imaging.Rotate180(d.img)
	case imager.Rotate270:
		d.img = imaging.Rotate270(d.img)
	}
	return nil
}

func (d *disImage) Write(destFile string, mtime time.Time, w, h, q int) error {
	// scale it down to fit within the width & height
	xr := float64(w) / float64(d.Width())
	yr := float64(h) / float64(d.Height())
	var img image.Image
	if xr < 1 || yr < 1 {
		// Image is larger than requested size, so scale it down
		if xr < yr {
			img = imaging.Resize(d.img, w, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(d.img, 0, h, imaging.Lanczos)
		}
	} else {
		img = d.img
	}
	if err := imaging.Save(img, destFile, imaging.JPEGQuality(q)); err != nil {
		return err
	}
	return os.Chtimes(destFile, mtime, mtime)
}
