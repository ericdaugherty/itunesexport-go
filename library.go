package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	plist "howett.net/plist"
)

// Prevent illegal characters in playlist filenames
var illegalChars = regexp.MustCompile(`[\[\]\\:/*?<>|]`)

type Library struct {
	MajorVersion        int `plist:"Major Version"`
	MinorVersion        int `plist:"Minor Version"`
	Date                time.Time
	ApplicationVersion  int
	Features            int
	ShowContentRating   bool   `plist:"Show Content Ratings"`
	MusicFolder         string `plist:"Music Folder"`
	LibraryPersistentId string `plist:"Library Persistent ID"`
	Tracks              map[string]Track
	Playlists           []Playlist
	PlaylistMap         map[string]Playlist
	PlaylistIdMap       map[string]Playlist
}

type Track struct {
	TrackId             int `plist:"Track ID"`
	Name                string
	Artist              string
	AlbumArtist         string `plist:"Album Artist"`
	Composer            string
	Album               string
	Genre               string
	Kind                string
	Size                int
	TotalTime           int `plist:"Total Time"`
	StartTime           int `plist:"Start Time"`
	StopTime            int `plist:"Stop Time"`
	TrackNumber         int `plist:"Track Number"`
	TrackCount          int `plist:"Track Count"`
	DiscNumber          int `plist:"Disc Number"`
	DiscCount           int `plist:"Disc Count"`
	Year                int
	DateModified        time.Time `plist:"Date Modified"`
	DateAdded           time.Time `plist:"Date Added"`
	BitRate             int       `plist:"Bit Rate"`
	SampleRate          int       `plist:"Sample Rate"`
	PlayCount           int       `plist:"Play Count"`
	PlayDate            int       `plist:"Play Date"`
	PlayDateUTC         time.Time `plist:"Play Date UTC"`
	SkipCount           int       `plist:"Skip Count"`
	SkipDate            time.Time `plist:"Skip Date"`
	Rating              int
	AlbumRating         int    `plist:"Album Rating"`
	AlbumRatingComputed bool   `plist:"Album Rating Computed"`
	ArtworkCount        int    `plist:"Artwork Count"`
	PersistentId        string `plist:"Persistent ID"`
	TrackType           string `plist:"Track Type"`
	Location            string
	FileFolderCount     int `plist:"File Folder Count"`
	LibraryFolderCount  int `plist:"Library Folder Count"`
	Loved               bool
	Disabled            bool
	Comments            string
	SortName            string `plist:"Sort Name"`
	SortAlbum           string `plist:"Sort Album"`
	SortAlbumArtist     string `plist:"Sort Album Artist"`
	SortArtist          string `plist:"Sort Artist"`
	SortComposer        string `plist:"Sort Composer"`
	Work                string
	Grouping            string
	VolumeAdjustment    int `plist:"Volume Adjustment"`
}

type Playlist struct {
	Name                 string
	Master               bool
	PlaylistId           int    `plist:"Playlist ID"`
	PlaylistPersistentId string `plist:"Playlist Persistent ID"`
	ParentPersistentId   string `plist:"Parent Persistent ID"`
	DistinguishedKind    int    `plist:"Distinguished Kind"`
	Visible              bool
	AllItems             bool           `plist:"All Items"`
	Folder               bool           `plist:"Folder"`
	SmartInfo            []byte         `plist:"Smart Info"`
	SmartCriteria        []byte         `plist:"Smart Criteria"`
	PlaylistItems        []PlaylistItem `plist:"Playlist Items"`
}

func (p Playlist) SafeName() string {
	return illegalChars.ReplaceAllString(p.Name, "_")
}

type PlaylistItem struct {
	TrackId int `plist:"Track ID"`
}

func LoadLibrary(fileLocation string) (*Library, error) {
	if _, statErr := os.Stat(fileLocation); os.IsNotExist(statErr) {
		return nil, statErr
	}

	file, pathErr := os.Open(fileLocation)
	if pathErr != nil {
		return nil, pathErr
	}

	decoder := plist.NewDecoder(file)

	var library Library
	decodeErr := decoder.Decode(&library)
	if decodeErr != nil {
		return nil, decodeErr
	}

	library.PlaylistMap = make(map[string]Playlist)
	library.PlaylistIdMap = make(map[string]Playlist)
	for _, value := range library.Playlists {
		library.PlaylistMap[value.Name] = value
		library.PlaylistIdMap[value.PlaylistPersistentId] = value
	}

	// The iTunes library file does not encode the + character correctly in Location
	//so if we do not escape it here, it will get removed later.
	for k, v := range library.Tracks {
		v.Location = strings.ReplaceAll(v.Location, "+", "%2B")
		library.Tracks[k] = v
	}

	return &library, nil
}

func (playlist *Playlist) Tracks(library *Library) []Track {
	var tracks []Track
	for _, item := range playlist.PlaylistItems {
		track, ok := library.Tracks[strconv.FormatInt(int64(item.TrackId), 10)]
		if ok {
			tracks = append(tracks, track)
		}
	}
	return tracks
}
