package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

func main() {
	w := wasm.GetWindow()
	p := wasm.NewPage()
	w.LoadStyle("gallery-style.css")
	aXml, err := wasm.GetContent(data.GalleryFile)
	if err != nil {
		w.SetTitle("No gallery!")
		p.Wr("<h1>No photo gallery!</h1>")
		w.Display(p)
		return
	}
	w.OnResize(resized)
	var gallery data.Gallery
	err = xml.Unmarshal(aXml, &gallery)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		p.Wr("<h1>Bad gallery data!</h1>")
		return
	}
	displayGallery(&gallery, w, p)
	w.Wait()
}

func resized(w, h int) {
	fmt.Printf("Resize callback, w = %d, h = %d\n", w, h)
}
