package main

import (
	"fmt"
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
	baseName    string // Base filename

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
	baseName := name
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
		baseName:    baseName,
	}, nil
}

// GetExif returns the EXIF data for the picture, loading it from the
// file if it is not already loaded.
func (p *Pict) GetExif() (*Exif, error) {
	if p.exif == nil {
		var err error
		if p.exif, err = ReadExif(p.srcFile); err != nil {
			return nil, fmt.Errorf("%s: exif read %v", p.srcFile, err)
		}
		if p.exif.ts.IsZero() {
			// Use file timestamp
			p.exif.ts = p.mtime
		}
	}
	return p.exif, nil
}

// MustExif returns the EXIF data for the picture, or exits on error.
func (p *Pict) MustExif() *Exif {
	exif, err := p.GetExif()
	if err != nil {
		panic(err.Error())
	}
	return exif
}

// AddGallery adds this picture to the gallery XML structure.
func (p *Pict) AddToGallery(g *data.Gallery, download int) error {
	var ph data.Photo
	exif, err := p.GetExif()
	if err != nil {
		return err
	}
	ph.Name = p.baseName
	ph.Filename = p.destFile
	ph.Date = exif.ts.Format("03:04 PM Monday, 02 January 2006")
	ph.Original.Width = p.width
	ph.Original.Height = p.height
	ph.Title = exif.title
	ph.Caption = exif.caption
	ph.Exposure = exif.exposure
	ph.ISO = exif.iso
	ph.Aperture = exif.fstop
	ph.FocalLength = exif.focal_len
	if download != DL_NONE {
		ph.Download = p.dlFile
	}
	g.Photos = append(g.Photos, ph)
	return nil
}

// Resize resizes this picture to a thumbnail size, a preview size, and
// a web page size. A resizer function is provided to perform the action
// to allow selection of different image processors.
func (p *Pict) Resize(handler NewImage, tw, th, pw, ph, iw, ih int) error {
	exif, err := p.GetExif()
	if err != nil {
		return err
	}
	if exif.width != 0 && exif.height != 0 {
		p.width = exif.width
		p.height = exif.height
	}
	// Check whether timestamp is the same, and we have the original resolution
	destPath := path.Join(p.destDir, p.destFile)
	mt, _ := getMtime(destPath)
	if mt == p.mtime && p.width > 0 && p.height > 0 {
		if *verbose {
			fmt.Printf("Skipping read/decode of %s\n", p.destFile)
		}
		return nil
	}
	img, err := handler(p.srcFile)
	if err != nil {
		return err
	}
	p.width = img.Width()
	p.height = img.Height()
	if mt == p.mtime {
		if *verbose {
			fmt.Printf("Skipping resize of %s\n", p.destFile)
		}
		return nil
	}
	if *verbose {
		fmt.Printf("Resizing %s from %d x %d\n", p.srcFile, img.Width(), img.Height())
	}
	switch exif.orientation {
	case "8":
		img.Rotate(Rotate90)
	case "3":
		img.Rotate(Rotate180)
	case "6":
		img.Rotate(Rotate270)
	}
	if err := img.Write(destPath, p.mtime, iw, ih, 90); err != nil {
		return err
	}
	if err := img.Write(path.Join(p.destDir, p.previewFile), p.mtime, pw, ph, 80); err != nil {
		return err
	}
	return img.Write(path.Join(p.destDir, p.thumbFile), p.mtime, tw, th, 80)
}
