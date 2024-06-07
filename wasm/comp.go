package main

import (
	"fmt"
	"strconv"
	"strings"

	_ "syscall/js"
)

// Comp is a simple composer that understands the JS values
type Comp struct {
	strings.Builder
}

// Wr writes the value to the string builder
func (c *Comp) Wr(s any) *Comp {
	switch v := s.(type) {
	case string:
		c.WriteString(v)
	case fmt.Stringer:
		c.WriteString(v.String())
	case []byte:
		c.Write(v)
	case rune:
		c.WriteRune(v)
	case int:
		c.WriteString(strconv.FormatInt(int64(v), 10))
	default:
		fmt.Println("Wr: Unknown type")
	}
	return c
}
