package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aamcrae/pweb/data"
)

// UpdateAlbum will read the album XML file that references this
// gallery, and will add or update it if there is no matching entry.
// If reverse is set, then the entry will be added at the end
func UpdateAlbum(back, dest, dir, title string, reverse bool) {
	// Map to XML file from back href.
	albumDir := path.Dir(path.Join(dest, dir, back))
	album := path.Join(albumDir, data.AlbumFile)
	// Whatever happens with the album file, make sure that the album HTML is up to date.
	cpFile(path.Join(*assets, "index.html"), path.Join(albumDir, "index.html"))
	// The back reference (usually "../index.html") may actually refer
	// back to deeper levels, so to create the link forward, copy as
	// many directory elements as necessary from the destination directory
	d := dir
	link := "index.html"
	for i := strings.Count(back, "../"); i > 0; i-- {
		var ele string
		d, ele = path.Split(path.Clean(d))
		link = path.Join(ele, link)
	}
	var adata data.AlbumPage
	var exists bool
	err := ReadXml(album, &adata)
	if err == nil {
		// Search for gallery in album
		for ind, al := range adata.Albums {
			if al.Id == dir {
				// Gallery reference already exists, check that the data is the same,
				// otherwise rewrite it.
				if link == al.Link && al.Title == title {
					if *verbose {
						fmt.Printf("Gallery %s already present in %s\n", dir, album)
					}
					return
				}
				adata.Albums[ind].Link = link
				adata.Albums[ind].Title = title
				exists = true
				break
			}
		}
		if *verbose && !exists {
			fmt.Printf("Gallery %s being added to %s\n", dir, album)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path.Dir(album), 0755); err != nil {
			log.Fatalf("%s: %v", album, err)
		}
		log.Printf("%s: New album, please set title etc.", album)
		// Preload album data from template
		ReadXml(path.Join(*assets, data.TemplateAlbumFile), &adata)
	} else {
		log.Fatalf("%s: %v", album, err)
	}
	if !exists {
		newAlbum := data.Album{Link: link, Id: dir, Title: title}
		// Add new gallery reference either at the front (reverse false)
		// or appended.
		if reverse {
			adata.Albums = append(adata.Albums, newAlbum)
		} else {
			adata.Albums = append([]data.Album{newAlbum}, adata.Albums...)
		}
	}
	if newData, err := xml.MarshalIndent(&adata, "", " "); err != nil {
		log.Fatalf("%s: Marshal %v", album, err)
	} else {
		if err := os.WriteFile(album, newData, 0664); err != nil {
			log.Fatalf("%s: Write %v", album, err)
		}
	}
}

// ReadXml reads the XML from the file and populates the structure passed.
func ReadXml(file string, d any) error {
	if af, err := os.ReadFile(file); err != nil {
		return err
	} else {
		return xml.Unmarshal(af, d)
	}
}
