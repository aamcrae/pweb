package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/aamcrae/pweb/data"
)

func RunAlbum(w *Window, ax []byte) {
	w.LoadStyle("/pweb/album-style.css")
	var album data.AlbumPage
	err := xml.Unmarshal(ax, &album)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		w.Display(H1("Bad album data!"))
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
		w.OnSwipe(func(d Direction) bool {
			if d == Down {
				w.Goto(album.Back)
				return true
			}
			return false
		})
	}
	w.Wait()
}

// displayAlbum generates the HTML for the album from the XML data.
func displayAlbum(w *Window, a *data.AlbumPage) string {
	var c strings.Builder
	if len(a.Title) > 0 {
		w.SetTitle(a.Title)
		h1 := H1(a.Title)
		
		if len(a.Back) > 0 {
			c.WriteString(A(Href(a.Back), h1))
		} else {
			c.WriteString(h1)
		}
	}
	c.WriteString(Table(Open(), Summary("Album"), Id("albumTab")))
	for _, entry := range a.Albums {
		c.WriteString(Tr(Td(Class("albumName"), A(Href(entry.Link), entry.Title))))
	}
	c.WriteString(Table(Close()))
	c.WriteString(Copyright(a.Copyright))
	return c.String()
}
