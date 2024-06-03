package main

import (
	"image"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

type disImage struct {
	Image
	img image.Image
}

func NewDisImage(src string) (Image, error) {
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

func (d *disImage) Rotate(deg RotateDegrees) error {
	switch deg {
	case Rotate90:
		d.img = imaging.Rotate90(d.img)
	case Rotate180:
		d.img = imaging.Rotate180(d.img)
	case Rotate270:
		d.img = imaging.Rotate270(d.img)
	}
	return nil
}

func (d *disImage) Write(destFile string, mtime time.Time, w, h, q int) {
	// scale it down to fit within the width & height
	xr := float64(w) / float64(d.Width())
	yr := float64(h) / float64(d.Height())
	var img image.Image
	if xr < 1 || yr < 1 {
		// Image is larger than requested size, so scale it down
		if xr < yr {
			img = imaging.Resize(d.img, 0, h, imaging.Lanczos)
		} else {
			img = imaging.Resize(d.img, w, 0, imaging.CatmullRom)
		}
	} else {
		img = d.img
	}
	err := imaging.Save(img, destFile, imaging.JPEGQuality(q))
	if err != nil {
		log.Fatalf("%s: save: %v", destFile, err)
	}
	if err := os.Chtimes(destFile, mtime, mtime); err != nil {
		log.Fatalf("%s: chtime: %v", destFile, err)
	}
}
