package main

import (
	"testing"
	"os"
	"path"
	"bufio"
)

func TestItCreatesFileWhenKnownReleasesFileMissing(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{
		{DummyFileInfo{"Mouldy Old Dough", true}, DummyFileInfo{"Lieutenant Pigeon", true}},
	}, knownReleasesFile(), 0)
	_, err := os.Stat(knownReleasesFile())
	if (err != nil) {
		t.Error("Expected UpdateKnownReleases to create known releases file but it didn't")
	}

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "known releases on creation", 1, len(lines))
	expect(t, "known release", "Lieutenant Pigeon - Mouldy Old Dough", lines[len(lines)-1])
}

func TestItDoesNotChangeFileWhenNoNewReleases(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, DummyFileInfo{"Abba", true}},
		{DummyFileInfo{"It's Fan-dabi-dozi!", true}, DummyFileInfo{"The Krankies", true}},
		{DummyFileInfo{"Discipline", true}, DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 3, len(lines))
	expect(t, "known release 1", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-1])
	expect(t, "known release 2", "Throbbing Gristle - Discipline", lines[len(lines)-2])
	expect(t, "known release 3", "Abba - I Do, I Do, I Do, I Do, I Do", lines[len(lines)-3])
}

func TestItAddsReleasesToEndOfFile(t *testing.T) {
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
	expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	expect(t, "known release 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
}

func TestItIgnoresReleasesAlreadyKnown(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Iconoclastic Diaries", true}, DummyFileInfo{"Shake", true}},
		{DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, DummyFileInfo{"Abba", true}},
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
		{DummyFileInfo{"Discipline", true}, DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 5, len(lines))
	expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	expect(t, "known release 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
}

func TestItHandlesWhenKnownReleasesFileMayNotEndInNewline(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "knownreleases_endwithoutnewline"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	expectInt(t, "number of known releases", 2, len(lines))
	expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	expect(t, "known release 2", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-2])
}

func TestUpdateKnownReleasesReturnValueWhenNoNewReleasesAndKnownReleasesAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	latestAdditions := UpdateKnownReleases([]FolderInfo{}, knownReleasesFile(), 2)

	expectInt(t, "number of latest additions", 2, len(latestAdditions))
	expect(t, "latest addition 1", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[0])
	expect(t, "latest addition 2", "Throbbing Gristle - Discipline", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewReleasesAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Iconoclastic Diaries", true}, DummyFileInfo{"Shake", true}},
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
		{DummyFileInfo{"Mouldy Old Dough", true}, DummyFileInfo{"Lieutenant Pigeon", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), 2)

	expectInt(t, "number of latest additions", 2, len(latestAdditions))
	expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	expect(t, "latest addition 2", "Lieutenant Pigeon - Mouldy Old Dough", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewReleasesBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), 2)

	expectInt(t, "number of latest additions", 2, len(latestAdditions))
	expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewAndKnownReleasesCombinedAreBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	copyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{DummyFileInfo{"Vent", true}, DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), 5)

	expectInt(t, "number of latest additions", 4, len(latestAdditions))
	expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
	expect(t, "latest addition 3", "Throbbing Gristle - Discipline", latestAdditions[2])
	expect(t, "latest addition 4", "Abba - I Do, I Do, I Do, I Do, I Do", latestAdditions[3])
}

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