package main

import (
	"encoding/xml"

	"github.com/aamcrae/pweb/data"
	html "github.com/aamcrae/wasm"
)

func RunAlbum(w *html.Window, ax []byte) {
	w.LoadStyle("/pweb/album-style.css")
	var album data.AlbumPage
	err := xml.Unmarshal(ax, &album)
	if err != nil {
		w.Display(new(html.HTML).H1("Bad album data!").String())
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
		w.OnSwipe(func(d html.Direction) bool {
			if d == html.Down {
				w.Goto(album.Back)
				return true
			}
			return false
		})
	}
	w.Wait()
}

// displayAlbum generates the HTML for the album from the XML data.
func displayAlbum(w *html.Window, a *data.AlbumPage) string {
	h := new(html.HTML)
	if len(a.Title) > 0 {
		w.SetTitle(a.Title)
		title := h.H1(a.Title).String()

		if len(a.Back) > 0 {
			h.Wr(h.A(h.Href(a.Back), title))
		} else {
			h.Wr(title)
		}
	}
	h.Wr(h.Table(h.Open(), h.Summary("Album"), h.Id("albumTab")))
	for _, entry := range a.Albums {
		h.Wr(h.Tr(h.Td(h.Class("albumName"), h.A(h.Href(entry.Link), entry.Title))))
	}
	h.Wr(h.Table(h.Close()))
	h.Wr(Copyright(a.Copyright))
	return h.String()
}
