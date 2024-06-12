package main

import (
	"github.com/aamcrae/pweb/data"
	html "github.com/aamcrae/wasm"
)

func main() {
	w := html.GetWindow()
	// Try to concurrently load both album.xml and gallery.xml
	f1 := w.Fetcher(data.AlbumFile)
	f2 := w.Fetcher(data.GalleryFile)
	aData, _ := f1.Get()
	gData, _ := f2.Get()
	if len(aData) > 0 {
		RunAlbum(w, aData)
	} else if len(gData) > 0 {
		RunGallery(w, gData)
	} else {
		w.SetTitle("No album or gallery")
		h := new(html.HTML)
		w.Display(h.H1("No photo album or gallery!").String())
	}
}

// Add a copyright string
func Copyright(owner string) string {
	h := new(html.HTML)
	return h.Div(h.If(len(owner) > 0), h.Id("copyright"), "&nbsp; &copy; Copyright ", owner).String()
}
