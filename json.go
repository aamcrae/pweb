package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// readMeta reads the JSON from the file and populates the structure passed.
func readMeta(file string, d any) error {
	if af, err := os.ReadFile(file); err != nil {
		return err
	} else {
		return json.Unmarshal(af, d)
	}
}

// writeMeta writes the marshaled JSON to the file.
func writeMeta(file string, d any) error {
	if s, err := json.MarshalIndent(&d, "", " "); err != nil {
		return fmt.Errorf("%s: marshal %w", file, err)
	} else {
		return os.WriteFile(file, s, 0664)
	}
}
