package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	UsageMessage = `usage: %v [<flags>] [include <playlist name>...] [exclude <playlist name>...]

Specify one of the -include<All|AllWithBuiltin|PlaylistWithRegex> flags or use 
the include parameter with playlist names to specify the playlist to export.

Usage of exclude parameter will override any playlist included using the flag 
or parameter.

Flags:
    -library <file path>        Path to iTunes Music Library XML File.
    -output <file path>         Path where the playlists should be written.
    -type <M3U|EXT|WPL|ZPL>     Type of playlist file to write.  Defaults to M3U
                                EXT = M3U Extended, WPL = Windows Playlist, ZPL = Zune Playlist
    -includeAll                 Include all user defined playlists.
    -includeAllWithBuiltin      Include All playlists, including iTunes defined playlists
    -includePlaylistWithRegex   Include all playlists matching the provided regular expression
    -copy <COPY TYPE>           Copy the music tracks as well, according the the COPY TYPE scheme...
        NONE                    (default) The music files will not be copied.	                            
        PLAYLIST                Copies the music into a folder for each playlist.
        ITUNES                  Copies using the itunes music/<Artist>/<Album>/<Track> structure.
        FLAT                    Copies all the music into the output folder.
    -musicPath <new path>       Base path to the music files. This will override the Music Folder path from iTunes.
	-musicPathOrig <path>       When using -musicPath this allows you to override the Music Folder value that is replaced.
`
	UsageErrorMessage = `Unable to parse command line parameters.
%v
`
	ModeUnknown = 0
	ModeInclude = 1
	ModeExclude = 2
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
	includePlaylistWithRegex       string
	excludePlaylistNames           []string
	copyType                       string
	musicPath                      string
	musicPathOrig                  string

	exportSettings ExportSettings
)

func main() {

	fmt.Printf("\niTunes Export (Go Version %v)\nSee http://www.ericdaugherty.com/dev/itunesexport/ for detailed instructions.\n\n", Version)

	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	flags.StringVar(&libraryPath, "library", "", "")
	flags.StringVar(&outputPath, "output", "", "")
	flags.StringVar(&exportType, "type", "M3U", "")
	flags.BoolVar(&includeAllPlaylists, "includeAll", false, "")
	flags.BoolVar(&includeAllWithBuiltinPlaylists, "includeAllWithBuiltin", false, "")
	flags.StringVar(&includePlaylistWithRegex, "includePlaylistWithRegex", "", "")
	flags.StringVar(&copyType, "copy", "NONE", "")
	flags.StringVar(&musicPath, "musicPath", "", "")
	flags.StringVar(&musicPathOrig, "musicPathOrig", "", "")

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
		case "exclude":
			mode = ModeExclude
		default:
			switch mode {
			case ModeUnknown:
				commandLineError = true
				commandLineErrorMessage = fmt.Sprintf("Unexpected paramter %v\n", flagValue)
			case ModeInclude:
				includePlaylistNames = append(includePlaylistNames, flagValue)
			case ModeExclude:
				excludePlaylistNames = append(excludePlaylistNames, flagValue)
			}
		}
	}

	if commandLineError {
		fmt.Printf(UsageMessage, "itunesexport")
		fmt.Printf(UsageErrorMessage, commandLineErrorMessage)
		return
	}

	if libraryPath == "" {
		libraryPath, err = defaultLibraryPath()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	libraryPath = filepath.Clean(libraryPath)

	fmt.Printf("Include: %v, Exclude %v", includePlaylistNames, excludePlaylistNames)

	fmt.Println("Loading Library:", libraryPath)
	library, err := LoadLibrary(libraryPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	exportSettings.Library = library
	fmt.Printf("Library loaded successfully with %v playlists and %v tracks.\n", len(library.Playlists), len(library.Tracks))

	if musicPath != "" {
		if musicPathOrig != "" {
			exportSettings.OriginalMusicPath = musicPathOrig
		} else {
			origMusicPath, err := url.QueryUnescape(library.MusicFolder)
			if err != nil {
				fmt.Printf("Error parsing Music Folder from library: %v\n", err)
				return
			}
			exportSettings.OriginalMusicPath = trimTrackLocationPrefix(origMusicPath)
		}
	}
	exportSettings.NewMusicPath = musicPath

	exportSettings.OutputPath = outputPath
	exportSettings.Playlists = parsePlaylists(exportSettings.Library)

	fmt.Printf("Exporting %v playlists...\n", len(exportSettings.Playlists))
	err = ExportPlaylists(&exportSettings, library)
	if err != nil {
		fmt.Printf("Error Exporting Playlist: %v\n", err)
		return
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
	case "FLAT":
		exportSettings.CopyType = COPY_FLAT
	default:
		return errors.New("Unknown Copy Type: " + copyType)
	}
	return nil
}

func parsePlaylists(library *Library) []Playlist {
	var playlists []Playlist

	if includeAllPlaylists {
		for _, playlist := range library.Playlists {
			if playlist.DistinguishedKind == 0 && playlist.Name != "Library" {
				playlists = append(playlists, playlist)
			}
		}
	} else if includeAllWithBuiltinPlaylists {
		playlists = library.Playlists
	} else if len(includePlaylistWithRegex) > 0 {
		for _, playlist := range library.Playlists {
			match, _ := regexp.MatchString(includePlaylistWithRegex, playlist.Name)
			if match {
				playlists = append(playlists, playlist)
			}
		}
	} else if len(includePlaylistNames) > 0 {
		for _, playlistName := range includePlaylistNames {
			playlist, ok := library.PlaylistMap[playlistName]
			if ok {
				playlists = append(playlists, playlist)
			} else {
				fmt.Printf("Unable to find matching playlist for name: %q. Skipping Playlist.\n", playlistName)
			}
		}
	}

	var filteredPlaylists []Playlist
	for _, playlist := range playlists {
		remove := false
		for _, removePlaylistName := range excludePlaylistNames {
			if playlist.Name == removePlaylistName {
				remove = true
				break
			}
		}
		if !remove {
			filteredPlaylists = append(filteredPlaylists, playlist)
		}
	}

	return filteredPlaylists

}
