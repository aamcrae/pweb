package wasm

import (
	"fmt"
	"strconv"
	"strings"

	_ "syscall/js"
)

type Comp struct {
	strings.Builder
}

func (c *Comp) Wr(s any) *Comp {
	switch v := s.(type) {
	case string:
		c.Write([]byte(v))
	case fmt.Stringer:
		c.Write([]byte(v.String()))
	case []byte:
		c.Write(v)
	case int:
		c.Write([]byte(strconv.FormatInt(int64(v), 10)))
	default:
		fmt.Println("Wr: Unknown type")
	}
	return c
}

func (c *Comp) Copyright(owner string) *Comp {
	if len(owner) > 0 {
		c.Wr("<div id=\"copyright\">&nbsp; &copy; Copyright ").Wr(owner).Wr("</div>")
	}
	return c
}
