package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyPlaylistTrack(t *testing.T) {
	sourceFileLocation := createSourceFile()
	defer os.Remove(sourceFileLocation)

	outputPath := createOutputPath()
	defer os.RemoveAll(outputPath)

	exportSettings := ExportSettings{
		OutputPath:        outputPath,
		CopyType:          COPY_PLAYLIST,
		OriginalMusicPath: "",
		NewMusicPath:      "",
	}

	playlist := Playlist{
		Name: "MyPlaylist",
	}

	destinationPath, err := copyTrack(nil, &exportSettings, &playlist, nil, sourceFileLocation)

	if err != nil {
		t.Fatalf("copyTrack failed with error: %s", err)
	}

	sourceFileName := filepath.Base(sourceFileLocation)
	expectedPath := filepath.Join(outputPath, "MyPlaylist", sourceFileName)

	if expectedPath != destinationPath {
		t.Fatalf("destination path '%s' not identical to expected path '%s'", destinationPath, expectedPath)
	}

	if !pathExists(destinationPath) {
		t.Fatalf("File can't be found at destination path: %s", destinationPath)
	}
}

func pathExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}

func createSourceFile() string {
	sourceFile, err := os.CreateTemp("", "sourceFile")
	if err != nil {
		log.Fatal(err)
	}
	return sourceFile.Name()
}

func createOutputPath() string {
	outputPath, err := os.MkdirTemp("", "outputPath")
	if err != nil {
		log.Fatal(err)
	}
	return outputPath
}
