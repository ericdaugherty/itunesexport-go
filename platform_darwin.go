package main

import (
	"fmt"
	"os"
	"strings"
)

func isWindows() bool {
	return false
}

func defaultLibraryPath() string {
	return fmt.Sprintf("/Users/%v/Music/iTunes/iTunes Music Library.xml", os.Getenv("USER"))
}

func trimTrackLocationPrefix(path string) string {
	return strings.TrimPrefix(path, "file://localhost")
}
