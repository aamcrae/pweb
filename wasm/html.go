package main

import (
	"fmt"
	"strconv"
	"strings"

	_ "syscall/js"
)

type Attr string
type flag int

const (
	f_drop flag = 1 << iota
	f_no_open
	f_no_close
)

func H1(elems ...any) string {
	return tag("h1", elems)
}

func H2(elems ...any) string {
	return tag("h2", elems)
}

func H3(elems ...any) string {
	return tag("h2", elems)
}

func Img(elems ...any) string {
	return tag("img", elems)
}

func Div(elems ...any) string {
	return tag("div", elems)
}

func A(elems ...any) string {
	return tag("a", elems)
}

func Span(elems ...any) string {
	return tag("span", elems)
}

func Table(elems ...any) string {
	return tag("table", elems)
}

func Tr(elems ...any) string {
	return tag("tr", elems)
}

func Td(elems ...any) string {
	return tag("td", elems)
}

func Br(elems ...any) string {
	return emptyTag("br", elems)
}

func P(elems ...any) string {
	return emptyTag("p", elems)
}

// Attributes

func Alt(s ...any) Attr {
	return attr("alt", s)
}

func Title(s ...any) Attr {
	return attr("title", s)
}

func Src(s ...any) Attr {
	return attr("src", s)
}

func Onclick(s ...any) Attr {
	return attr("onclick", s)
}

func Href(s ...any) Attr {
	return attr("href", s)
}

func Border(s ...any) Attr {
	return attr("border", s)
}

func Summary(s ...any) Attr {
	return attr("summary", s)
}

func Class(s ...any) Attr {
	return attr("class", s)
}

func Id(s ...any) Attr {
	return attr("id", s)
}

func Style(s ...any) Attr {
	return attr("style", s)
}

func Download(s ...any) Attr {
	return attr("download", s)
}

func If(c bool) flag {
	if !c {
		return f_drop
	} else {
		return 0
	}
}

func Open() flag {
	return f_no_close
}

func Close() flag {
	return f_no_open
}

func Text(s ...any) string {
	var b strings.Builder
	wrAll(&b, s, false)
	return b.String()
}
func tag(nm string, elems []any) string {
	return wrTag(nm, elems, false)
}

func emptyTag(nm string, elems []any) string {
	return wrTag(nm, elems, true)
}

func wrTag(nm string, elems []any, empty bool) string {
	atrs, other, flags := unpack(elems)
	if (flags & f_drop) != 0 {
		return ""
	}
	var sb strings.Builder
	if (flags & f_no_open) == 0 {
		sb.WriteRune('<')
		sb.WriteString(nm)
		wrAll(&sb, atrs, true)
		sb.WriteRune('>')
	}
	wrAll(&sb, other, false)
	if !empty && (flags & f_no_close)==0 {
		sb.WriteString("</")
		sb.WriteString(nm)
		sb.WriteString(">")
	}
	return sb.String()
}

func attr(nm string, elems []any) Attr {
	atrs, other, flags := unpack(elems)
	if (flags & f_drop) != 0 || len(atrs) > 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteRune(' ')
	sb.WriteString(nm)
	sb.WriteString("=\"")
	wrAll(&sb, other, false)
	sb.WriteString("\"")
	return Attr(sb.String())
}

func unpack(s []any) ([]any, []any, flag) {
	var other []any
	var atrs []any
	var flags flag
	for _, ele := range s {
		switch v := ele.(type) {
		case Attr:
			atrs = append(atrs, ele)
		case flag:
			flags |= v
		default:
			other = append(other, ele)
		}
	}
	return atrs, other, flags
}

func wrAll(sb *strings.Builder, s []any, space bool) {
	for _, ele := range s {
		if space {
			sb.WriteRune(' ')
		}
		wr(sb, ele)
	}
}

func wr(sb *strings.Builder, s any) {
	switch v := s.(type) {
	case string:
		sb.WriteString(v)
	case Attr:
		sb.WriteString(string(v))
	case fmt.Stringer:
		sb.WriteString(v.String())
	case []byte:
		sb.Write(v)
	case rune:
		sb.WriteRune(v)
	case int:
		sb.WriteString(strconv.FormatInt(int64(v), 10))
	default:
		panic("wr: Unknown type")
	}
}
