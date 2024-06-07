package main

import (
	"github.com/aamcrae/pweb/data"
)

func main() {
	w := GetWindow()
	if aData, err := GetContent(data.AlbumFile); err == nil {
		RunAlbum(w, aData)
	} else if gData, err := GetContent(data.GalleryFile); err == nil {
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
