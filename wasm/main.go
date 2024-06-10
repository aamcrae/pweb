package main

import (
	"github.com/aamcrae/pweb/data"
	h "github.com/aamcrae/wasm"
)

func main() {
	w := h.GetWindow()
	// Try to concurrently load both album.xml and gallery.xml
	f1 := h.NewFetcher(w, data.AlbumFile)
	f2 := h.NewFetcher(w, data.GalleryFile)
	aData, _ := f1.Get()
	gData, _ := f2.Get()
	if len(aData) > 0 {
		RunAlbum(w, aData)
	} else if len(gData) > 0 {
		RunGallery(w, gData)
	} else {
		w.SetTitle("No album or gallery")
		w.Display(h.H1("No photo album or gallery!")
	}
}

// Add a copyright string
func Copyright(owner string) string {
	return h.Div(h.If(len(owner) > 0), h.Id("copyright"), "&nbsp; &copy; Copyright ", owner)
}
