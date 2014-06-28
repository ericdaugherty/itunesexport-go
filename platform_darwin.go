package main

import (
	"fmt"
	"os"
)

func isWindows() bool {
	return false
}

func defaultLibraryPath() string {
	return fmt.Sprintf("/Users/%v/Music/iTunes/iTunes Music Library.xml", os.Getenv("USER"))
}
