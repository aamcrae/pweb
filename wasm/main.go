package main

import (
	"github.com/aamcrae/pweb/data"
)

func main() {
	w := GetWindow()
	// Try to concurrently load both album.xml and gallery.xml
	f1 := NewFetcher(w, data.AlbumFile)
	f2 := NewFetcher(w, data.GalleryFile)
	aData, _ := f1.Get()
	gData, _ := f2.Get()
	if len(aData) > 0 {
		RunAlbum(w, aData)
	} else if len(gData) > 0 {
		RunGallery(w, gData)
	} else {
		w.SetTitle("No album or gallery")
		w.Display("<h1>No photo album or gallery!</h1>")
	}
}

// Add a copyright string
func Copyright(c *Comp, owner string) {
	if len(owner) > 0 {
		c.Wr("<div id=\"copyright\">&nbsp; &copy; Copyright ").Wr(owner).Wr("</div>")
	}
}
