package main

import (
	"fmt"
	"os"
	"strings"
)

func defaultLibraryPath() (string, error) {
	return fmt.Sprintf("%v%v\\Music\\iTunes\\iTunes Music Library.xml", os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH")), nil
}

func trimTrackLocationPrefix(path string) string {
	return strings.TrimPrefix(path, "file://localhost/")
}
