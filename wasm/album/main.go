package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

func main() {
	w := wasm.GetWindow()
	w.LoadStyle("/pweb/album-style.css")
	var c wasm.Comp
	defer func() {
		w.Display(c.String())
	}()
	aXml, err := wasm.GetContent(data.AlbumFile)
	if err != nil {
		w.SetTitle("No album")
		c.Wr("<h1>No photo album!</h1>")
		return
	}
	var album data.AlbumPage
	err = xml.Unmarshal(aXml, &album)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		c.Wr("<h1>Bad album data!</h1>")
		return
	}
	displayAlbum(&album, w, &c)
}

// displayAlbum generates the HTML for the album from the XML data.
func displayAlbum(a *data.AlbumPage, w *wasm.Window, c *wasm.Comp) {
	if len(a.Title) > 0 {
		w.SetTitle(a.Title)
		h1 := "<h1>" + a.Title + "</h1>"
		if len(a.Back) > 0 {
			c.Wr("<a href=\"").Wr(a.Back).Wr("\">").Wr(h1).Wr("</a>")
		} else {
			c.Wr(h1)
		}
	}
	c.Wr("<table summary=\"Album\" id=\"albumTab\">")
	for _, entry := range a.Albums {
		c.Wr("<tr><td class=\"albumName\">")
		c.Wr("<a href=\"").Wr(entry.Link).Wr("\">").Wr(entry.Title).Wr("</a>")
	}
	c.Wr("</table>")
	c.Copyright(a.Copyright)
}
