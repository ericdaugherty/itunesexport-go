package main

import (
	"flag"
	"fmt"
)

// compile passing -ldflags "-X main.Build <build number>"
// Must be var not const so it can be set by build flags.
var Version string = "DEV"

func main() {

	fmt.Printf("iTunes Export (Go Version %v)\n", Version)

	var libraryPath string

	flag.StringVar(&libraryPath, "library", "", "location of the iTunes Library XML file.")
	flag.Parse()

	// TODO: Remove
	libraryPath = "/Users/eric/Music/iTunes/iTunes Music Library.xml"

	fmt.Println("Library Location:", libraryPath)

	library, err := LoadLibrary(libraryPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Library loaded successfully with %v playlists and %v tracks.\n", len(library.Playlists), len(library.Tracks))
}
