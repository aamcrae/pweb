package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aamcrae/pweb/shared"
)

var input = flag.String("input", "", "Input XML file")
var write = flag.Bool("write", false, "Write JSON file")

func main() {
	flag.Parse()
	switch filepath.Base(*input) {
	case shared.AlbumFileXML:
		var adata shared.AlbumPage
		read(*input, *write, &adata)
	case shared.GalleryFileXML:
		var gdata shared.Gallery
		read(*input, *write, &gdata)
	default:
		log.Fatalf("Unknown XML file: %s", *input)
	}
}

func read(in string, write bool, d any) {
	var err error
	f, err := os.ReadFile(in)
	if err != nil {
		log.Fatalf("%s: %v", in, err)
	}
	if err = xml.Unmarshal(f, d); err != nil {
		log.Fatalf("XML unmarshal %s: %v", in, err)
	}
	var of *os.File
	if write {
		fn := strings.TrimSuffix(in, filepath.Ext(in)) + ".json"
		of, err = os.Create(fn)
		if err != nil {
			log.Fatalf("File create %s: %v", fn, err)
		}
		defer of.Close()
	} else {
		of = os.Stdout
	}
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal %s: %v", in, err)
	}
	of.Write(b)
}
