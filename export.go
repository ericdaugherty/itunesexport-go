package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type playlistWriter func(io.Writer, *ExportSettings, *Playlist) error
type trackWriter func(io.Writer, *ExportSettings, *Playlist, *Track, string) error

type ExportSettings struct {
	Library    *Library
	Playlists  []Playlist
	OutputPath string
	Extension  string
}

func ExportPlaylists(exportSettings *ExportSettings) error {

	for _, playlist := range exportSettings.Playlists {
		fmt.Printf("Exporting Playlist %v\n", playlist.Name)

		fileName := filepath.Join(exportSettings.OutputPath, playlist.Name+"."+exportSettings.Extension)

		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer file.Close()
		if err != nil {
			return err
		}

		header, entry, footer := m3uPlaylistWriters()

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
	}

	fmt.Printf("\n\nExport Complete.\n")

	return nil
}

func m3uPlaylistWriters() (header playlistWriter, entry trackWriter, footer playlistWriter) {

	const headerString = "# M3U Playlist '%v' exported %v by iTunes Export v. %v (http://www.ericdaugherty.com/dev/itunesexport/)\n"
	const entryString = "%v\n"

	header = func(w io.Writer, exportSettings *ExportSettings, playlist *Playlist) error {
		_, err := w.Write([]byte(fmt.Sprintf(headerString, playlist.Name, time.Now().Format("2006-01-02 3:04PM"), Version)))
		return err
	}

	entry = func(w io.Writer, exporterSetting *ExportSettings, playlist *Playlist, track *Track, fileLocation string) error {
		_, err := w.Write([]byte(fmt.Sprintf(entryString, fileLocation)))
		return err
	}

	footer = func(w io.Writer, exportSettings *ExportSettings, playlist *Playlist) error {
		return nil
	}

	return
}
