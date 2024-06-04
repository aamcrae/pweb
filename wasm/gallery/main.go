package main

import (
	"encoding/xml"
	"fmt"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

type Image struct {
	base       string
	title      string
	date       string
	thumbEntry string
}

type Gallery struct {
	w          *wasm.Window
	title      string
	header     string
	imagePage  bool
	curPage    int
	curImage   int
	owner      string
	th, tw     int
	pw, ph     int
	iw, ih     int
	rows, cols int
	images     []*Image
}

func main() {
	w := wasm.GetWindow()
	w.LoadStyle("/pweb/gallery-style.css")
	gXml, err := wasm.GetContent(data.GalleryFile)
	if err != nil {
		var c wasm.Comp
		w.SetTitle("No gallery!")
		c.Wr("<h1>No photo gallery!</h1>")
		w.Display(&c)
		return
	}
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
	w.OnResize(g.Resize)
	g.Display()
	w.Wait()
}

func newGallery(xmlData *data.Gallery, w *wasm.Window) *Gallery {
	g := &Gallery{w: w,
		owner: xmlData.Copyright,
		tw:    xmlData.Thumb.Width,
		th:    xmlData.Thumb.Height,
		pw:    xmlData.Preview.Width,
		ph:    xmlData.Preview.Height,
		iw:    xmlData.Image.Width,
		ih:    xmlData.Image.Height,
	}
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
	for i, entry := range xmlData.Photos {
		img := &Image{base: entry.Name, title: entry.Title, date: entry.Date}
		var ct wasm.Comp
		ct.Wr(fmt.Sprintf("<div class=\"holder\"><div id=\"slide%d\" class=\"slideshow\">", i))
		ct.Wr(fmt.Sprintf("<a onclick=\"return showPic(%d)\" href=\"h/%s\" target=\"_top\">", i, img.base))
		ct.Wr("<img src=\"t/").Wr(img.base).Wr("\" title=\"").Wr(img.title).Wr("\">")
		ct.Wr("</a>")
		if len(img.title) > 0 {
			ct.Wr("<div class=thumbName>").Wr(img.title).Wr("</div>")
		}
		ct.Wr("</div> </div>")
		img.thumbEntry = ct.String()
		g.images = append(g.images, img)
	}
	return g
}

func (g *Gallery) Resize() {
	rows, cols := g.tableSize()
	if rows != g.rows || cols != g.cols {
		fmt.Printf("Resize callback, cols = %d, rows = %d\n", cols, rows)
		g.Display()
	}
}

func (g *Gallery) Display() {
	var c wasm.Comp
	g.w.SetTitle(g.title)
	c.Wr(g.header)
	g.rows, g.cols = g.tableSize()
	c.Wr("<table>")
	i := 0
	for y := 0; y < g.cols; y++ {
		c.Wr("<tr>")
		for x := 0; x < g.rows; x++ {
			c.Wr("<td>")
			if i < len(g.images) {
				c.Wr(g.images[i].thumbEntry)
				i++
			}
			c.Wr("</td>")
		}
		c.Wr("</tr>")
	}
	c.Wr("</table>")
	c.Copyright(g.owner)
	g.w.Display(&c)
}

func (g *Gallery) tableSize() (int, int) {
	return g.w.Width / (g.tw + 33), (g.w.Height - 200) / (g.th + 33)
}
