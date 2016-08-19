package main

import (
	"testing"
	"os"
	"path"
	"bufio"
)

func _TestWhenFileMissingItCreatesFile(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{{DummyFileInfo{"Lieutenant Pigeon", true}, nil}}, knownReleasesFile(), 0)
	_, err := os.Stat(knownReleasesFile())
	if (err != nil) {
		t.Error("Expected UpdateKnownReleases to create known releases file but it didn't")
	}
}

func _TestItDoesNotChangeFileWhenNoNewReleases(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, DummyFileInfo{"Abba", true}},
		{DummyFileInfo{"It's Fan-Dabi-Dozi!", true}, DummyFileInfo{"The Krankies", true}},
		{DummyFileInfo{"Discipline", true}, DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 3, len(lines))
	expect(t, "new release 1", "Abba - I Do, I Do, I Do, I Do, I Do", lines[len(lines)-3])
	expect(t, "new release 2", "Throbbing Gristle - Discipline", lines[len(lines)-2])
	expect(t, "new release 3", "The Krankies - It's Fan-Dabi-Dozi!", lines[len(lines)-1])
}

func _TestItAddsReleasesToEndOfFile(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
		{DummyFileInfo{"Iconoclastic Diaries", true}, DummyFileInfo{"Shake", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 5, len(lines))
	expect(t, "new release 1", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
	expect(t, "new release 2", "Daniel Menche - Vent", lines[len(lines)-1])
}

func TestItIgnoresReleasesAlreadyKnown(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
		{DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, DummyFileInfo{"Abba", true}},
		{DummyFileInfo{"Iconoclastic Diaries", true}, DummyFileInfo{"Shake", true}},
		{DummyFileInfo{"Discipline", true}, DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 5, len(lines))
	expect(t, "new release 1", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
	expect(t, "new release 2", "Daniel Menche - Vent", lines[len(lines)-1])
}

//func TestItHandlesWhenKnownReleasesFileMayNotEndInNewline(t *testing.T) {}

func setUp() {
	os.Mkdir(tempDir(), os.ModePerm)
}
func tearDown() {
	os.RemoveAll(tempDir())
}

func knownReleasesFile() string {
	return path.Join(tempDir(), "knownreleases")
}

func knownReleasesLines() (file *os.File, lines []string) {
	file, _ = os.Open(knownReleasesFile())
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return file, lines
}

func testData() string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata")
}

func tempDir() string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/temp")
}