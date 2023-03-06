package main

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"regexp"
)

const (
	M3U = iota
	EXT
	WPL
	ZPL
)

const (
	COPY_NONE = iota
	COPY_PLAYLIST
	COPY_ITUNES
	COPY_FLAT
)

type playlistWriter func(io.Writer, *ExportSettings, *Playlist) error
type trackWriter func(io.Writer, *ExportSettings, *Playlist, *Track, string) error

type ExportSettings struct {
	Library           *Library
	Playlists         []Playlist
	ExportType        int
	OutputPath        string
	Extension         string
	CopyType          int
	OriginalMusicPath string
	NewMusicPath      string
}

func ExportPlaylists(exportSettings *ExportSettings, library *Library) error {
	start := time.Now()

	for _, playlist := range exportSettings.Playlists {
		fmt.Printf("Exporting Playlist %v\n", playlist.Name)

		// Prevent illegal characters in playlist filenames
		illegalChars := regexp.MustCompile(`[*?<>|]`)
		safePlaylistName := illegalChars.ReplaceAllString(playlist.Name, "_")
		safePlaylistName = strings.Replace(safePlaylistName, ":", " - ", -1)

		fileName := filepath.Join(exportSettings.OutputPath, safePlaylistName+"."+exportSettings.Extension)

		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		var header playlistWriter
		var entry trackWriter
		var footer playlistWriter
		switch exportSettings.ExportType {
		case M3U:
			header, entry, footer = m3uPlaylistWriters()
		case EXT:
			header, entry, footer = extPlaylistWriters()
		case WPL:
			header, entry, footer = wplPlaylistWriters()
		case ZPL:
			header, entry, footer = zplPlaylistWriters()
		default:
			return errors.New("export type not implemented")
		}

		// Write out the Header
		err = header(file, exportSettings, &playlist)
		if err != nil {
			return err
		}

		// Write the body of the playlist
		for _, track := range playlist.Tracks(exportSettings.Library) {

			sourceFileLocation, errParse := url.QueryUnescape(track.Location)
			sourceFileLocation = trimTrackLocationPrefix(sourceFileLocation)

			destFileLocation, err := copyTrack(exportSettings, &playlist, &track, sourceFileLocation)
			if err != nil {
				fmt.Printf("Unable to copy file %v: %v\n", sourceFileLocation, err.Error())
				continue
			}

			if errParse != nil {
				fmt.Printf("Skipping Track %v because an error occured parsing the location: %v\n", track.Name, errParse.Error())
				continue
			}

			err = entry(file, exportSettings, &playlist, &track, destFileLocation)
			if err != nil {
				return err
			}
		}

		// Write the footer.
		err = footer(file, exportSettings, &playlist)
		if err != nil {
			return err
		}

	}

	fmt.Printf("\n\nExport Complete.\n")
	fmt.Println(time.Since(start).String())
	return nil
}

// copyTrack copies a file from the provided sourceFileLocation to another location. The new location
// depends on the CopyType selected in exportSettings. If COPY_NONE is selected, the sourceFileLocation is returned.
func copyTrack(exportSettings *ExportSettings, playlist *Playlist, track *Track, sourceFileLocation string) (string, error) {
	var destinationPath string
	switch exportSettings.CopyType {
	case COPY_PLAYLIST:
		destinationPath = filepath.Join(exportSettings.OutputPath, playlist.Name)
	case COPY_ITUNES:
		destinationPath = filepath.Join(exportSettings.OutputPath, track.Artist, track.Album)
	case COPY_FLAT:
		destinationPath = exportSettings.OutputPath
	case COPY_NONE:
		if exportSettings.NewMusicPath != "" {
			return strings.Replace(sourceFileLocation, exportSettings.OriginalMusicPath, exportSettings.NewMusicPath, 1), nil
		}
		return sourceFileLocation, nil
	default:
		return "", errors.New("unknown copy type")
	}
	dest := filepath.Join(destinationPath, filepath.Base(sourceFileLocation))

	if err := copyFile(sourceFileLocation, dest); err != nil {
		return "", err
	}
	return dest, nil
}

func copyFile(src, dest string) error {
	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileInfo.Mode().IsRegular() {
		return errors.New("source file is not a regular file.")
	}

	_, err = os.Stat(dest)
	if err == nil {
		// No need to copy.
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	destDir := filepath.Dir(dest)
	_, err = os.Stat(destDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(destDir, 0777)
			if err != nil {
				return nil
			}
		} else {
			return err
		}
	}

	return copyFileData(src, dest)
}

func copyFileData(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
