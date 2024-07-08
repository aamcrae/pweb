package main

import (
	_ "fmt"
	"github.com/dsoprea/go-exif/v3"
)

type Goexif struct {
	tagMap map[string]*exif.ExifTag
}

// GoexifReader creates and returns a new EXIF reader.
func GoexifReader() ExifReader {
	return &Goexif{}
}

func (r *Goexif) Open(file string) error {
	edata, err := exif.SearchFileAndExtractExif(file)
	if err != nil {
		return err
	}
	r.tagMap = make(map[string]*exif.ExifTag)
	// Extract all tags
	tags, _, err := exif.GetFlatExifDataUniversalSearch(edata, nil, true)
	if err != nil {
		return err
	}
	for _, t := range tags {
		r.tagMap[t.TagName] = &t
	}
	return nil
}

func (r *Goexif) Get(keys ...string) string {
	for _, k := range keys {
		v, ok := r.tagMap[k]
		if ok {
			return v.Formatted
		}
	}
	return ""
}
