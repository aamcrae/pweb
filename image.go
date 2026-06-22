package main

import (
	"time"
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
