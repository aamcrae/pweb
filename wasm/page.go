package wasm

import (
	"fmt"
	"strings"

	_ "syscall/js"
)

type Page struct {
	strings.Builder
}

func NewPage() *Page {
	p := &Page{}
	return p
}

func (p *Page) Wr(s any) *Page {
	switch v := s.(type) {
	case string:
		p.Write([]byte(v))
	case fmt.Stringer:
		p.Write([]byte(v.String()))
	case []byte:
		p.Write(v)
	default:
		fmt.Println("Wr: Unknown type")
	}
	return p
}

func (p *Page) Copyright(c string) *Page {
	if len(c) > 0 {
		p.Wr("<div id=\"copyright\">&nbsp; &copy; Copyright ").Wr(c).Wr("</div>")
	}
	return p
}
