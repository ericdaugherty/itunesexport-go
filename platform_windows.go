package main

import (
	"fmt"
	"os"
	"strings"
)

func isWindows() bool {
	return true
}

func defaultLibraryPath() string {
	return fmt.Sprintf("%v%v\\Music\\iTunes\\iTunes Music Library.xml", os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
}

func trimTrackLocation(path string) string {
	return strings.TrimPrefix(path, "file://localhost/"
}