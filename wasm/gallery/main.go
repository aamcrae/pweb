package main

import (
	"encoding/xml"
	"fmt"

	"syscall/js"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

type Image struct {
	base       string
	title      string
	date       string
	thumbEntry string
	imagePage  string
	original   data.Size
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
		w.Display(c.String())
		return
	}
	var gallery data.Gallery
	err = xml.Unmarshal(gXml, &gallery)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		var c wasm.Comp
		c.Wr("<h1>Bad or no gallery data!</h1>")
		w.Display(c.String())
		return
	}
	g := newGallery(&gallery, w)
	// Add some callbacks
	w.OnResize(g.Resize)
	w.AddJSFunction("showPict", g.ShowPic)
	w.AddJSFunction("showThumbs", g.ShowThumbs)
	g.ShowPage()
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
		img := &Image{base: entry.Name, title: entry.Title, date: entry.Date, original: entry.Original}
		var ct wasm.Comp
		ct.Wr(fmt.Sprintf("<div class=\"holder\"><div id=\"slide%d\" class=\"slideshow\">", i))
		ct.Wr(fmt.Sprintf("<a onclick=\"return showPict(%d)\" href=\"#\">", i))
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
		g.ShowPage()
	}
}

func (g *Gallery) ShowPic(this js.Value, p []js.Value) any {
	g.curImage = p[0].Int()
	fmt.Printf("Showing picture %d\n", p[0].Int())
	img := g.images[g.curImage]
	if len(img.imagePage) == 0 {
		g.BuildPict(g.curImage)
	}
	g.w.Display(img.imagePage)
	return js.ValueOf(false)
}

func (g *Gallery) ShowThumbs(this js.Value, p []js.Value) any {
	fmt.Printf("Showing thumbnails\n")
	g.curImage = p[0].Int()
	g.ShowPage()
	return js.ValueOf(false)
}

func (g *Gallery) ShowPage() {
	var c wasm.Comp
	g.w.SetTitle(g.title)
	c.Wr(g.header)
	g.rows, g.cols = g.tableSize()
	pageNo := g.curImage / (g.rows * g.cols)
	c.Wr("<table>")
	i := pageNo * g.rows * g.cols
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
	g.w.Display(c.String())
}

func (g *Gallery) BuildPict(index int) {
	img := g.images[index]
	c := new(wasm.Comp)
	g.LinkToPict(c, "prev", index-1)
	c.Wr("<div id=\"home\"><a onclick=\"return showThumbs(").Wr(index).Wr(")\" href=\"#\">back to index</a></div>")
	g.LinkToPict(c, "next", index+1)
	c.Wr("<h1>")
	if len(img.title) > 0 {
		c.Wr(img.title)
	} else {
		c.Wr(g.title)
	}
	c.Wr("</h1>")
	c.Wr("<div id=\"mainimage\"><img src=\"").Wr(img.base).Wr("\" alt=\"").Wr(img.title).Wr("\"></div>")
	// Show image properties etc.
	c.Wr("<div class=\"properties\"><table summary=\"image properties\" border=\"0\">")
	g.Property(c, "Date", img.date)
	g.Property(c, "Filename", img.base)
	if img.original.Width != 0 && img.original.Height != 0 {
		g.Property(c, "Original resolution", fmt.Sprintf("%d x %d", img.original.Width, img.original.Height))
	}
	c.Wr("</table></div>")
	c.Copyright(g.owner)
	img.imagePage = c.String()
}

func (g *Gallery) Property(c *wasm.Comp, n, val string) {
	if val != "" {
		c.Wr("<tr><td class=\"exifname\">").Wr(n).Wr("</td><td class=\"exifdata\">").Wr(val).Wr("</td></tr>")
	}
}

func (g *Gallery) LinkToPict(c *wasm.Comp, n string, index int) {
	c.Wr("<div id=\"").Wr(n).Wr("\">")
	if index < 0 || index == len(g.images) {
		c.Wr("&nbsp")
	} else {
		c.Wr("<a onclick=\"return showPict(").Wr(index).Wr(")\" href=\"#\">")
		if len(g.images[index].title) == 0 {
			c.Wr(g.images[index].base)
		} else {
			c.Wr(g.images[index].title)
		}
		c.Wr("</a>")
	}
	c.Wr("</div>")
}

func (g *Gallery) tableSize() (int, int) {
	return max(g.w.Width/(g.tw+33), 1), max((g.w.Height-200)/(g.th+33), 1)
}
