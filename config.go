package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	C_UP = iota
	C_TITLE
	C_DIR
	C_INCLUDE
	C_EXCLUDE
	C_STYLE
	C_AFTER
	C_BEFORE
	C_RATING
	C_SELECT
	C_DOWNLOAD
	C_NOCAPTION
	C_SORT
	C_REVERSE
	C_LARGE
	C_CAPTION
	C_NOZIP
	C_THUMB
)

// configOptions contains some options for the configuration keywords.
type configOptions struct {
	code  int
	min   int
	multi bool
	count int
}

var configKeywords = map[string]*configOptions{
	"up":        &configOptions{code: C_UP, min: 1},
	"title":     &configOptions{code: C_TITLE, min: 1},
	"dir":       &configOptions{code: C_DIR, min: 1},
	"include":   &configOptions{code: C_INCLUDE, min: 1, multi: true},
	"exclude":   &configOptions{code: C_EXCLUDE, min: 1, multi: true},
	"style":     &configOptions{code: C_STYLE, min: 1},
	"after":     &configOptions{code: C_AFTER, min: 2, multi: true},
	"before":    &configOptions{code: C_BEFORE, min: 2, multi: true},
	"rating":    &configOptions{code: C_RATING, min: 1},
	"select":    &configOptions{code: C_SELECT, min: 1},
	"download":  &configOptions{code: C_DOWNLOAD},
	"nocaption": &configOptions{code: C_NOCAPTION},
	"sort":      &configOptions{code: C_SORT, min: 1},
	"reverse":   &configOptions{code: C_REVERSE},
	"large":     &configOptions{code: C_LARGE},
	"caption":   &configOptions{code: C_CAPTION, min: 2, multi: true},
	"nozip":     &configOptions{code: C_NOZIP},
	"thumb":     &configOptions{code: C_THUMB},
}

type Config map[int][]string

// ReadConfig parses the config file and stores the parameters into
// a map. The map value is the parameters for the keyword.
// Some keywords may have multiple entries - these are added to the
// string slice for the keyword.
func ReadConfig(f string) Config {
	conf := make(Config)
	b, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("%s: %v", f, err)
	}
	for i, l := range strings.Split(string(b), "\n") {
		if len(l) == 0 || l[0] == '#' {
			continue
		}
		cmd := strings.SplitN(l, ":", 2)
		if len(cmd) != 2 {
			log.Fatalf("%s: line %d, Illegal config", f, i+1)
		}
		if c, ok := configKeywords[cmd[0]]; !ok {
			log.Fatalf("%s: line %d, unknown keyword (%s)", f, i+1, cmd[0])
		} else {
			if !c.multi && c.count > 0 {
				log.Fatalf("%s: line %d, duplicate keyword (%s)", f, i+1, cmd[0])
			}
			if len(strings.Fields(cmd[1])) < c.min {
				log.Fatalf("%s: line %d, not enough arguments for '%s'", f, i, cmd[0])
			}
			trimmed := strings.TrimLeft(cmd[1], " ")
			conf[c.code] = append(conf[c.code], trimmed)
			c.count++
			if *verbose {
				fmt.Printf("%s: line %d, keyword %s, args=<%s>\n", f, i, cmd[0], trimmed)
			}
		}
	}
	return conf
}
