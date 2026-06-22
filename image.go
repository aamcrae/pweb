package main

import (
	"time"

	"github.com/aamcrae/pweb/imaging/dis"
	"github.com/aamcrae/pweb/imaging/vips"
)

// Image defines the interface to an image processor for an image.
type Image interface {
	Width() int
	Height() int
	Rotate(degrees int) error // Should only be 90, 180, 270
	Write(dest string, mtime time.Time, width, height, quality int) error
}

// Function to create an image using a particular processor
type NewImage func(src string) (Image, error)

func selectImager(name string) NewImage {
	switch *imager {
	case "vips":
		return func(s string) (Image, error) {
				return vips.NewVipsImage(s)
			}
		vips.vipsInit()
	case "dis":
		return func(s string) (Image, error) {
			return dis.NewDisImage(s)
		}
	default:
		log.Fatalf("%s: Unknown imager", *imager)
	}
}
