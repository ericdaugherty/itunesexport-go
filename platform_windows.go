package main

import (
	"fmt"
	"os"
)

func isWindows() bool {
	return true
}

func defaultLibraryPath() string {
	return fmt.Sprintf("%v%v\\Music\\iTunes\\iTunes Music Library.xml", os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
}
