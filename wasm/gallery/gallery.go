package main

import (
	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

type Image struct {
	base    string
	caption string
}

type Gallery struct {
	images []Image
}

func displayGallery(g *data.Gallery, w *wasm.Window, p *wasm.Page) {
	if len(g.Title) > 0 {
		w.SetTitle(g.Title)
		if len(g.Back) > 0 {
			p.Wr("<a href=\"").Wr(g.Back).Wr("\"<h1>").Wr(g.Title).Wr("</h1></a>")
		} else {
			p.Wr("<h1>").Wr(g.Title).Wr("</h1>")
		}
	}
	p.Wr("<table>")
	p.Wr("<tr><th>Name</th><th>Title</th><th>Date</th></tr>")
	for _, entry := range g.Photos {
		p.Wr("<tr>")
		p.Wr("<td>").Wr(entry.Name).Wr("</td>")
		p.Wr("<td>").Wr(entry.Title).Wr("</td>")
		p.Wr("<td>").Wr(entry.Date).Wr("</td>")
		p.Wr("</tr>")
	}
	p.Wr("</table>")
	p.Copyright(g.Copyright)
	w.Display(p)
}
