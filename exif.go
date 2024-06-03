package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kolesa-team/goexiv"
)

// date/time layouts for the EXIF date objects.
const (
	canonLayout = "2006:01:02 15:04:05"
	stdLayout   = "2006-01-02 15:04:05"
)

// Exif holds the EXIF date read from a file.
type Exif struct {
	title       string
	caption     string
	orientation string
	ts          time.Time
	rating      string
	iso         string
	exposure    string
	fstop       string
	focal_len   string
}

// ReadExif reads the file and extracts the EXIF data from the file.
// goexiv (a cgo binding to libexiv2) is used, but sometimes it seems
// this binding doesn't handle concurrency reliably, and will sometimes crash unexpectedly.
func ReadExif(srcFile string) (*Exif, error) {
	if *verbose {
		fmt.Printf("%s: reading exif\n", srcFile)
	}
	var exif Exif
	idata, err := os.ReadFile(srcFile)
	if err != nil {
		return nil, err
	}
	// Read exif and extract relevant tags.
	img, err := goexiv.OpenBytes(idata)
	if err != nil {
		// Assume no exif
		log.Printf("%s: exif read (%v)", srcFile, err)
		return &exif, nil
	} else {
		err := img.ReadMetadata()
		if err != nil {
			// Unable to parse exif
			return nil, err
		}
		exif.title = getIptc(img, "Iptc.Application2.ObjectName", "Iptc.Application2.Headline")
		exif.caption = getIptc(img, "Iptc.Application2.Caption")
		exif.exposure = getExif(img, "Exif.Photo.ExposureTime")
		exif.iso = getExif(img, "Exif.Photo.ISOSpeedRatings")
		exif.fstop = rational(getExif(img, "Exif.Photo.FNumber"))
		exif.focal_len = rational(getExif(img, "Exif.Photo.FocalLength"))
		exif.orientation = getExif(img, "Exif.Image.Orientation")
		date := getExif(img, "Exif.Photo.DateTimeDigitized", "Exif.Photo.DateTimeOriginal", "Exif.Image.DateTime")
		if len(date) > 0 {
			// The date should be in ISO 8601 format, but Canon uses ':' instead of '-'
			if exif.ts, err = time.ParseInLocation(canonLayout, date, time.Local); err != nil {
				if exif.ts, err = time.ParseInLocation(canonLayout, date, time.Local); err != nil {
					fmt.Printf("Unable to parse date (%s): %v\n", date, err)
				}
			}
		}
		xmp := img.GetXmpData()
		if r, err := xmp.FindKey("Xmp.xmp.Rating"); err != nil || r == nil {
			exif.rating = ""
		} else {
			exif.rating = r.String()
		}
	}
	return &exif, nil
}

func getIptc(img *goexiv.Image, keys ...string) string {
	idata := img.GetIptcData()
	for _, s := range keys {
		if v, err := idata.FindKey(s); err == nil && v != nil {
			return v.String()
		}
	}
	return ""
}

func getExif(img *goexiv.Image, keys ...string) string {
	edata := img.GetExifData()
	for _, s := range keys {
		if v, err := edata.FindKey(s); err == nil && v != nil {
			return v.String()
		}
	}
	return ""
}

// Convert rational to FP
func rational(in string) string {
	if len(in) == 0 {
		return ""
	}
	var v1, v2 float64
	n, err := fmt.Sscanf(in, "%f/%f", &v1, &v2)
	if err != nil || n != 2 {
		fmt.Printf("Failed to convert EXIF ration value %s\n", in)
		return in
	} else {
		s := strconv.FormatFloat(v1/v2, 'f', 2, 64)
		s = strings.TrimRight(s, "0")
		// Check for no trailing digits
		last := len(s) - 1
		if s[last] == '.' {
			s = s[:last]
		}
		return s
	}
}
