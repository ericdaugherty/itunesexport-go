package main

import (
	"fmt"
	"io"
	"time"
)

func m3uPlaylistWriters() (header playlistWriter, entry trackWriter, footer playlistWriter) {

	const headerString = "# M3U Playlist '%v' exported %v by iTunes Export v. %v (http://www.ericdaugherty.com/dev/itunesexport/)\n"
	const entryString = "%v\n"

	header = func(w io.Writer, _ *ExportSettings, playlist *Playlist) error {
		_, err := w.Write([]byte(fmt.Sprintf(headerString, playlist.Name, time.Now().Format("2006-01-02 3:04PM"), Version)))
		return err
	}

	entry = func(w io.Writer, _ *ExportSettings, _ *Playlist, _ *Track, fileLocation string) error {
		_, err := w.Write([]byte(fmt.Sprintf(entryString, fileLocation)))
		return err
	}

	footer = func(_ io.Writer, _ *ExportSettings, _ *Playlist) error {
		return nil
	}

	return
}

func extPlaylistWriters() (header playlistWriter, entry trackWriter, footer playlistWriter) {

	const headerString = "#EXTM3U\n"
	const entryString = "#EXTINF:%v,%v - %v\n%v\n"

	header = func(w io.Writer, _ *ExportSettings, _ *Playlist) error {
		_, err := w.Write([]byte(fmt.Sprint(headerString)))
		return err
	}

	entry = func(w io.Writer, _ *ExportSettings, _ *Playlist, track *Track, fileLocation string) error {
		_, err := w.Write([]byte(fmt.Sprintf(entryString, track.TotalTime/1000, track.Artist, track.Name, fileLocation)))
		return err
	}

	footer = func(_ io.Writer, _ *ExportSettings, _ *Playlist) error {
		return nil
	}

	return
}

func wplPlaylistWriters() (header playlistWriter, entry trackWriter, footer playlistWriter) {

	const headerString = `<?wpl version=\"1.0\"?>
<smil>
  <head>
    <author />
    <title>%v</title>
  </head>
  <body>
    <seq>
`

	const entryString = "      <media src=%v></media>\n"
	const footerString = `    </seq>
  </body>
</smil>
`

	header = func(w io.Writer, _ *ExportSettings, playlist *Playlist) error {
		_, err := w.Write([]byte(fmt.Sprintf(headerString, playlist.Name)))
		return err
	}

	entry = func(w io.Writer, _ *ExportSettings, _ *Playlist, _ *Track, fileLocation string) error {
		_, err := w.Write([]byte(fmt.Sprintf(entryString, fileLocation)))
		return err
	}

	footer = func(w io.Writer, _ *ExportSettings, _ *Playlist) error {
		_, err := w.Write([]byte(footerString))
		return err
	}

	return
}

func zplPlaylistWriters() (header playlistWriter, entry trackWriter, footer playlistWriter) {

	const headerString = `<?zpl version=\"1.0\"?>
<smil>
  <head>
    <meta name="Generator" content="Zune -- 1.3.5728.0" />
    <author />
    <title>%v</title>
  </head>
  <body>
    <seq>
`

	const entryString = "      <media src=%v></media>\n"
	const footerString = `    </seq>
  </body>
</smil>
`

	header = func(w io.Writer, _ *ExportSettings, playlist *Playlist) error {
		_, err := w.Write([]byte(fmt.Sprintf(headerString, playlist.Name)))
		return err
	}

	entry = func(w io.Writer, _ *ExportSettings, _ *Playlist, _ *Track, fileLocation string) error {
		_, err := w.Write([]byte(fmt.Sprintf(entryString, fileLocation)))
		return err
	}

	footer = func(w io.Writer, _ *ExportSettings, _ *Playlist) error {
		_, err := w.Write([]byte(footerString))
		return err
	}

	return
}
