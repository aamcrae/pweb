package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/aamcrae/pweb/data"
	h "github.com/aamcrae/wasm"
)

func RunAlbum(w *h.Window, ax []byte) {
	w.LoadStyle("/pweb/album-style.css")
	var album data.AlbumPage
	err := xml.Unmarshal(ax, &album)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		w.Display(h.H1("Bad album data!"))
		return
	}
	w.Display(displayAlbum(w, &album))
	if len(album.Back) > 0 {
		w.OnKey(func(key string) {
			switch key {
			case "ArrowLeft", "ArrowUp":
				w.Goto(album.Back)
			}
		})
		w.OnSwipe(func(d h.Direction) bool {
			if d == h.Down {
				w.Goto(album.Back)
				return true
			}
			return false
		})
	}
	w.Wait()
}

// displayAlbum generates the HTML for the album from the XML data.
func displayAlbum(w *h.Window, a *data.AlbumPage) string {
	var c strings.Builder
	if len(a.Title) > 0 {
		w.SetTitle(a.Title)
		h1 := h.H1(a.Title)
		
		if len(a.Back) > 0 {
			c.WriteString(h.A(h.Href(a.Back), h1))
		} else {
			c.WriteString(h1)
		}
	}
	c.WriteString(h.Table(h.Open(), h.Summary("Album"), h.Id("albumTab")))
	for _, entry := range a.Albums {
		c.WriteString(h.Tr(h.Td(h.Class("albumName"), h.A(h.Href(entry.Link), entry.Title))))
	}
	c.WriteString(h.Table(h.Close()))
	c.WriteString(Copyright(a.Copyright))
	return c.String()
}
