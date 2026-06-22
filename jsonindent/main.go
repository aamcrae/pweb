package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aamcrae/pweb/shared"
)

var input = flag.String("input", "", "Input JSON file")
var write = flag.Bool("write", false, "Rewrite JSON file")

func main() {
	flag.Parse()
	switch filepath.Base(*input) {
	case shared.AlbumFileJSON:
		var adata shared.AlbumPage
		read(*input, *write, &adata)
	case shared.GalleryFileJSON:
		var gdata shared.Gallery
		read(*input, *write, &gdata)
	default:
		log.Fatalf("Unknown JSON file: %s", *input)
	}
}

func read(in string, write bool, d any) {
	var err error
	f, err := os.ReadFile(in)
	if err != nil {
		log.Fatalf("%s: %v", in, err)
	}
	if err = json.Unmarshal(f, d); err != nil {
		log.Fatalf("JSON unmarshal %s: %v", in, err)
	}
	var of *os.File
	if write {
		of, err = os.Create(in)
		if err != nil {
			log.Fatalf("File create %s: %v", in, err)
		}
		defer of.Close()
	} else {
		of = os.Stdout
	}
	b, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		log.Fatalf("JSON marshal %s: %v", in, err)
	}
	of.Write(b)
	of.Write([]byte{'\n'})
}
