package main

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestExportPlaylists(t *testing.T) {
	// arrange
	outputDir := createTempDir(t, "itunes-exporter-test")
	defer os.RemoveAll(outputDir)

	musicFile := createTempFile(t, "Some_Song_*.mp3")
	defer os.Remove(musicFile)
	
	uniqueContent := currentUnixNano()
	// write some unique data to the original music file to later check if the copy was successful
	writeFile(t, musicFile, uniqueContent)
	
	// make sure we have a path only containing "/" as separators
	musicFilePath := filepath.ToSlash(musicFile)
	itunesDbFile := prepareItunesDbFile(t, musicFilePath)
	defer os.Remove(itunesDbFile)

	// act

	// Save the real os.Args and defer the restoration.
	realArgs := os.Args
	defer func() { os.Args = realArgs }()

	// Set the necessary parameters to simulate command line arguments.
	os.Args = []string{
		"itunesexport", // The program name (os.Args[0]).
		"-library", itunesDbFile,
		"-output", outputDir,
		"-type", "M3U",
		"-includeAll",
		"-copy", "PLAYLIST",
	}
	main()

	// assert
	expectedPlaylistDir := filepath.Join(outputDir, "My Playlist")
	assertPathExists(t, expectedPlaylistDir)

	expectedCopiedMusicFilePath := filepath.Join(expectedPlaylistDir, filepath.Base(musicFilePath))
	assertPathExists(t, expectedCopiedMusicFilePath)

	copiedMusicFileContent := readFile(t, expectedCopiedMusicFilePath)
	if copiedMusicFileContent != uniqueContent {
		t.Errorf("Content of copied file not as expected. Expected: %s, Got: %s", uniqueContent, copiedMusicFileContent)
	}

	expectedPlaylistFilePath := filepath.Join(outputDir, "My Playlist.m3u")
	assertPlaylistFileCorrectlyWritten(t, expectedPlaylistFilePath, expectedCopiedMusicFilePath)
}


func TestExportPlaylistsWithAdjustedMusicPath(t *testing.T) {
	// arrange
	outputDir := createTempDir(t, "itunes-exporter-test")
	defer os.RemoveAll(outputDir)

	originalMusicFile := createTempFile(t, "Some_Song_*.mp3")
	defer os.Remove(originalMusicFile)

	uniqueContent := currentUnixNano()
	// write some unique data to the original music file to later check if the copy was successful
	writeFile(t, originalMusicFile, uniqueContent)

	originalMusicFileDir := filepath.Dir(originalMusicFile)
	originalMusicName := filepath.Base(originalMusicFile)
	invalidMusicFilePath := filepath.ToSlash(filepath.Join("/invalid", "path", originalMusicName))
		
	itunesDbFile := prepareItunesDbFile(t, invalidMusicFilePath)
	defer os.Remove(itunesDbFile)

	// act

	// Save the real os.Args and defer the restoration.
	realArgs := os.Args
	defer func() { os.Args = realArgs }()

	// Set the necessary parameters to simulate command line arguments.
	os.Args = []string{
		"itunesexport", // The program name (os.Args[0]).
		"-library", itunesDbFile,
		"-output", outputDir,
		"-type", "M3U",
		"-includeAll",
		"-copy", "PLAYLIST",
		"-musicPath", originalMusicFileDir,
		"-musicPathOrig", "/invalid/path",
	}
	main()

	// assert
	expectedPlaylistDir := filepath.Join(outputDir, "My Playlist")
	assertPathExists(t, expectedPlaylistDir)

	expectedCopiedMusicFilePath := filepath.Join(expectedPlaylistDir, originalMusicName)
	assertPathExists(t, expectedCopiedMusicFilePath)

	copiedMusicFileContent := readFile(t, expectedCopiedMusicFilePath)
	if copiedMusicFileContent != uniqueContent {
		t.Errorf("Content of copied file not as expected. Expected: %s, Got: %s", uniqueContent, copiedMusicFileContent)
	}

	expectedPlaylistFilePath := filepath.Join(outputDir, "My Playlist.m3u")
	assertPlaylistFileCorrectlyWritten(t, expectedPlaylistFilePath, expectedCopiedMusicFilePath)
}


func assertPathExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("File '%s' does not exist: %v", path, err)
	}
}

func createTempFile(t *testing.T, pattern string) string {
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatal(err)
	}
	return tmpFile.Name()
}

func createTempDir(t *testing.T, pattern string) string {
	tmpDirPath, err := os.MkdirTemp("", pattern)
	if err != nil {
		t.Fatal(err)
	}
	return tmpDirPath
}

func writeFile(t *testing.T, filePath string, content string) {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write '%s' to file '%s': %v", content, filePath, err)
	}
}

func readFile(t *testing.T, filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file '%s': %v", filePath, err)
	}
	return string(content)
}

func currentUnixNano() string {
	currentTimeNano := time.Now().UnixNano()
	return strconv.FormatInt(currentTimeNano, 10)
}

func prepareItunesDbFile(t *testing.T, musicFilePath string) string {
	itunesDbContent := readFile(t, "fixture/example-itunes-db.xml")
	itunesDbContentAdjusted := strings.ReplaceAll(string(itunesDbContent), "REPLACE_ME_EXAMPLE_SONG_LOCATION", "file://"+musicFilePath)

	itunesDbFile := createTempFile(t, "testItunesDb_*.xml")
	writeFile(t, itunesDbFile, itunesDbContentAdjusted)

	return itunesDbFile
}

func assertPlaylistFileCorrectlyWritten(t *testing.T, playlistPath string, singleLineContent string) {
	assertPathExists(t, playlistPath)

	playlistFileContents := readFile(t, playlistPath)
	re := buildStringOnSingleLineRegex(singleLineContent)
	matches := re.FindAllString(string(playlistFileContents), -1)

	if len(matches) != 1 {
		t.Errorf("Expected playlist to contain '%s' exactly once, but found %d", singleLineContent, len(matches))
	}
}

// e.g. ...\n/path/to/file.mp3\n
func buildStringOnSingleLineRegex(s string) *regexp.Regexp {
	pattern := "\r?\n" + regexp.QuoteMeta(s) + "\r?\n"
	return regexp.MustCompile(pattern)
}
