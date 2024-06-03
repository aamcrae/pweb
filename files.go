package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	be "github.com/thomasheller/braceexpansion"
)

// globFiles expands the wildcard file list, and returns
// the list of files matching the wildcards.
func globFiles(in []string) []string {
	var files []string
	for _, f := range in {
		for _, splitF := range strings.Fields(f) {
			tree, err := be.New().Parse(splitF)
			if err != nil {
				log.Fatalf("Brace expand error %s: %v", f, err)
			}
			for _, exp := range tree.Expand() {
				fl, err := filepath.Glob(exp)
				if err != nil {
					log.Fatalf("Error in expanding %s: %v", f, err)
				}
				files = append(files, fl...)
			}
		}
	}
	return files
}

// Add the list of file names to an existing list, using the first
// filename in each entry as an anchor. The list may be added
// before the anchor, or after, depending on the argument.
func insert(flist []string, list []string, before bool) []string {
	m := make(map[string][]string)
	for _, il := range list {
		iEntry := strings.Fields(il)
		m[iEntry[0]] = append(m[iEntry[0]], iEntry[1:]...)
	}
	var newFiles []string
	for _, f := range flist {
		if v, ok := m[f]; ok {
			if before {
				newFiles = append(newFiles, globFiles(v)...)
			}
			newFiles = append(newFiles, f)
			if !before {
				newFiles = append(newFiles, globFiles(v)...)
			}
			delete(m, f)
		} else {
			newFiles = append(newFiles, f)
		}
	}
	// Check if any files left over.
	if len(m) != 0 {
		for k, _ := range m {
			log.Printf("Could not locate %s", k)
		}
		log.Fatalf("Exiting")
	}
	return newFiles
}

// find does a linear search of the list for the requested filename.
func find(list []string, name string) (int, bool) {
	for i, v := range list {
		if v == name {
			return i, true
		}
	}
	return -1, false
}

// makeDirs ensures that the directories passed exist, creating them if necessary.
func makeDirs(dir ...string) {
	for _, d := range dir {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("%s: Mkdir %v", d, err)
		}
	}
}

// cpMaybe will copy the file if the source exists
func cpMaybe(src, dst string) error {
	_, err := os.Stat(src)
	if err == nil {
		return cpFile(src, dst)
	}
	// No error if source does not exist
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// cpFile will copy the src file to the dst filename if the
// modify time is different.
func cpFile(src, dst string) error {
	if st, err := getMtime(src); err != nil {
		return err
	} else {
		if dt, err := getMtime(dst); err == nil && dt == st {
			// mtimes are the same
			return nil
		} else {
			// File is either non-existent or out of date
			if b, err := os.ReadFile(src); err != nil {
				return err
			} else {
				return cp(b, dst, st)
			}
		}
	}
}

// cpBytes will copy the src data to the dst filename if the
// modify time is different.
func cpBytes(src []byte, mtime time.Time, dst string) error {
	if t, err := getMtime(dst); err == nil {
		if t == mtime {
			return nil
		}
	}
	// dst file is either non-existent or a different time.
	return cp(src, dst, mtime)
}

// cp copies the byte slice src to dest, and adjusts the mtime to match.
func cp(src []byte, dst string, mtime time.Time) error {
	if err := os.WriteFile(dst, src, 0644); err != nil {
		return err
	}
	return os.Chtimes(dst, mtime, mtime)
}

// getMtime gets the modified time of the file.
// If an error is returned, it must be a not-existing file.
func getMtime(f string) (time.Time, error) {
	if fi, err := os.Stat(f); err != nil {
		var zeroT time.Time
		if errors.Is(err, os.ErrNotExist) {
			return zeroT, err
		}
		log.Fatalf("%s: %v", f, err)
		return zeroT, err
	} else {
		return fi.ModTime(), nil
	}
}
