package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"syscall/js"

	"github.com/aamcrae/pweb/data"
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
	w          *Window
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

func RunGallery(w *Window, gx []byte) {
	w.LoadStyle("/pweb/gallery-style.css")
	var gallery data.Gallery
	err := xml.Unmarshal(gx, &gallery)
	if err != nil {
		fmt.Printf("unmarshal: %v\n", err)
		w.Display(H1("Bad or no gallery data!"))
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
func newGallery(xmlData *data.Gallery, w *Window) *Gallery {
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
			Div(Class("holder"),
				Div(Class("slideshow"), Id(Text("slide", i)),
				A(Onclick(Text("return showPict(", i, ")")),
					Href("#"),
					Img(Title(img.title), Src(Text("t/", img.filename)))),
				Div(If(len(img.title) > 0), Class("thumbName"), img.title)))
		g.images = append(g.images, img)
	}
	// Install some style elements now that we know the thumbnail sizes
	g.w.AddStyle(Text(".holder {width:", g.tw + 10, "px;height:", g.th + 30, "px} .thumbName{width:", g.tw + 10, "px}"))
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
func (g *Gallery) Swipe(d Direction) bool {
	if g.imagePage {
		switch d {
		case Down:
			g.SelectThumb(g.curImage)
		case Right:
			g.ImageDisplay(g.curImage - 1)
		case Left:
			g.ImageDisplay(g.curImage + 1)
		default:
			return false
		}
	} else {
		perPage := g.rows * g.cols
		switch d {
		case Down:
			if len(g.back) > 0 {
				g.w.Goto(g.back)
			}
		case Right:
			g.SelectThumb(g.curImage - perPage)
		case Left:
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
		c.WriteString(Div(Open(), Id("navlinks"), "Pages: "))
		for i := 0; i < nPages; i++ {
			c.WriteString(g.LinkToPage(Text(i+1), i, i*perPage, Text(If(curPage == i), "current")))
		}
		c.WriteString(Div(Close()))
	}
	c.WriteString(g.header)
	c.WriteString(Div(Id("thumbpage"), Open()))
	i := curPage * g.rows * g.cols
	g.firstImage = i
	for x := 0; x < g.rows; x++ {
		for y := 0; y < g.cols; y++ {
			if i < len(g.images) {
				c.WriteString(g.images[i].thumbEntry)
				i++
			}
		}
		c.WriteString(Br(Style("clear: left")))
	}
	g.lastImage = i - 1
	c.WriteString(Div(Close()))
	c.WriteString(Br(Style("clear: both")))
	c.WriteString(Copyright(g.owner))
	g.w.Display(c.String())
	g.updateThumb(g.curImage, thumbOn)
}

// LinkToPage generates HTML for a link to a thumbnail page.
func (g *Gallery) LinkToPage(txt string, pageNo, index int, class string) string {
	return A(Class(class), Id(If(pageNo >= 0), Text("navlink", pageNo)),
			Onclick(Text("return showThumbs(", index, ")")),
			Href("#"),
			txt)
}

// BuildPict creates the full page HTML for this image
func (g *Gallery) BuildPict(index int) {
	img := g.images[index]
	var c strings.Builder
	g.LinkToPict(&c, "prev", index-1)
	c.WriteString(Div(Id("home"), A(Onclick("return showThumbs(", index, ")"), Href("#"), "back to index")))
	g.LinkToPict(&c, "next", index+1)
	var t string
	if len(img.title) > 0 {
		t = img.title
	} else {
		t = g.title
	}
	c.WriteString(g.HeaderDownload(t, "", img.download))
	c.WriteString(Div(Id("mainimage"), Img(Src(img.filename), Alt(t))))
	// Show image properties etc.
	c.WriteString(Div(Open(), Class("properties"), Table(Open(), Summary("image properties"), Border(0))))
	c.WriteString(g.Property("Date", img.date))
	c.WriteString(g.Property("Filename", img.name))
	if img.original.Width != 0 && img.original.Height != 0 {
		c.WriteString(g.Property("Original resolution", Text(img.original.Width, " x ", img.original.Height)))
	}
	c.WriteString(g.Property("Exposure", img.exposure))
	c.WriteString(g.Property("Aperture", img.aperture))
	c.WriteString(g.Property("ISO", img.iso))
	c.WriteString(g.Property("Focal length", img.flen))
	c.WriteString(Table(Close()))
	c.WriteString(Div(Close()))
	c.WriteString(Copyright(g.owner))
	img.imagePage = c.String()
}

// Property generates HTML for the image metadata (name, size etc.)
func (g *Gallery) Property(n, val string) string {
	return Tr(If(val != ""), Td(Class("exifname"), n), Td(Class("exifdata"), val))
}

// LinkToPict generates HTML for a link to the selected picture.
func (g *Gallery) LinkToPict(c *strings.Builder, n string, index int) {
	c.WriteString(Div(Open(), Id(n)))
	if index < 0 || index == len(g.images) {
		c.WriteString("&nbsp")
	} else {
		c.WriteString(A(Open(), Onclick("return showPict(", index, ")"), Href("#")))
		if len(g.images[index].title) == 0 {
			c.WriteString(g.images[index].name)
		} else {
			c.WriteString(g.images[index].title)
		}
		c.WriteString(A(Close()))
	}
	c.WriteString(Div(Close()))
}

func (g *Gallery) HeaderDownload(title, back, download string) string {
	var c strings.Builder
	c.WriteString(H1(Open()))
	if download != "" {
		c.WriteString(hSpace)
	}
	if back != "" {
		c.WriteString(A(Open(), Href(back)))
	}
	c.WriteString(title)
	if back != "" {
		c.WriteString(A(Close()))
	}
	if download != "" {
		c.WriteString(Span(hSpace, A(Download(), Href(download), rune(0x21A7))))
	}
	c.WriteString(H1(Close()))
	return c.String()
}

// updateThumb sets the class for the selected image (used to
// highlight the current image).
func (g *Gallery) updateThumb(i int, cl string) {
	id := Text("slide", i)
	current := g.w.GetById(id)
	if current.IsUndefined() || current.IsNull() {
		fmt.Printf("Can't find %s\n", id)
	} else {
		current.Set("className", js.ValueOf(cl))
	}
}

// tableSize returns the column and row count for the thumbnail page.
func (g *Gallery) tableSize() (int, int) {
	return max(g.w.Width/(g.tw+14), 1), max((g.w.Height-170)/(g.th+33), 1)
}
