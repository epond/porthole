package foldermusic

import (
	"bufio"
	"os"
	"path"
	"testing"

	"github.com/epond/porthole/test"
)

func TestItCreatesFileWhenKnownAlbumsFileMissing(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownAlbums([]FolderInfo{
		{test.DummyFileInfo{"Mouldy Old Dough", true}, test.DummyFileInfo{"Lieutenant Pigeon", true}},
	}, knownAlbumsFile(), knownAlbumsBackup(), 0)
	_, lines := knownAlbumsLines()
	if len(lines) == 0 {
		t.Error("Expected UpdateKnownAlbums to create known albums file but it didn't")
	}

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "known albums on creation", 1, len(lines))
	test.Expect(t, "known album", "Lieutenant Pigeon - Mouldy Old Dough", lines[len(lines)-1])
}

func TestItDoesNotChangeFileWhenNoNewAlbums(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, test.DummyFileInfo{"Abba", true}},
		{test.DummyFileInfo{"It's Fan-dabi-dozi!", true}, test.DummyFileInfo{"The Krankies", true}},
		{test.DummyFileInfo{"Discipline", true}, test.DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "number of known albums", 3, len(lines))
	test.Expect(t, "known album 1", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-1])
	test.Expect(t, "known album 2", "Throbbing Gristle - Discipline", lines[len(lines)-2])
	test.Expect(t, "known album 3", "Abba - I Do, I Do, I Do, I Do, I Do", lines[len(lines)-3])
}

func TestItAddsAlbumsToEndOfFile(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}
	UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "number of known albums", 5, len(lines))
	test.Expect(t, "known album 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known album 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
	test.Expect(t, "known album 3", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-3])
}

func TestItIgnoresAlbumsAlreadyKnown(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
		{test.DummyFileInfo{"I Do, I Do, I Do, I Do, I Do", true}, test.DummyFileInfo{"Abba", true}},
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Discipline", true}, test.DummyFileInfo{"Throbbing Gristle", true}},
	}
	UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "number of known albums", 5, len(lines))
	test.Expect(t, "known album 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known album 2", "Shake - Iconoclastic Diaries", lines[len(lines)-2])
}

func TestItHandlesWhenKnownAlbumsFileMayNotEndInNewline(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "knownalbums_endwithoutnewline"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "number of known albums", 2, len(lines))
	test.Expect(t, "known album 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known album 2", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-2])
}

func TestUpdateKnownAlbumsReturnValueWhenNoNewAlbumsAndKnownAlbumsAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	latestAdditions := UpdateKnownAlbums([]FolderInfo{}, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[0])
	test.Expect(t, "latest addition 2", "Throbbing Gristle - Discipline", latestAdditions[1])
}

func TestUpdateKnownAlbumsReturnValueWhenNewAlbumsAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Mouldy Old Dough", true}, test.DummyFileInfo{"Lieutenant Pigeon", true}},
	}
	latestAdditions := UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "Lieutenant Pigeon - Mouldy Old Dough", latestAdditions[1])
}

func TestUpdateKnownAlbumsReturnValueWhenNewAlbumsBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
}

func TestUpdateKnownAlbumsReturnValueWhenNewAndKnownAlbumsCombinedAreBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
	}
	latestAdditions := UpdateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 5)

	test.ExpectInt(t, "number of latest additions", 4, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0])
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1])
	test.Expect(t, "latest addition 3", "Throbbing Gristle - Discipline", latestAdditions[2])
	test.Expect(t, "latest addition 4", "Abba - I Do, I Do, I Do, I Do, I Do", latestAdditions[3])
}

func TestItBacksUpKnownAlbumsWhenChanged(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	UpdateKnownAlbums([]FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}, knownAlbumsFile(), knownAlbumsBackup(), 0)

	backupFile, backupLines := knownAlbumsBackupLines()
	defer backupFile.Close()
	if len(backupLines) == 0 {
		t.Error("Expected UpdateKnownAlbums to create backup but it didn't")
	}

	test.ExpectInt(t, "number of lines in backup", 5, len(backupLines))
	test.Expect(t, "backup album 1", "Daniel Menche - Vent", backupLines[len(backupLines)-1])
	test.Expect(t, "backup album 2", "Shake - Iconoclastic Diaries", backupLines[len(backupLines)-2])
	test.Expect(t, "backup album 3", "The Krankies - It's Fan-dabi-dozi!", backupLines[len(backupLines)-3])
}

func TestItDoesntBackUpKnownAlbumsWhenUnchanged(t *testing.T) {
	setUp()
	defer tearDown()
	UpdateKnownAlbums([]FolderInfo{}, knownAlbumsFile(), knownAlbumsBackup(), 0)

	_, backup := knownAlbumsBackupLines()
	if len(backup) > 0 {
		t.Error("Expected UpdateKnownAlbums to not create backup but it did")
	}
}

func TestNoMissingAlbums(t *testing.T) {
	scanned := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}
	known := []string{
		"Daniel Menche - Vent",
		"Shake - Iconoclastic Diaries",
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 0, len(missing))
}

func TestOneMissingAlbum(t *testing.T) {
	scanned := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}
	known := []string{
		"Daniel Menche - Vent",
		"The Krankies - It's Fan-dabi-dozi!",
		"Shake - Iconoclastic Diaries",
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 1, len(missing))
	test.Expect(t, "missing album", "The Krankies - It's Fan-dabi-dozi!", missing[0])
}

func TestTwoMissingAlbums(t *testing.T) {
	scanned := []FolderInfo{
		{test.DummyFileInfo{"Vent", true}, test.DummyFileInfo{"Daniel Menche", true}},
		{test.DummyFileInfo{"Iconoclastic Diaries", true}, test.DummyFileInfo{"Shake", true}},
	}
	known := []string{
		"Daniel Menche - Vent",
		"The Krankies - It's Fan-dabi-dozi!",
		"Shake - Iconoclastic Diaries",
		"Throbbing Gristle - Discipline",
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 2, len(missing))
	test.Expect(t, "missing album 1", "The Krankies - It's Fan-dabi-dozi!", missing[0])
	test.Expect(t, "missing album 2", "Throbbing Gristle - Discipline", missing[1])
}

func setUp() {
	os.Mkdir(tempDir(), os.ModePerm)
}
func tearDown() {
	os.RemoveAll(tempDir())
}

func knownAlbumsFile() string {
	return path.Join(tempDir(), "knownalbums")
}

func knownAlbumsBackup() string {
	return path.Join(tempDir(), "knownalbums_backup")
}

func knownAlbumsLines() (file *os.File, lines []string) {
	return fileLines(knownAlbumsFile())
}

func knownAlbumsBackupLines() (file *os.File, lines []string) {
	return fileLines(knownAlbumsBackup())
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
