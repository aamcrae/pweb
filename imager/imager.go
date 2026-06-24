package imager

import (
	"time"
)

type RotateDegrees int

const (
	Rotate90 RotateDegrees = iota
	Rotate180
	Rotate270
)

// Image defines the interface to an image processor for an image.
type Image interface {
	Width() int
	Height() int
	Rotate(degrees RotateDegrees) error
	Write(dest string, mtime time.Time, width, height, quality int) error
}
