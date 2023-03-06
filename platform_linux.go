package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// as iTunes does not nativly run under Linux, 
// we assume the drive was mounted to this path
const DefaultLinuxDrive = "/mnt/itunes"

func defaultLibraryPath() (string, error) {
	return defaultLibraryPathInternal(execCmd)
}


func defaultLibraryPathInternal(cmdExecFunc func(command string) (string, error)) (string, error) {
	if (os.Getenv("WSLENV") != "") {
		return determineWslDefaultLibraryPath(cmdExecFunc)
	} else {
		return DefaultLinuxDrive, nil
	}
}

func determineWslDefaultLibraryPath(execCmdFunc func(command string) (string, error)) (string, error) {
	homeDrive, err := execCmdFunc("echo %HOMEDRIVE%")
	if err != nil {
		return "", err
	}
	homeDrive = strings.ToLower(strings.TrimRight(homeDrive, ":"))

	homePath, err := execCmdFunc("echo %HOMEPATH%")
	if err != nil {
		return "", err
	}
	homePath = strings.ReplaceAll(homePath, "\\", "/")

	return fmt.Sprintf("/mnt/%v%v/Music/iTunes/iTunes Music Library.xml", homeDrive, homePath), nil
}

func execCmd(command string) (string, error) {
	result, err := exec.Command("cmd.exe", "/c", command).Output()
	if (err != nil) {
		return "", err;
	}
	return strings.TrimSpace(string(result)), nil

}

func trimTrackLocationPrefix(path string) string {
	return strings.TrimPrefix(path, "file://localhost")
}
