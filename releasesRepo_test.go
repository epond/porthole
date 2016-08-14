package main

import (
	"testing"
	"os"
	"path"
)

func TestWhenFileMissingItCreatesFile(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{{DummyFileInfo{"Lieutenant Pigeon", false}, nil}}, knownReleasesFile(), 0)
	_, err := os.Stat(knownReleasesFile())
	if (err != nil) {
		t.Error("Expected UpdateKnownReleases to create known releases file but it didn't")
	}
}

//func TestItAddsReleasesToEndOfFile(t *testing.T) {}

func setUp() {
	os.Mkdir(tempDir(), os.ModePerm)
}
func tearDown() {
	os.RemoveAll(tempDir())
}

func knownReleasesFile() string {
	return path.Join(tempDir(), "knownreleases")
}

func testData() string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata")
}

func tempDir() string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/temp")
}