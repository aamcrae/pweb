package shared

import (
	"encoding/xml"
)

type Album struct {
	XMLName xml.Name `xml:"album" json:"-"`
	Link    string   `xml:"link" json:"link"`
	Title   string   `xml:"title,omitempty" json:"title,omitempty"`
	Id      string   `xml:"id,omitempty" json:"id,omitempty"`
}

type Size struct {
	Width  int `xml:"width" json:"width,omitzero"`
	Height int `xml:"height" json:"height,omitzero"`
}

type AlbumPage struct {
	XMLName   xml.Name `xml:"albumpage" json:"-"`
	Title     string   `xml:"title,omitempty" json:"title,omitempty"`
	Back      string   `xml:"back,omitempty" json:"back,omitempty"`
	Copyright string   `xml:"copyright,omitempty" json:"copyright,omitempty"`
	Albums    []Album  `xml:"album" json:"albums,omitempty"`
}

type Gallery struct {
	XMLName   xml.Name `xml:"gallery" json:"-"`
	Title     string   `xml:"title,omitempty" json:"title,omitempty"`
	Back      string   `xml:"back,omitempty" json:"back,omitempty"`
	Copyright string   `xml:"copyright,omitempty" json:"copyright,omitempty"`
	Download  string   `xml:"download,omitempty" json:"download,omitempty"`
	Thumb     Size     `xml:"thumb" json:"thumb"`
	Preview   Size     `xml:"preview" json:"preview"`
	Image     Size     `xml:"image" json:"image"`
	Photos    []Photo  `xml:"photo" json:"photos,omitempty"`
}

type Photo struct {
	XMLName     xml.Name `xml:"photo" json:"-"`
	Name        string   `xml:"name" json:"name"`
	Filename    string   `xml:"filename" json:"filename"`
	Original    Size     `xml:"original" json:"original"`
	Title       string   `xml:"title,omitempty" json:"title,omitempty"`
	Caption     string   `xml:"caption,omitempty" json:"caption,omitempty"`
	Date        string   `xml:"date,omitempty" json:"date,omitempty"`
	ISO         string   `xml:"iso,omitempty" json:"iso,omitempty"`
	Exposure    string   `xml:"exposure,omitempty" json:"exposure,omitempty"`
	Aperture    string   `xml:"aperture,omitempty" json:"aperture,omitempty"`
	FocalLength string   `xml:"length,omitempty" json:"length,omitempty"`
	Download    string   `xml:"download,omitempty" json:"download,omitempty"`
}
