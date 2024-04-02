package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestExportPlaylists tests the main functionality with a given XML library file.
func TestExportPlaylists(t *testing.T) {
	// Create a temporary directory to serve as the output directory.
	outputDir, err := ioutil.TempDir("", "itunes_test_output")
	if err != nil {
		t.Fatalf("Failed to create temporary output directory: %v", err)
	}
	defer os.RemoveAll(outputDir) // Clean up after the test.

	// Create a temporary empty file to simulate the music file.
	musicFile, err := ioutil.TempFile("", "Some_Song_*.mp3")
	if err != nil {
		t.Fatalf("Failed to create temporary music file: %v", err)
	}
	defer os.Remove(musicFile.Name()) // Clean up after the test.

	currentTimeNano := time.Now().UnixNano()
	unixNanoStr := strconv.FormatInt(currentTimeNano, 10)

	// Write the Unix nanoseconds string to the music file.
	err = ioutil.WriteFile(musicFile.Name(), []byte(unixNanoStr), 0644)
	if err != nil {
		t.Fatalf("Failed to write Unix nanoseconds to music file: %v", err)
	}

	// Read the XML content from the fixture file.
	xmlContent, err := ioutil.ReadFile("fixture/example-itunes-db.xml")
	if err != nil {
		t.Fatalf("Failed to read fixture XML file: %v", err)
	}

	// Adjust the XML content to point to the temporary music file's location.
	musicFilePath := filepath.ToSlash(musicFile.Name())
	xmlContentAdjusted := strings.ReplaceAll(string(xmlContent), "REPLACE_ME_EXAMPLE_SONG_LOCATION", "file://"+musicFilePath)

	// Write the adjusted XML content to a new temporary fixture file.
	fixtureFile, err := ioutil.TempFile("", "testItunesDb_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temporary fixture XML file: %v", err)
	}
	defer os.Remove(fixtureFile.Name()) // Clean up after the test.

	_, err = fixtureFile.WriteString(xmlContentAdjusted)
	if err != nil {
		t.Fatalf("Failed to write to temporary fixture XML file: %v", err)
	}

	// Close the file to ensure all writes are flushed to disk.
	if err := fixtureFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary fixture XML file: %v", err)
	}

	// Save the real os.Args and defer the restoration.
	realArgs := os.Args
	defer func() { os.Args = realArgs }()

	// Set the necessary parameters to simulate command line arguments.
	os.Args = []string{
		"itunesexport", // The program name (os.Args[0]).
		"-library", fixtureFile.Name(),
		"-output", outputDir,
		"-type", "M3U",
		"-includeAll",
		"-copy", "PLAYLIST",
	}

	// Run the main program to test the functionality.
	main()

	expectedPlaylistPath := filepath.Join(outputPath, "My Playlist")

	// if expectedPlaylistPath != destinationPath {
	// 	t.Fatalf("destination path '%s' not identical to expected path '%s'", destinationPath, expectedPlaylistPath)
	// }

	if !pathExists(expectedPlaylistPath) {
		t.Fatalf("File can't be found at destination path: %s", expectedPlaylistPath)
	}

	musicFileName := filepath.Base(musicFilePath)
	expectedMusicFilePath := filepath.Join(expectedPlaylistPath, musicFileName)
	if !pathExists(expectedMusicFilePath) {
		t.Fatalf("File can't be found at destination path: %s", expectedMusicFilePath)
	}

	musicFileContents, err := ioutil.ReadFile(musicFilePath)
	if err != nil {
		t.Fatalf("Failed to read music file: %v", err)
	}

	// Assert that the music file contains the Unix nanoseconds string.
	if string(musicFileContents) != unixNanoStr {
		t.Errorf("The music file does not contain the expected Unix nanoseconds string. Expected: %s, Got: %s", unixNanoStr, string(musicFileContents))
	}

	playlistFilePath := filepath.Join(outputDir, "My Playlist.m3u")

	// Read the contents of the "My Playlist.m3u" file.
	playlistFileContents, err := ioutil.ReadFile(playlistFilePath)
	if err != nil {
		t.Fatalf("Failed to read playlist file: %v", err)
	}

	pattern := "\r?\n" + regexp.QuoteMeta(expectedMusicFilePath) + "\r?\n"
	re := regexp.MustCompile(pattern)

	// Find all matches of the pattern in the playlist file contents.
	matches := re.FindAllString(string(playlistFileContents), -1)

	// Assert that there is exactly one match.
	if len(matches) != 1 {
		t.Errorf("Expected exactly one line with the music file path, but found %d", len(matches))
	}
}
