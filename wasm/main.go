package main

import (
	"sync"

	"github.com/aamcrae/pweb/data"
)

func main() {
	// Try to concurrently load both album.xml and gallery.xml
	var wg sync.WaitGroup
	wg.Add(2)
	w := GetWindow()
	var aData, gData []byte
	go func() {
		aData, _ = w.GetContent(data.AlbumFile)
		wg.Done()
	}()
	go func() {
		gData, _ = w.GetContent(data.GalleryFile)
		wg.Done()
	}()
	 wg.Wait()
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
