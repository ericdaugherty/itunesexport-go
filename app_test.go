package main

import (
	"testing"
)

func TestIncludeAllPlaylists(t *testing.T) {
	resetGlobalVars()

	library := &Library{
		Playlists: []Playlist{
			{Name: "Foo"},
			{Name: "Bar"},
			{Name: "Library", DistinguishedKind: 0},
		},
	}

	includeAllPlaylists = true
	playlists := parsePlaylists(library)

	if len(playlists) != 2 {
		t.Fatal("wrong playlist size")
	}
	if !(playlists[0].Name == "Foo" && playlists[1].Name == "Bar") {
		t.Fatal("Unexpected playlist names")
	}
}

func TestIncludeAllWinBuiltinPlaylists(t *testing.T) {
	resetGlobalVars()

	library := &Library{
		Playlists: []Playlist{
			{Name: "Foo"},
			{Name: "Bar"},
			{Name: "Library", DistinguishedKind: 0},
		},
	}

	includeAllWithBuiltinPlaylists = true
	playlists := parsePlaylists(library)

	if len(playlists) != 3 {
		t.Fatal("wrong playlist size")
	}
}

func TestIncludePlaylistNames(t *testing.T) {
	resetGlobalVars()

	library := &Library{
		PlaylistMap: map[string]Playlist{
			"Foo": {Name: "Foo"},
			"Bar": {Name: "Bar"},
		},
	}

	includePlaylistNames = []string{"Bar"}
	playlists := parsePlaylists(library)

	if len(playlists) != 1 && playlists[0].Name != "Bar" {
		t.Fatal("unexpected playlist")
	}
}

func TestPlaylistViaRegex(t *testing.T) {
	resetGlobalVars()

	library := &Library{
		Playlists: []Playlist{
			{Name: "Foo"},
			{Name: "Bar"},
			{Name: "Buzz"},
		},
	}

	includePlaylistWithRegex = "^B+"
	playlists := parsePlaylists(library)

	if len(playlists) != 2 {
		t.Fatal("unexpected playlist size")
	}
	if playlists[0].Name != "Bar" && playlists[1].Name != "Buzz" {
		t.Fatal("unexpected playlists")
	}
}

func TestExcludePlaylists(t *testing.T) {
	resetGlobalVars()

	library := &Library{
		Playlists: []Playlist{
			{Name: "Foo"},
			{Name: "Bar"},
			{Name: "Library", DistinguishedKind: 0},
		},
	}

	includeAllPlaylists = true
	excludePlaylistNames = []string{"Bar"}
	playlists := parsePlaylists(library)

	if len(playlists) != 1 {
		t.Fatal("wrong playlist size")
	}
	if !(playlists[0].Name == "Foo") {
		t.Fatal("Unexpected playlist names")
	}
}

func resetGlobalVars() {
	includeAllPlaylists = false
	includeAllWithBuiltinPlaylists = false
	includePlaylistNames = []string{}
}
