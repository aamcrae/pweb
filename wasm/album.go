package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
)

func RunAlbum(w *Window, ax []byte) {
	w.LoadStyle("/pweb/album-style.css")
	var c Comp
	defer func() {
		w.Display(c.String())
	}()
	var album data.AlbumPage
	err := xml.Unmarshal(ax, &album)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		c.Wr("<h1>Bad album data!</h1>")
		return
	}
	displayAlbum(&album, w, &c)
}

// displayAlbum generates the HTML for the album from the XML data.
func displayAlbum(a *data.AlbumPage, w *Window, c *Comp) {
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
	Copyright(c, a.Copyright)
}
