package main

import (
	"encoding/xml"
	"fmt"

	"syscall/js"

	"github.com/aamcrae/pweb/data"
	"github.com/aamcrae/pweb/wasm"
)

const (
	thumbOn  = "slideshowon"
	thumbOff = "slideshow"
)

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
	w          *wasm.Window
	title      string   // Title of gallery
	header     string   // HTML of title of page
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
	w.OnKey(g.KeyPress)
	w.OnSwipe(g.Swipe)
	w.AddJSFunction("showPict", g.ShowPic)
	w.AddJSFunction("showThumbs", g.ShowThumbs)
	g.ShowPage()
	w.Wait()
}

// newGallery creates a new gallery from the XML data provided.
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
			c.Wr("<a href=\"").Wr(xmlData.Back).Wr("\"><h1>").Wr(g.title).Wr("</h1></a>")
		} else {
			c.Wr("<h1>").Wr(g.title).Wr("</h1>")
		}
		g.header = c.String()
	}
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
		var ct wasm.Comp
		ct.Wr(fmt.Sprintf("<div class=\"holder\"><div id=\"slide%d\" class=\"slideshow\">", i))
		ct.Wr(fmt.Sprintf("<a onclick=\"return showPict(%d)\" href=\"#\">", i))
		ct.Wr("<img src=\"t/").Wr(img.filename).Wr("\" title=\"").Wr(img.title).Wr("\">")
		ct.Wr("</a>")
		if len(img.title) > 0 {
			ct.Wr("<div class=thumbName>").Wr(img.title).Wr("</div>")
		}
		ct.Wr("</div> </div>")
		img.thumbEntry = ct.String()
		g.images = append(g.images, img)
	}
	// Install some style elements now that we know the thumbnail sizes
	c := new(wasm.Comp)
	c.Wr(".holder {width:").Wr(g.tw + 10).Wr("px;height:").Wr(g.th + 30).Wr("px}")
	c.Wr(".thumbName{width:").Wr(g.tw + 10).Wr("px}")
	g.w.AddStyle(c.String())
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
func (g *Gallery) Swipe(d wasm.Direction) {
	if g.imagePage {
		switch d {
		case wasm.Down:
			g.SelectThumb(g.curImage)
		case wasm.Right:
			g.ImageDisplay(g.curImage - 1)
		case wasm.Left:
			g.ImageDisplay(g.curImage + 1)
		}
	} else {
		perPage := g.rows * g.cols
		switch d {
		case wasm.Right:
			g.SelectThumb(g.curImage - perPage)
		case wasm.Left:
			g.SelectThumb(g.curImage + perPage)
		}
	}
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
		case "ArrowRight":
			g.ImageDisplay(g.curImage + 1)
		case "ArrowLeft":
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
	c := new(wasm.Comp)
	g.w.SetTitle(g.title)
	g.cols, g.rows = g.tableSize()
	perPage := g.rows * g.cols
	nPages := (len(g.images) + perPage - 1) / perPage
	curPage := g.curImage / perPage
	if nPages > 1 {
		c.Wr("<div id=\"navlinks\">Pages: ")
		for i := 0; i < nPages; i++ {
			class := ""
			if curPage == i {
				class = "current"
			}
			g.LinkToPage(c, fmt.Sprintf("%d", i+1), i, i*perPage, class)
		}
		c.Wr("</div>")
	}
	c.Wr(g.header)
	c.Wr("<div id=\"thumbpage\">")
	i := curPage * g.rows * g.cols
	g.firstImage = i
	for x := 0; x < g.rows; x++ {
		for y := 0; y < g.cols; y++ {
			if i < len(g.images) {
				c.Wr(g.images[i].thumbEntry)
				i++
			}
		}
		c.Wr("<br style=\"clear: left\" />\n")
	}
	g.lastImage = i - 1
	c.Wr("</div><br style=\"clear: both\" />\n")
	c.Copyright(g.owner)
	g.w.Display(c.String())
	g.updateThumb(g.curImage, thumbOn)
}

// LinkToPage generates HTML for a link to a thumbnail page.
func (g *Gallery) LinkToPage(c *wasm.Comp, txt string, pageNo, index int, class string) {
	c.Wr("<a class=\"").Wr(class).Wr("\" ")
	if pageNo >= 0 {
		c.Wr("id=\"navlink").Wr(pageNo).Wr("\" ")
	}
	c.Wr("onclick=\"return showThumbs(").Wr(index).Wr(")\" href=\"#\">").Wr(txt).Wr("</a>")
}

// BuildPict creates the full page HTML for this image
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
	c.Wr("<div id=\"mainimage\"><img src=\"").Wr(img.filename).Wr("\" alt=\"").Wr(img.title).Wr("\"></div>")
	// Show image properties etc.
	c.Wr("<div class=\"properties\"><table summary=\"image properties\" border=\"0\">")
	g.Property(c, "Date", img.date)
	if img.download == "" {
		g.Property(c, "Filename", img.name)
	} else {
		g.Property(c, "Filename", fmt.Sprintf("%s <a href=\"%s\" download>[Download original]</a>", img.name, img.download))
	}
	if img.original.Width != 0 && img.original.Height != 0 {
		g.Property(c, "Original resolution", fmt.Sprintf("%d x %d", img.original.Width, img.original.Height))
	}
	g.Property(c, "Exposure", img.exposure)
	g.Property(c, "Aperture", img.aperture)
	g.Property(c, "ISO", img.iso)
	g.Property(c, "Focal length", img.flen)
	c.Wr("</table></div>")
	c.Copyright(g.owner)
	img.imagePage = c.String()
}

// Property generates HTML for the image metadata (name, size etc.)
func (g *Gallery) Property(c *wasm.Comp, n, val string) {
	if val != "" {
		c.Wr("<tr><td class=\"exifname\">").Wr(n).Wr("</td><td class=\"exifdata\">").Wr(val).Wr("</td></tr>")
	}
}

// LinkToPict generates HTML for a link to the selected picture.
func (g *Gallery) LinkToPict(c *wasm.Comp, n string, index int) {
	c.Wr("<div id=\"").Wr(n).Wr("\">")
	if index < 0 || index == len(g.images) {
		c.Wr("&nbsp")
	} else {
		c.Wr("<a onclick=\"return showPict(").Wr(index).Wr(")\" href=\"#\">")
		if len(g.images[index].title) == 0 {
			c.Wr(g.images[index].name)
		} else {
			c.Wr(g.images[index].title)
		}
		c.Wr("</a>")
	}
	c.Wr("</div>")
}

// updateThumb sets the class for the selected image (used to
// highlight the current image).
func (g *Gallery) updateThumb(i int, cl string) {
	id := fmt.Sprintf("slide%d", i)
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
