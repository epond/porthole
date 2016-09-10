package music

import (
	"testing"
	"os"
	"path"
	"bufio"
	"github.com/epond/porthole/test"
)

func TestItCreatesFileWhenKnownReleasesFileMissing(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{
		{test.DummyFileInfo{"Mouldy Old Dough", true}, test.DummyFileInfo{"Lieutenant Pigeon", true}},
	}, knownReleasesFile(), knownReleasesBackup(), 0)
	_, lines := knownReleasesLines()
	if (len(lines) == 0) {
		t.Error("Expected UpdateKnownReleases to create known releases file but it didn't")
	}

	file, lines := knownReleasesLines()
	defer file.Close()

	test.ExpectInt(t, "known releases on creation", 1, len(lines))
	test.Expect(t, "known release", "Lieutenant Pigeon - Mouldy Old Dough", lines[len(lines)-1])
}

func TestItDoesNotChangeFileWhenNoNewReleases(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, test.DummyFileInfo{"Abba", true}},
		{test.DummyFileInfo{"It's Fan-dabi-dozi!", true}, test.DummyFileInfo{"The Krankies", true}},
		{test.DummyFileInfo{"Discipline", true}, test.DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	test.ExpectInt(t, "number of known releases", 3, len(lines))
	test.Expect(t, "known release 1", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-1])
	test.Expect(t, "known release 2", "Throbbing Gristle - Discipline", lines[len(lines)-2])
	test.Expect(t, "known release 3", "Abba - I Do, I Do, I Do, I Do, I Do", lines[len(lines)-3])
}

func TestItAddsReleasesToEndOfFile(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	test.ExpectInt(t, "number of known releases", 5, len(lines))
	test.Expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known release 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
}

func TestItIgnoresReleasesAlreadyKnown(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
		{test.DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, test.DummyFileInfo{"Abba", true}},
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Discipline", true}, test.DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	test.ExpectInt(t, "number of known releases", 5, len(lines))
	test.Expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known release 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
}

func TestItHandlesWhenKnownReleasesFileMayNotEndInNewline(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "knownreleases_endwithoutnewline"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 0)

	file, lines := knownReleasesLines()
	defer file.Close()

	test.ExpectInt(t, "number of known releases", 2, len(lines))
	test.Expect(t, "known release 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known release 2", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-2])
}

func TestUpdateKnownReleasesReturnValueWhenNoNewReleasesAndKnownReleasesAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	latestAdditions := UpdateKnownReleases([]FolderInfo{}, knownReleasesFile(), knownReleasesBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[0])
	test.Expect(t, "latest addition 2", "Throbbing Gristle - Discipline", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewReleasesAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Mouldy Old Dough", true}, test.DummyFileInfo{"Lieutenant Pigeon", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "Lieutenant Pigeon - Mouldy Old Dough", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewReleasesBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
}

func TestUpdateKnownReleasesReturnValueWhenNewAndKnownReleasesCombinedAreBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownReleasesFile(), path.Join(testData(), "3knownreleases"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownReleases(currentScan, knownReleasesFile(), knownReleasesBackup(), 5)

	test.ExpectInt(t, "number of latest additions", 4, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
	test.Expect(t, "latest addition 3", "Throbbing Gristle - Discipline", latestAdditions[2])
	test.Expect(t, "latest addition 4", "Abba - I Do, I Do, I Do, I Do, I Do", latestAdditions[3])
}

func TestItBacksUpKnownReleasesWhenChanged(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{
		{test.DummyFileInfo{"Mouldy Old Dough", true}, test.DummyFileInfo{"Lieutenant Pigeon", true}},
	}, knownReleasesFile(), knownReleasesBackup(), 0)

	_, backup := knownReleasesBackupLines()
	if (len(backup) == 0) {
		t.Error("Expected UpdateKnownReleases to create known releases backup but it didn't")
	}
}

func TestItDoesntBackUpKnownReleasesWhenUnchanged(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownReleases([]FolderInfo{}, knownReleasesFile(), knownReleasesBackup(), 0)

	_, backup := knownReleasesBackupLines()
	if (len(backup) > 0) {
		t.Error("Expected UpdateKnownReleases to not create known releases backup but it did")
	}
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

func knownReleasesBackup() string {
	return path.Join(tempDir(), "knownreleases_backup")
}

func knownReleasesLines() (file *os.File, lines []string) {
	return fileLines(knownReleasesFile())
}

func knownReleasesBackupLines() (file *os.File, lines []string) {
	return fileLines(knownReleasesBackup())
}

func fileLines(filePath string) (file *os.File, lines []string) {
	var err error
	file, err = os.Open(filePath)
	if err != nil {
		return nil, []string{}
	}

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