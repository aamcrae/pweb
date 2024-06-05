package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/aamcrae/pweb/data"
)

const (
	thumbWidth    = 160
	thumbHeight   = 160
	previewWidth  = 320
	previewHeight = 240
)

var imageWidth int = 1500
var imageHeight int = 1200

const configDefault = ".web"

var verbose = flag.Bool("verbose", false, "Verbose output")
var force = flag.Bool("force", false, "Force rebuild")
var baseDir = flag.String("base", "/var/www/html/nphotos", "Base directory of web pages")
var assets = flag.String("assets", "/usr/share/pweb", "Source directory of web assets")
var imager = flag.String("imager", "dis", "Select the image handler")
var watchdog = flag.Int("watchdog", 20, "Timeout in seconds of watchdog")

// rScaleMap maps a selected rating to photo ratings that will be accepted
// e.g a rating of '3' will select photos with a rating of '3', '4' and '5'.
var rScaleMap = map[string][]string{
	"0": {"0", "1", "2", "3", "4", "5"},
	"1": {"1", "2", "3", "4", "5"},
	"2": {"2", "3", "4", "5"},
	"3": {"3", "4", "5"},
	"4": {"4", "5"},
	"5": {"5"},
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	var conf Config
	if len(args) == 0 {
		conf = ReadConfig(configDefault)
	} else if len(args) == 1 {
		conf = ReadConfig(args[0])
	} else {
		flag.Usage()
		log.Fatalf("Exiting...")
	}
	d, ok := conf[C_DIR]
	if !ok {
		log.Fatalf("%s: missing 'dir' config", args[0])
	}
	dir := d[0]
	destDir := path.Join(*baseDir, dir)
	if *verbose {
		fmt.Printf("Directory set to %s\n", destDir)
	}
	var files []string
	if incList, ok := conf[C_INCLUDE]; !ok {
		files = append(files, globFiles([]string{"*.jpg"})...)
	} else {
		files = append(files, globFiles(incList)...)
	}
	if *verbose {
		fmt.Printf("Include list: %v\n", files)
	}
	if excArg, ok := conf[C_EXCLUDE]; ok {
		for _, ex := range globFiles(excArg) {
			if ind, ok := find(files, ex); ok {
				files = append(files[:ind], files[ind+1:]...)
			} else {
				log.Printf("Cannot find %s in file list, ignored", ex)
			}
		}
	}
	if afterList, ok := conf[C_AFTER]; ok {
		files = insert(files, afterList, false)
	}
	if beforeList, ok := conf[C_BEFORE]; ok {
		files = insert(files, beforeList, true)
	}
	if *verbose {
		fmt.Printf("Before ratings and sorting: %v\n", files)
	}
	srcDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getwd %v", err)
	}
	// If a rating config is set, build a map of
	// allowed ratings (either as a scale or as selected
	// ratings)
	ratingMap := make(map[string]struct{})
	ratings, useRating := conf[C_RATING]
	sel, useSelect := conf[C_SELECT]
	if useRating && useSelect {
		log.Fatalf("Cannot use both select and rating")
	}
	if useRating {
		buildRatings(ratings, ratingMap, true)
	}
	if useSelect {
		buildRatings(sel, ratingMap, false)
	}
	// Build map of captions
	capt := make(map[string]string)
	cl, ok := conf[C_CAPTION]
	if ok {
		buildCaptions(cl, capt)
	}
	// If configured, sort by date.
	var dateSort bool
	if skey, ok := conf[C_SORT]; ok {
		if skey[0] == "date" {
			dateSort = true
		} else {
			log.Fatalf("Unknown sort order: %s", skey[0])
		}
	}
	picts := readPicts(files, srcDir, destDir, useSelect || useRating || dateSort || len(capt) > 0)
	if useSelect || useRating {
		picts = filterPicts(picts, ratingMap)
	}
	if len(capt) > 0 {
		addCaptions(picts, capt)
	}
	if dateSort {
		slices.SortStableFunc(picts, func(a, b *Pict) int {
			return a.GetExif().ts.Compare(b.GetExif().ts)
		})
	}
	if *verbose {
		fmt.Printf("Final list:")
		for _, p := range picts {
			fmt.Printf(" %s", p.srcFile)
		}
		fmt.Printf("\n")
	}
	// If force is on, delete the entire destination directory
	if *force {
		if err := os.RemoveAll(destDir); err != nil {
			log.Fatalf("%s: %v", destDir, err)
		}
	} else {
		// Remove any images no longer wanted
		removeUnwanted(destDir, picts)
	}
	if _, ok := conf[C_LARGE]; ok {
		imageWidth = 1800
		imageHeight = 1500
	}
	title, _ := conf[C_TITLE]
	up, upConfigured := conf[C_UP]
	_, reverse := conf[C_REVERSE]
	if upConfigured {
		UpdateAlbum(up[0], *baseDir, dir, title[0], reverse)
	}
	_, download := conf[C_DOWNLOAD]
	_, nozip := conf[C_NOZIP]
	// Ensure base page, thumbnail, preview and (optionally) download directories exist.
	makeDirs(destDir, path.Join(destDir, "t"), path.Join(destDir, "p"))
	if download {
		dlDir := path.Join(destDir, "d")
		makeDirs(dlDir)
		// If there is a .htaccess file required, copy it.
		if err := cpMaybe(path.Join(*assets, "download-htaccess"), path.Join(dlDir, ".htaccess")); err != nil {
			log.Fatalf("Write htaccess %v", err)
		}
	}
	var g data.Gallery
	// Preload gallery XML from template (to set copyright etc.)
	ReadXml(path.Join(*assets, data.TemplateGalleryFile), &g)
	g.Title = title[0]
	if download && !nozip {
		g.Download = path.Join("d", "photos.zip")
	}
	if upConfigured {
		g.Back = up[0]
	}
	g.Thumb.Width = thumbWidth
	g.Thumb.Height = thumbHeight
	g.Preview.Width = previewWidth
	g.Preview.Height = previewHeight
	g.Image.Width = imageWidth
	g.Image.Height = imageHeight
	var imgHandler NewImage
	switch *imager {
	case "vips":
		imgHandler = NewVipsImage
		vipsInit()
	case "dis":
		imgHandler = NewDisImage
	default:
		log.Fatalf("%s: Unknown imager", *imager)
	}
	// Now generate the scaled images that will appear on the web site.
	resizePhotos(imgHandler, picts, download)
	// Add the images to the gallery - this is done after the
	// resize in order to capture the original resolution dimensions, which is
	// only know after the image is processed.
	for _, p := range picts {
		p.AddToGallery(&g, download)
	}
	// Write the gallery XML file
	gFile := path.Join(destDir, data.GalleryFile)
	if gData, err := xml.MarshalIndent(&g, "", " "); err != nil {
		log.Fatalf("%s: Marshal %v", gFile, err)
	} else {
		if err := os.WriteFile(gFile, gData, 0664); err != nil {
			log.Fatalf("%s: Write %v", gFile, err)
		}
	}
	if download && !nozip {
		updateZip(path.Join(destDir, "d"))
	}
	// Conditionally copy the main index.html file.
	if err := cpFile(path.Join(*assets, "gallery-index.html"), path.Join(destDir, "index.html")); err != nil {
		log.Fatalf("index.html: Update %v", err)
	}
}

// readPicts will create a photo object and optionally read the EXIF (if the EXIF
// data is required for further processing)
func readPicts(files []string, srcDir, destDir string, readExif bool) []*Pict {
	// Create a worker pool to read the EXIF data
	var unratedPicts []*Pict
	pWork := NewWorker(time.Second * time.Duration(*watchdog))
	for _, f := range files {
		p, err := NewPict(f, srcDir, destDir)
		if err != nil {
			log.Fatalf("%s: %v", f, err)
		}
		unratedPicts = append(unratedPicts, p)
		// Read the EXIF if required
		if readExif {
			pWork.Run(func() {
				_ = p.GetExif()
			})
		}
	}
	pWork.Wait()
	return unratedPicts
}

func filterPicts(inPicts []*Pict, ratingMap map[string]struct{}) []*Pict {
	var outPicts []*Pict
	for _, p := range inPicts {
		rating := p.GetExif().rating
		_, ok := ratingMap[rating]
		if !ok {
			if *verbose {
				fmt.Printf("%s: Skipping due to rating (%s)\n", p.srcFile, rating)
			}
			continue
		}
		outPicts = append(outPicts, p)
	}
	return outPicts
}

func addCaptions(pl []*Pict, capt map[string]string) {
	for _, p := range pl {
		if c, ok := capt[p.srcFile]; ok {
			if *verbose {
				fmt.Printf("%s: Setting title to <%s>\n", p.srcFile, c)
			}
			p.GetExif().title = c
		}
	}
}

// buildCaptions will build a map of image filenames to
// any captions that are defined in the config file.
func buildCaptions(cl []string, capt map[string]string) {
	for _, c := range cl {
		// Caption is of the form <img_file Caption to be added>
		cimg, caption, found := strings.Cut(c, " ")
		if found {
			capt[cimg] = caption
		}
	}
}

// buildRatings builds a map of ratings. The rating is optionally
// a minimum rating value (if scale is false), so that photos with this rating
// or higher are selected, or as a list of ratings, and only those photos
// that have the rating will be selected.
func buildRatings(ratings []string, ratingMap map[string]struct{}, scale bool) {
	rList := strings.Fields(ratings[0])
	if !scale {
		// Allow all of the listed ratings to be included.
		for _, r := range rList {
			ratingMap[r] = struct{}{}
		}
	} else {
		if len(rList) > 1 {
			log.Fatalf("Can only have one rating value")
		}
		// If rating scale is set, then only allow a single rating value,
		// which is treated as a minimum rating - any ratings this rating and
		// higher are included
		for _, v := range rScaleMap[rList[0]] {
			ratingMap[v] = struct{}{}
		}
	}
}

func resizePhotos(handler NewImage, picts []*Pict, download bool) {
	resizers := NewWorker(time.Second * time.Duration(*watchdog))
	defer resizers.Wait()
	for _, p := range picts {
		resizers.Run(func() {
			p.Resize(handler, thumbWidth, thumbHeight, previewWidth, previewHeight, imageWidth, imageHeight)
			if download {
				// create symlink in the download directory to the original file, if not already existing
				dlPath := path.Join(p.destDir, p.dlFile)
				if _, err := os.Stat(dlPath); err != nil {
					if err := os.Symlink(p.srcPath, dlPath); err != nil {
						log.Fatalf("%s: symlink %v", p.srcFile, err)
					}
				}
			}
		})
	}
}

func updateZip(destDir string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("(cd %s; zip -FSq photos.zip *)", destDir))
	if err := cmd.Run(); err != nil {
		log.Fatalf("update zip: %v", err)
	}
}

func removeUnwanted(destDir string, plist []*Pict) {
	files := make(map[string]struct{})
	// Get the list of all files in the thumbnail directory, and
	// add them to the map.
	dentries, err := os.ReadDir(path.Join(destDir, "t"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Fatalf("%s/t: %v", destDir, err)
	}
	for _, d := range dentries {
		if !d.IsDir() {
			files[d.Name()] = struct{}{}
		}
	}
	// Remove all pictures from the map that are going to be added
	for _, p := range plist {
		delete(files, p.destFile)
	}
	// The entries remaining are unwanted, so remove them.
	for k, _ := range files {
		os.Remove(path.Join(destDir, k))
		os.Remove(path.Join(destDir, "t", k))
		os.Remove(path.Join(destDir, "p", k))
		os.Remove(path.Join(destDir, "d", k))
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] config-file\n", os.Args[0])
	flag.PrintDefaults()
}
