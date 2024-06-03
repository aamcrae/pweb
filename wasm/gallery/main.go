package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

type Image struct {
	base    string
	title string
	date string
}

type Gallery struct {
	w *wasm.Window
	title string
	header string
	imagePage bool
	curPage int
	curImage int
	owner string
	images []*Image
}

func main() {
	w := wasm.GetWindow()
	w.LoadStyle("gallery-style.css")
	gXml, err := wasm.GetContent(data.GalleryFile)
	if err != nil {
		var c wasm.Comp
		w.SetTitle("No gallery!")
		c.Wr("<h1>No photo gallery!</h1>")
		w.Display(&c)
		return
	}
	w.OnResize(resized)
	var gallery data.Gallery
	err = xml.Unmarshal(gXml, &gallery)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		var c wasm.Comp
		c.Wr("<h1>Bad or no gallery data!</h1>")
		w.Display(&c)
		return
	}
	g := newGallery(&gallery, w)
	g.Display()
	w.Wait()
}

func resized(w, h int) {
	fmt.Printf("Resize callback, w = %d, h = %d\n", w, h)
}

func newGallery(xmlData *data.Gallery, w *wasm.Window) *Gallery {
	g := &Gallery{w: w, owner: xmlData.Copyright}
	if len(xmlData.Title) > 0 {
		var c wasm.Comp
		g.title = xmlData.Title
		if len(xmlData.Back) > 0 {
			c.Wr("<a href=\"").Wr(xmlData.Back).Wr("\"<h1>").Wr(g.title).Wr("</h1></a>")
		} else {
			c.Wr("<h1>").Wr(g.title).Wr("</h1>")
		}
		g.header = c.String()
	}
	for _, entry := range xmlData.Photos {
		g.images = append(g.images, &Image{base: entry.Name, title: entry.Title, date: entry.Date})
	}
	return g
}

func (g *Gallery) Display() {
	var c wasm.Comp
	g.w.SetTitle(g.title)
	c.Wr(g.header)
	c.Wr("<table>")
	c.Wr("<tr><th>Name</th><th>Title</th><th>Date</th></tr>")
	for _, entry := range g.images {
		c.Wr("<tr>")
		c.Wr("<td>").Wr(entry.base).Wr("</td>")
		c.Wr("<td>").Wr(entry.title).Wr("</td>")
		c.Wr("<td>").Wr(entry.date).Wr("</td>")
		c.Wr("</tr>")
	}
	c.Wr("</table>")
	c.Copyright(g.owner)
	g.w.Display(&c)
}
