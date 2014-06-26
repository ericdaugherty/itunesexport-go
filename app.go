package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	UsageMessage = `usage: %v [<flags>] [include <playlist name>...]
	
Flags:
    -library <file path>        Path to iTunes Music Libary XML File.
    -includeAll                 Include all user defined playlists.
    -includeAllWithBuiltin      Include All playlists, including iTunes defined playlists

`
	UsageErrorMessage = `Unable to parse command line parameters.
%v
`
	ModeUnknown = 0
	ModeInclude = 1
)

// compile passing -ldflags "-X main.Build <build number>"
// Must be var not const so it can be set by build flags.
var Version string = "DEV"

var (
	commandLineError        = false
	commandLineErrorMessage = ""

	libraryPath                    string
	includeAllPlaylists            bool
	includeAllWithBuiltinPlaylists bool
	includePlaylistNames           []string
)

func main() {

	fmt.Printf("\niTunes Export (Go Version %v)\nSee http://www.ericdaugherty.com/dev/itunesexport/ for detailed usage instructions.\n\n", Version)

	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	//TODO Remove default Library Path
	flags.StringVar(&libraryPath, "library", "/Users/eric/Music/iTunes/iTunes Music Library.xml", "location of the iTunes Library XML file.")
	flags.BoolVar(&includeAllPlaylists, "includeAll", false, "includes all user defined playlists.")
	flags.BoolVar(&includeAllWithBuiltinPlaylists, "includeAllWithBuiltin", false, "includes all playlists in the export, including built in iTunes Playlists.")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		commandLineError = true
		commandLineErrorMessage = err.Error()
	}

	var mode = ModeUnknown
	for _, flagValue := range flags.Args() {
		switch flagValue {
		case "include":
			mode = ModeInclude
		default:
			switch mode {
			case ModeUnknown:
				commandLineError = true
				commandLineErrorMessage = fmt.Sprintf("Unexpected paramter %v\n", flagValue)
				break
			case ModeInclude:
				includePlaylistNames = append(includePlaylistNames, flagValue)
			}
		}
	}

	if commandLineError {
		fmt.Printf(UsageMessage, "itunesexport")
		fmt.Printf(UsageErrorMessage, commandLineErrorMessage)
	}

	fmt.Println("Loading Library:", libraryPath)

	library, err := LoadLibrary(libraryPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Library loaded successfully with %v playlists and %v tracks.\n", len(library.Playlists), len(library.Tracks))

	var exportSettings ExportSettings
	if includeAllPlaylists {
		var playlists []Playlist
		for _, value := range library.Playlists {
			if value.DistinguishedKind == 0 && value.Name != "Library" {
				playlists = append(playlists, value)
			}
		}
		exportSettings.Playlists = playlists
	} else if includeAllWithBuiltinPlaylists {
		exportSettings.Playlists = library.Playlists
	}

	if len(includePlaylistNames) > 0 {
		var playlists []Playlist

		for _, playlistName := range includePlaylistNames {
			playlist, ok := library.PlaylistMap[playlistName]
			if ok {
				playlists = append(playlists, playlist)
			} else {
				fmt.Printf("Unable to find matching playlist for name: %v  Skipping Playlist.\n", playlistName)
			}

		}
		exportSettings.Playlists = append(exportSettings.Playlists, playlists...)
	}

	fmt.Printf("Exporting %v playlists...\n", len(exportSettings.Playlists))
}
