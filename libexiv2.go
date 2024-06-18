package main

import (
	"os"
	"strings"

	"github.com/kolesa-team/goexiv"
)

type Exiv2 struct {
	img *goexiv.Image
}

// Exiv2Reader creates and returns a new EXIF reader.
// goexiv (a cgo binding to libexiv2) is used, but sometimes it seems
// this binding doesn't handle concurrency reliably, and will sometimes crash unexpectedly.
func Exiv2Reader() ExifReader {
	return &Exiv2{}
}

func (r *Exiv2) Open(file string) error {
	idata, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	// Read exif and extract relevant tags.
	if r.img, err = goexiv.OpenBytes(idata); err != nil {
		// No exif. Allow this.
		return err
	}
	err = r.img.ReadMetadata()
	if err != nil {
		// Unable to parse exif
		return err
	}
	return nil
}

func (r *Exiv2) Get(keys ...string) string {
	for _, k := range keys {
		b, _, _ := strings.Cut(k, ".")
		switch b {
		case "Iptc":
			idata := r.img.GetIptcData()
			if v, err := idata.FindKey(k); err == nil && v != nil {
				return v.String()
			}
			return ""
		case "Exif":
			edata := r.img.GetExifData()
			if v, err := edata.FindKey(k); err == nil && v != nil {
				return v.String()
			}
			return ""
		case "Xmp":
			xmp := r.img.GetXmpData()
			if v, err := xmp.FindKey(k); err == nil && v != nil {
				return v.String()
			}
			return ""
		default:
			return ""
		}
	}
	return ""
}
