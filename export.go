package main

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
)

type playlistWriter func(io.Writer, *ExportSettings, *Playlist) error
type trackWriter func(io.Writer, *ExportSettings, *Playlist, *Track, string) error

type ExportSettings struct {
	Library    *Library
	Playlists  []Playlist
	ExportType int
	OutputPath string
	Extension  string
	CopyType   int
}

func ExportPlaylists(exportSettings *ExportSettings, library *Library) error {

	for _, playlist := range exportSettings.Playlists {
		fmt.Printf("Exporting Playlist %v\n", playlist.Name)

		fileName := filepath.Join(exportSettings.OutputPath, playlist.Name+"."+exportSettings.Extension)

		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer file.Close()
		if err != nil {
			return err
		}

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
			return errors.New("Export Type Not Implemented")
		}

		// Write out the Header
		err = header(file, exportSettings, &playlist)
		if err != nil {
			return err
		}

		// Write the body of the playlist
		for _, track := range playlist.Tracks(exportSettings.Library) {
			fileLocation, errParse := url.QueryUnescape(track.Location)
			fileLocation = strings.TrimPrefix(fileLocation, "file://localhost")
			if isWindows() {
				fileLocation = strings.TrimPrefix(fileLocation, "/")
			}
			if errParse != nil {
				fmt.Printf("Skipping Track %v because an error occured parsing the location: %v\n", track.Name, errParse.Error())
				continue
			}

			// TODO: Parse location here and pass it in to method.
			err := entry(file, exportSettings, &playlist, &track, fileLocation)
			if err != nil {
				return err
			}
		}

		// Write the footer.
		err = footer(file, exportSettings, &playlist)
		if err != nil {
			return err
		}

		// Copy the tracks (if needed)
		if exportSettings.CopyType != COPY_NONE {
			copyTracks(exportSettings, library, &playlist)
		}

	}

	fmt.Printf("\n\nExport Complete.\n")

	return nil
}

func copyTracks(exportSettings *ExportSettings, library *Library, playlist *Playlist) {

	var destinationPath = ""

	switch exportSettings.CopyType {
	case COPY_PLAYLIST:
		destinationPath = exportSettings.OutputPath + string(os.PathSeparator) + playlist.Name
	}

	for _, item := range playlist.Tracks(library) {
		src, err := url.QueryUnescape(trimTrackLocation(item.Location))
		if err != nil {
			fmt.Printf("Error copying source file.  Unable to decode location: %v\n", item.Location)
			continue
		}
		dest := destinationPath + string(os.PathSeparator) + filepath.Base(item.Location)

		err = copyFile(src, dest)
		if err != nil {
			fmt.Printf("Error copying source file %v.  %v\n", src, err.Error())
		}
	}
}

func copyFile(src, dest string) error {

	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileInfo.Mode().IsRegular() {
		errors.New("Source file is not a regular file.")
	}

	return errors.New("Not Implemented")
}
