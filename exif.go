package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

type ExifReader interface {
	Open(string) error
	Get(...string) string
}

// Factory function for getting a new EXIF reader
var NewExifReader func() ExifReader

// ReadExif reads the file and extracts the EXIF data from the file.
func ReadExif(srcFile string) (*Exif, error) {
	if *verbose {
		fmt.Printf("%s: reading exif\n", srcFile)
	}
	reader := NewExifReader()
	if err := reader.Open(srcFile); err != nil {
		return nil, err
	}
	var exif Exif
	exif.title = reader.Get("Iptc.Application2.ObjectName", "Iptc.Application2.Headline")
	exif.caption = reader.Get("Iptc.Application2.Caption")
	exif.exposure = reader.Get("Exif.Photo.ExposureTime")
	exif.iso = reader.Get("Exif.Photo.ISOSpeedRatings")
	exif.fstop = rational(reader.Get("Exif.Photo.FNumber"))
	exif.focal_len = rational(reader.Get("Exif.Photo.FocalLength"))
	exif.orientation = reader.Get("Exif.Image.Orientation")
	date := reader.Get("Exif.Photo.DateTimeDigitized", "Exif.Photo.DateTimeOriginal", "Exif.Image.DateTime")
	if len(date) > 0 {
		var err error
		// The date should be in ISO 8601 format, but Canon uses ':' instead of '-'
		if exif.ts, err = time.ParseInLocation(canonLayout, date, time.Local); err != nil {
			if exif.ts, err = time.ParseInLocation(canonLayout, date, time.Local); err != nil {
				fmt.Printf("Unable to parse date (%s): %v\n", date, err)
			}
		}
	}
	exif.rating = reader.Get("Xmp.xmp.Rating")
	return &exif, nil
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
