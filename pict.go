package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/aamcrae/pweb/data"
)

type Pict struct {
	srcFile     string // Source filename, relative to cwd
	srcPath     string // Full pathname of source file
	destDir     string // Destination directory for web page
	thumbFile   string // Thumbnail filename relative to destDir
	dlFile      string // Download filename relative to destDir
	previewFile string // Preview filename relative to destDir
	destFile    string // Image filename relative to destDir

	mtime         time.Time // File modified time
	exif          *Exif     // Lazily loaded Exif data
	width, height int
}

func NewPict(fname, srcDir, destDir string) (*Pict, error) {
	mtime, err := getMtime(fname)
	if err != nil {
		return nil, err
	}
	d, name := path.Split(fname)
	for d != "" {
		var f string
		d = path.Dir(d)
		d, f = path.Split(d)
		name = f + "_" + name
	}
	return &Pict{
		srcFile:     fname,
		srcPath:     path.Join(srcDir, fname),
		destDir:     destDir,
		thumbFile:   path.Join("t", name),
		dlFile:      path.Join("d", name),
		previewFile: path.Join("p", name),
		destFile:    name,
		mtime:       mtime,
	}, nil
}

// GetExif returns the EXIF data for the picture, loading it from the
// file if it is not already loaded.
func (p *Pict) GetExif() *Exif {
	if p.exif == nil {
		var err error
		if p.exif, err = ReadExif(p.srcFile); err != nil {
			log.Fatalf("%s: exif read %v", p.srcFile, err)
		}
		if p.exif.ts.IsZero() {
			// Use file timestamp
			p.exif.ts = p.mtime
		}
	}
	return p.exif
}

// AddGallery adds this picture to the gallery XML structure.
func (p *Pict) AddToGallery(g *data.Gallery, download bool) {
	var ph data.Photo
	exif := p.GetExif()
	ph.Name = p.destFile
	ph.Date = exif.ts.Format("03:04 PM Monday, 02 January 2006")
	ph.Original.Width = p.width
	ph.Original.Height = p.height
	ph.Title = exif.title
	ph.Caption = exif.caption
	ph.ISO = exif.iso
	ph.Aperture = exif.fstop
	ph.FocalLength = exif.focal_len
	if download {
		ph.Download = p.dlFile
	}
	g.Photos = append(g.Photos, ph)
}

// Resize resizes this picture to a thumbnail size, a preview size, and
// a web page size. A resizer function is provided to perform the action
// to allow selection of different image processors.
func (p *Pict) Resize(handler NewImage, tw, th, pw, ph, iw, ih int) {
	img, err := handler(p.srcFile)
	if err != nil {
		log.Fatalf("%s: %v", p.srcFile, err)
	}
	if *verbose {
		fmt.Printf("Resizing %s from %d x %d\n", p.srcFile, img.Width(), img.Height())
	}
	p.width = img.Width()
	p.height = img.Height()
	// Check whether timestamp is the same
	destPath := path.Join(p.destDir, p.destFile)
	mt, _ := getMtime(destPath)
	if mt == p.mtime {
		if *verbose {
			fmt.Printf("Skipping resize of %s\n", p.destFile)
		}
		return
	}
	switch p.GetExif().orientation {
	case "8":
		img.Rotate(Rotate90)
	case "3":
		img.Rotate(Rotate180)
	case "6":
		img.Rotate(Rotate270)
	}
	img.Write(destPath, p.mtime, iw, ih, 90)
	img.Write(path.Join(p.destDir, p.previewFile), p.mtime, pw, ph, 80)
	img.Write(path.Join(p.destDir, p.thumbFile), p.mtime, tw, th, 80)
}
