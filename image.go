package main

import (
	"log"

	"github.com/aamcrae/pweb/imager"
	"github.com/aamcrae/pweb/imager/dis"
	"github.com/aamcrae/pweb/imager/vips"
)

// Function to create an image using a particular processor
type NewImage func(src string) (imager.Image, error)

// selectImage returns a factory function for managing images
func selectImager(name string) NewImage {
	switch name {
	case "vips":
		vips.VipsInit()
		return vips.NewVipsImage
	case "dis":
		return dis.NewDisImage
	default:
		log.Fatalf("%s: Unknown imager", name)
	}
	return nil
}
