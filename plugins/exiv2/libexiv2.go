package exiv2

import (
	"os"
	"strings"

	"github.com/kolesa-team/goexiv"
)

type exiv2 struct {
	img *goexiv.Image
}

// goexiv (a cgo binding to libexiv2) is used, but sometimes it seems
// this binding doesn't handle concurrency reliably, and will sometimes crash unexpectedly.
func Exiv2Open(file string) (*exiv2, error) {
	idata, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Read exif and extract relevant tags.
	img, err := goexiv.OpenBytes(idata)
	if err != nil {
		// No exif. Allow this.
		return nil, err
	}
	if err = img.ReadMetadata(); err != nil {
		// Unable to parse exif
		return nil, err
	}
	return &exiv2{ img: img}, nil
}

func (r *exiv2) Get(keys ...string) string {
	for _, k := range keys {
		b, _, _ := strings.Cut(k, ".")
		switch b {
		case "Iptc":
			idata := r.img.GetIptcData()
			if v, err := idata.FindKey(k); err == nil && v != nil && v.String() != "" {
				return v.String()
			}
		case "Exif":
			edata := r.img.GetExifData()
			if v, err := edata.FindKey(k); err == nil && v != nil {
				return v.String()
			}
		case "Xmp":
			xmp := r.img.GetXmpData()
			if v, err := xmp.FindKey(k); err == nil && v != nil {
				return v.String()
			}
		}
	}
	return ""
}
