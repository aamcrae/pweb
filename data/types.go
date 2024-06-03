package data

import (
	"encoding/xml"
)

const AlbumFile = "album.xml"
const TemplateAlbumFile = "album-template.xml"
const GalleryFile = "gallery.xml"
const TemplateGalleryFile = "gallery-template.xml"

type Album struct {
	XMLName xml.Name `xml:"album"`
	Link    string   `xml:"link"`
	Title   string   `xml:"title,omitempty"`
	Id      string   `xml:"id,omitempty"`
}

type AlbumPage struct {
	XMLName xml.Name `xml:"albumpage"`
	Title   string   `xml:"title,omitempty"`
	Back	string   `xml:"back,omitempty"`
	Copyright     string   `xml:"copyright,omitempty"`
	Albums  []Album  `xml:"album"`
}

type Gallery struct {
	XMLName  xml.Name `xml:"gallery"`
	Title    string   `xml:"title,omitempty"`
	Back     string   `xml:"back,omitempty"`
	Copyright     string   `xml:"copyright,omitempty"`
	Download string   `xml:"download,omitempty"`
	Photos   []Photo  `xml:"photo"`
}

type Photo struct {
	XMLName     xml.Name `xml:"photo"`
	Name        string   `xml:"name"`
	Title       string   `xml:"title,omitempty"`
	Caption     string   `xml:"caption,omitempty"`
	Date        string   `xml:"date,omitempty"`
	ISO         string   `xml:"iso,omitempty"`
	Exposure    string   `xml:"exposure,omitempty"`
	Aperture    string   `xml:"aperture,omitempty"`
	FocalLength string   `xml:"length,omitempty"`
	Download    string   `xml:"download,omitempty"`
}
