package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	UsageMessage = `usage: %v [<flags>] [include <playlist name>...]
	
Flags:
    -library <file path>        Path to iTunes Music Libary XML File.
    -output <file path>         Path where the playlists should be written.
    -type <M3U|EXT|WPL|ZPL>     Type of playlist file to write.  Defaults to M3U
                                EXT = M3U Extended, WPL = Windows Playlist, ZPL = Zune Playlist
    -includeAll                 Include all user defined playlists.
    -includeAllWithBuiltin      Include All playlists, including iTunes defined playlists
    -copy <COPY TYPE>           Copy the music tracks as well, according the the COPY TYPE scheme...
        PLAYLIST                Copies the music into a folder for each playlist.
        ITUNES                  Copies using the itunes music/<Artist>/<Album>/<Track> structure.

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
	outputPath                     string
	exportType                     string
	includeAllPlaylists            bool
	includeAllWithBuiltinPlaylists bool
	includePlaylistNames           []string
	copyType                       string

	exportSettings ExportSettings
)

func main() {

	fmt.Printf("\niTunes Export (Go Version %v)\nSee http://www.ericdaugherty.com/dev/itunesexport/ for detailed usage instructions.\n\n", Version)

	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	flags.StringVar(&libraryPath, "library", "", "")
	flags.StringVar(&outputPath, "output", "", "")
	flags.StringVar(&exportType, "type", "M3U", "")
	flags.BoolVar(&includeAllPlaylists, "includeAll", false, "")
	flags.BoolVar(&includeAllWithBuiltinPlaylists, "includeAllWithBuiltin", false, "")
	flags.StringVar(&copyType, "copy", "NONE", "")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		commandLineError = true
		commandLineErrorMessage = err.Error()
	}

	err = parseExportType()
	if err != nil {
		commandLineError = true
		commandLineErrorMessage = fmt.Sprintf("%v\n", err.Error())
	}

	err = parseCopyType()
	if err != nil {
		commandLineError = true
		commandLineErrorMessage = fmt.Sprintf("%v\n", err.Error())
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
		return
	}

	if libraryPath == "" {
		libraryPath = defaultLibraryPath()
	}
	libraryPath = filepath.Clean(libraryPath)

	fmt.Println("Loading Library:", libraryPath)
	library, err := LoadLibrary(libraryPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	exportSettings.Library = library
	fmt.Printf("Library loaded successfully with %v playlists and %v tracks.\n", len(library.Playlists), len(library.Tracks))

	exportSettings.OutputPath = outputPath
	exportSettings.Playlists = parsePlaylists(exportSettings.Library)

	fmt.Printf("Exporting %v playlists...\n", len(exportSettings.Playlists))
	err = ExportPlaylists(&exportSettings, library)
	if err != nil {
		fmt.Printf("Error Exporting Playlist: %v\n", err.Error())
	}
}

func parseExportType() error {
	switch strings.ToUpper(exportType) {
	case "M3U":
		exportSettings.ExportType = M3U
		exportSettings.Extension = "m3u"
	case "EXT":
		exportSettings.ExportType = EXT
		exportSettings.Extension = "m3u"
	case "WPL":
		exportSettings.ExportType = WPL
		exportSettings.Extension = "wpl"
	case "ZPL":
		exportSettings.ExportType = ZPL
		exportSettings.Extension = "zpl"
	default:
		return errors.New("Unknown Export Type: " + exportType)
	}
	return nil
}

func parseCopyType() error {
	switch strings.ToUpper(copyType) {
	case "NONE":
		exportSettings.CopyType = COPY_NONE
	case "PLAYLIST":
		exportSettings.CopyType = COPY_PLAYLIST
	case "ITUNES":
		exportSettings.CopyType = COPY_ITUNES
	default:
		return errors.New("Unknown Copy Type: " + copyType)
	}
	return nil
}

func parsePlaylists(library *Library) (playlists []Playlist) {
	if includeAllPlaylists {
		for _, value := range library.Playlists {
			if value.DistinguishedKind == 0 && value.Name != "Library" {
				playlists = append(playlists, value)
			}
		}
		exportSettings.Playlists = playlists
	} else if includeAllWithBuiltinPlaylists {
		playlists = library.Playlists
	}

	if len(includePlaylistNames) > 0 {
		for _, playlistName := range includePlaylistNames {
			playlist, ok := library.PlaylistMap[playlistName]
			if ok {
				playlists = append(playlists, playlist)
			} else {
				fmt.Printf("Unable to find matching playlist for name: %v  Skipping Playlist.\n", playlistName)
			}

		}
	}
	return playlists
}
