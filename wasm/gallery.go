package main

import (
	"encoding/xml"
	"strings"

	"syscall/js"

	"github.com/aamcrae/pweb/data"
	h "github.com/aamcrae/wasm"
)

const (
	thumbOn  = "slideshowon"
	thumbOff = "slideshow"
)

const hSpace = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"

/*
 * Image holds the data for a single photo
 */
type Image struct {
	name       string    // Base filename (that may not be unique)
	filename   string    // Unique filename that may include appended directory names
	title      string    // Headline or title
	date       string    // Date photo was taken
	thumbEntry string    // The HTML used to display the thumbnail
	imagePage  string    // The HTML used to display the full sized image
	download   string    // If set, the file for download
	original   data.Size // The original image's resolution
	exposure   string    // EXIF data
	aperture   string
	iso        string
	flen       string
}

// Gallery holds the collection of images that form a photo gallery
type Gallery struct {
	w          *h.Window
	title      string   // Title of gallery
	header     string   // HTML of title of page
	back string // Referring link
	imagePage  bool     // If set, displaying the full sized image, otherwise showing thumbnails
	firstImage int      // The index of the first thumbnail displayed on the page
	lastImage  int      // The index of the last thumbnail displayed on the page
	curImage   int      // The current image
	owner      string   // If set, display a copyright notice
	th, tw     int      // Size of thumbnail image
	pw, ph     int      // Size of preview image
	iw, ih     int      // Size of full image
	rows, cols int      // Number of thumbnail rows and columns being displayed
	images     []*Image // slice of images in the gallery
}

func RunGallery(w *h.Window, gx []byte) {
	w.LoadStyle("/pweb/gallery-style.css")
	var gallery data.Gallery
	err := xml.Unmarshal(gx, &gallery)
	if err != nil {
		w.Display(h.H1("Bad or no gallery data!"))
		return
	}
	g := newGallery(&gallery, w)
	// Add some callbacks
	w.OnResize(g.Resize)
	w.OnKey(g.KeyPress)
	w.OnSwipe(g.Swipe)
	w.AddJSFunction("showPict", g.ShowPic)
	w.AddJSFunction("showThumbs", g.ShowThumbs)
	g.ShowPage()
	w.Wait()
}

// newGallery creates a new gallery from the XML data provided.
func newGallery(xmlData *data.Gallery, w *h.Window) *Gallery {
	g := &Gallery{w: w,
		title: xmlData.Title,
		back: xmlData.Back,
		owner: xmlData.Copyright,
		tw:    xmlData.Thumb.Width,
		th:    xmlData.Thumb.Height,
		pw:    xmlData.Preview.Width,
		ph:    xmlData.Preview.Height,
		iw:    xmlData.Image.Width,
		ih:    xmlData.Image.Height,
	}
	if g.title == "" {
		g.title = "Gallery"
	}
	g.header = g.HeaderDownload(g.title, g.back, xmlData.Download)
	for i, entry := range xmlData.Photos {
		img := &Image{name: entry.Name,
			filename: entry.Filename,
			title:    entry.Title,
			date:     entry.Date,
			download: entry.Download,
			original: entry.Original,
			aperture: entry.Aperture,
			exposure: entry.Exposure,
			iso:      entry.ISO,
			flen:     entry.FocalLength}
		img.thumbEntry = 
			h.Div(h.Class("holder"),
				h.Div(h.Class("slideshow"), h.Id(h.Text("slide", i)),
				h.A(h.Onclick(h.Text("return showPict(", i, ")")),
					h.Href("#"),
					h.Img(h.Title(img.title), h.Src(h.Text("t/", img.filename)))),
				h.Div(h.If(len(img.title) > 0), h.Class("thumbName"), img.title)))
		g.images = append(g.images, img)
	}
	// Install some style elements now that we know the thumbnail sizes
	g.w.AddStyle(h.Text(".holder {width:", g.tw + 10, "px;height:", g.th + 30, "px} .thumbName{width:", g.tw + 10, "px}"))
	return g
}

// Resize will redisplay the thumbnail page when the window is resized.
func (g *Gallery) Resize() {
	if !g.imagePage {
		cols, rows := g.tableSize()
		if rows != g.rows || cols != g.cols {
			g.ShowPage()
		}
	}
}

// Swipe handles a touch swipe action
func (g *Gallery) Swipe(d h.Direction) bool {
	if g.imagePage {
		switch d {
		case h.Down:
			g.SelectThumb(g.curImage)
		case h.Right:
			g.ImageDisplay(g.curImage - 1)
		case h.Left:
			g.ImageDisplay(g.curImage + 1)
		default:
			return false
		}
	} else {
		perPage := g.rows * g.cols
		switch d {
		case h.Down:
			if len(g.back) > 0 {
				g.w.Goto(g.back)
			}
		case h.Right:
			g.SelectThumb(g.curImage - perPage)
		case h.Left:
			g.SelectThumb(g.curImage + perPage)
		default:
			return false
		}
	}
	return true
}

// KeyPress handles keyboard shortcuts
func (g *Gallery) KeyPress(key string) {
	if g.imagePage {
		// Image page
		switch key {
		case "Home":
			g.ImageDisplay(0)
		case "End":
			g.ImageDisplay(len(g.images) - 1)
		case "ArrowRight", "PageDown":
			g.ImageDisplay(g.curImage + 1)
		case "ArrowLeft", "PageUp":
			g.ImageDisplay(g.curImage - 1)
		case "ArrowUp":
			g.SelectThumb(g.curImage)
		}
	} else {
		// Thumbnail page
		switch key {
		case "Enter":
			g.ImageDisplay(g.curImage)
		case "Home":
			g.SelectThumb(0)
		case "End":
			g.SelectThumb(len(g.images) - 1)
		case "ArrowRight":
			g.SelectThumb(g.curImage + 1)
		case "ArrowLeft":
			g.SelectThumb(g.curImage - 1)
		case "ArrowUp":
			g.SelectThumb(g.curImage - g.cols)
		case "ArrowDown":
			g.SelectThumb(g.curImage + g.cols)
		case "PageDown":
			g.SelectThumb(g.curImage + g.rows * g.cols)
		case "PageUp":
			g.SelectThumb(g.curImage - g.rows * g.cols)
		}
	}
}

// SelectThumb sets the current image and displays the thumbnail page
func (g *Gallery) SelectThumb(index int) {
	if index < 0 {
		index = len(g.images) - 1
	} else if index >= len(g.images) {
		index = 0
	}
	if index == g.curImage && !g.imagePage {
		return
	}
	index, g.curImage = g.curImage, index
	// Check if same page.
	if !g.imagePage && g.curImage >= g.firstImage && g.curImage <= g.lastImage {
		g.updateThumb(index, thumbOff)
		g.updateThumb(g.curImage, thumbOn)
	} else {
		g.ShowPage()
	}
}

// ShowPic is a callback from a javascript onclick handler.
// The full sized image is shown for this image.
func (g *Gallery) ShowPic(this js.Value, p []js.Value) any {
	g.ImageDisplay(p[0].Int())
	return js.ValueOf(false)
}

// ImageDisplay displays the full sized image, lazily building
// the HTML page for this image as required.
func (g *Gallery) ImageDisplay(index int) {
	if index < 0 {
		index = 0
	} else if index >= len(g.images) {
		index = len(g.images) - 1
	}
	if g.imagePage && index == g.curImage {
		return
	}
	g.imagePage = true
	g.curImage = index
	img := g.images[index]
	if len(img.imagePage) == 0 {
		g.BuildPict(index)
	}
	g.w.Display(img.imagePage)
}

// ShowThumbs is a callback from a JS onclick event,
// and will set the current image and show the thumbnail page.
func (g *Gallery) ShowThumbs(this js.Value, p []js.Value) any {
	g.curImage = p[0].Int()
	g.ShowPage()
	return js.ValueOf(false)
}

// ShowPage displays the thumbnail page of the current image.
func (g *Gallery) ShowPage() {
	g.imagePage = false
	var c strings.Builder
	g.w.SetTitle(g.title)
	g.cols, g.rows = g.tableSize()
	perPage := g.rows * g.cols
	nPages := (len(g.images) + perPage - 1) / perPage
	curPage := g.curImage / perPage
	if nPages > 1 {
		c.WriteString(h.Div(h.Open(), h.Id("navlinks"), "Pages: "))
		for i := 0; i < nPages; i++ {
			c.WriteString(g.LinkToPage(h.Text(i+1), i, i*perPage, h.Text(h.If(curPage == i), "current")))
		}
		c.WriteString(h.Div(h.Close()))
	}
	c.WriteString(g.header)
	c.WriteString(h.Div(h.Id("thumbpage"), h.Open()))
	i := curPage * g.rows * g.cols
	g.firstImage = i
	for x := 0; x < g.rows; x++ {
		for y := 0; y < g.cols; y++ {
			if i < len(g.images) {
				c.WriteString(g.images[i].thumbEntry)
				i++
			}
		}
		c.WriteString(h.Br(h.Style("clear: left")))
	}
	g.lastImage = i - 1
	c.WriteString(h.Div(h.Close()))
	c.WriteString(h.Br(h.Style("clear: both")))
	c.WriteString(Copyright(g.owner))
	g.w.Display(c.String())
	g.updateThumb(g.curImage, thumbOn)
}

// LinkToPage generates HTML for a link to a thumbnail page.
func (g *Gallery) LinkToPage(txt string, pageNo, index int, class string) string {
	return h.A(h.Class(class), h.Id(h.If(pageNo >= 0), h.Text("navlink", pageNo)),
			h.Onclick(h.Text("return showThumbs(", index, ")")),
			h.Href("#"),
			txt)
}

// BuildPict creates the full page HTML for this image
func (g *Gallery) BuildPict(index int) {
	img := g.images[index]
	var c strings.Builder
	g.LinkToPict(&c, "prev", index-1)
	c.WriteString(h.Div(h.Id("home"), h.A(h.Onclick("return showThumbs(", index, ")"), h.Href("#"), "back to index")))
	g.LinkToPict(&c, "next", index+1)
	var t string
	if len(img.title) > 0 {
		t = img.title
	} else {
		t = g.title
	}
	c.WriteString(g.HeaderDownload(t, "", img.download))
	c.WriteString(h.Div(h.Id("mainimage"), h.Img(h.Src(img.filename), h.Alt(t))))
	// Show image properties etc.
	c.WriteString(h.Div(h.Open(), h.Class("properties"), h.Table(h.Open(), h.Summary("image properties"), h.Border(0))))
	c.WriteString(g.Property("Date", img.date))
	c.WriteString(g.Property("Filename", img.name))
	if img.original.Width != 0 && img.original.Height != 0 {
		c.WriteString(g.Property("Original resolution", h.Text(img.original.Width, " x ", img.original.Height)))
	}
	c.WriteString(g.Property("Exposure", img.exposure))
	c.WriteString(g.Property("Aperture", img.aperture))
	c.WriteString(g.Property("ISO", img.iso))
	c.WriteString(g.Property("Focal length (mm)", img.flen))
	c.WriteString(h.Table(h.Close()))
	c.WriteString(h.Div(h.Close()))
	c.WriteString(Copyright(g.owner))
	// Add links for prefetching the next and previous images
	if index > 0 {
		c.WriteString(h.Link(h.Rel("prefetch"), h.Type("image/jpeg"), h.Href(g.images[index-1].filename)))
	}
	if index < len(g.images)-1 {
		c.WriteString(h.Link(h.Rel("prefetch"), h.Type("image/jpeg"), h.Href(g.images[index+1].filename)))
	}
	img.imagePage = c.String()
}

// Property generates HTML for the image metadata (name, size etc.)
func (g *Gallery) Property(n, val string) string {
	return h.Tr(h.If(val != ""), h.Td(h.Class("exifname"), n), h.Td(h.Class("exifdata"), val))
}

// LinkToPict generates HTML for a link to the selected picture.
func (g *Gallery) LinkToPict(c *strings.Builder, n string, index int) {
	c.WriteString(h.Div(h.Open(), h.Id(n)))
	if index < 0 || index == len(g.images) {
		c.WriteString("&nbsp")
	} else {
		c.WriteString(h.A(h.Open(), h.Onclick("return showPict(", index, ")"), h.Href("#")))
		if len(g.images[index].title) == 0 {
			c.WriteString(g.images[index].name)
		} else {
			c.WriteString(g.images[index].title)
		}
		c.WriteString(h.A(h.Close()))
	}
	c.WriteString(h.Div(h.Close()))
}

func (g *Gallery) HeaderDownload(title, back, download string) string {
	var c strings.Builder
	c.WriteString(h.H1(h.Open()))
	if download != "" {
		c.WriteString(hSpace)
	}
	if back != "" {
		c.WriteString(h.A(h.Open(), h.Href(back)))
	}
	c.WriteString(title)
	if back != "" {
		c.WriteString(h.A(h.Close()))
	}
	if download != "" {
		c.WriteString(h.Span(hSpace, h.A(h.Download(), h.Href(download), rune(0x21A7))))
	}
	c.WriteString(h.H1(h.Close()))
	return c.String()
}

// updateThumb sets the class for the selected image (used to
// highlight the current image).
func (g *Gallery) updateThumb(i int, cl string) {
	id := h.Text("slide", i)
	g.w.GetById(id).Set("className", js.ValueOf(cl))
}

// tableSize returns the column and row count for the thumbnail page.
func (g *Gallery) tableSize() (int, int) {
	return max(g.w.Width/(g.tw+14), 1), max((g.w.Height-170)/(g.th+33), 1)
}
