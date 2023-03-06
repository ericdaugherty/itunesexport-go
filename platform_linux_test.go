package main

import (
	"errors"
	"os"
	"testing"
)

func TestGetDefaultLibraryInWsl(t *testing.T) {

	cnt := 0

	fakeExecCmdFunc := func(_ string) (string, error) {
		switch (cnt) {
		case 0:
			cnt++
			return "C:", nil
		case 1:
			cnt++
			return "\\Users\\SomeUser", nil
		default:
			return "", errors.New("function called too often")
		}
	}

	os.Setenv("WSLENV", "FOO")
	result, err := defaultLibraryPathInternal(fakeExecCmdFunc)
	if (err != nil) {
		t.Fail()
		t.Logf("function return error")
	}
	
	expected := "/mnt/c/Users/SomeUser/Music/iTunes/iTunes Music Library.xml"
	if (result != expected) {
		t.Fail()
		t.Logf("expected %v, got %v", expected, result)
	}
}

func TestGetDefaultLibraryInLinux(t *testing.T) {

	stubFunc := func(_ string) (string, error) {
		return "", nil
	}

	os.Unsetenv("WSLENV")
	result, err := defaultLibraryPathInternal(stubFunc)
	if (err != nil) {
		t.Fail()
		t.Logf("function return error")
	}
	
	expected := "/mnt/itunes"
	if (result != expected) {
		t.Fail()
		t.Logf("expected %v, got %v", expected, result)
	}
}

