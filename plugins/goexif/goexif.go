package goexif

import (
	_ "fmt"
	"github.com/dsoprea/go-exif/v3"
)

type goexif struct {
	tagMap map[string]*exif.ExifTag
}

func GoExifOpen(file string) (*goexif, error) {
	edata, err := exif.SearchFileAndExtractExif(file)
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]*exif.ExifTag)
	// Extract all tags
	tags, _, err := exif.GetFlatExifDataUniversalSearch(edata, nil, true)
	if err != nil {
		return nil, err
	}
	for _, t := range tags {
		tagMap[t.TagName] = &t
	}
	return &goexif{tagMap: tagMap}, nil
}

func (r *goexif) Get(keys ...string) string {
	for _, k := range keys {
		v, ok := r.tagMap[k]
		if ok {
			return v.Formatted
		}
	}
	return ""
}
