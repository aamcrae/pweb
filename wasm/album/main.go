package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

func main() {
	w := wasm.GetWindow()
	w.LoadStyle("album-style.css")
	p := wasm.NewPage()
	defer w.Display(p)
	aXml, err := wasm.GetContent(data.AlbumFile)
	if err != nil {
		w.SetTitle("No album")
		p.Wr("<h1>No photo album!</h1>")
		return
	}
	var album data.AlbumPage
	err = xml.Unmarshal(aXml, &album)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		p.Wr("<h1>Bad album data!</h1>")
		return
	}
	displayAlbum(&album, w, p)
}

func displayAlbum(a *data.AlbumPage, w *wasm.Window, p *wasm.Page) {
	if len(a.Title) > 0 {
		w.SetTitle(a.Title)
		h1 := "<h1>" + a.Title + "</h1>"
		if len(a.Back) > 0 {
			p.Wr("<a href=\"").Wr(a.Back).Wr("\">").Wr(h1).Wr("</a>")
		} else {
			p.Wr(h1)
		}
	}
	p.Wr("<table summary=\"Album\" id=\"albumTab\">")
	for _, entry := range a.Albums {
		p.Wr("<tr><td class=\"albumName\">")
		p.Wr("<a href=\"").Wr(entry.Link).Wr("\">").Wr(entry.Title).Wr("</a>")
	}
	p.Wr("</table>")
	p.Copyright(a.Copyright)
}
