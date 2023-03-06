package main

import (
	"fmt"
	"os"
	"strings"
)

func defaultLibraryPath() (string, error) {
	return fmt.Sprintf("/Users/%v/Music/iTunes/iTunes Music Library.xml", os.Getenv("USER")), nil
}

func trimTrackLocationPrefix(path string) string {
	return strings.TrimPrefix(path, "file://localhost")
}
