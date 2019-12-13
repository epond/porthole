package foldermusic

import (
	"bufio"
	"os"
	"path"
	"testing"

	"github.com/epond/porthole/shared"
	"github.com/epond/porthole/test"
)

func TestItCreatesFileWhenKnownAlbumsFileMissing(t *testing.T) {
	setUp()
	defer tearDown()
	updateKnownAlbums([]shared.Album{shared.Album{"Lieutenant Pigeon - Mouldy Old Dough"}}, knownAlbumsFile(), knownAlbumsBackup(), 0)
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
	currentScan := []shared.Album{
		shared.Album{"Abba - I Do, I Do, I Do, I Do, I Do"},
		shared.Album{"The Krankies - It's Fan-dabi-dozi!"},
		shared.Album{"Throbbing Gristle - Discipline"},
	}
	updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

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
	currentScan := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

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
	currentScan := []shared.Album{
		shared.Album{"Shake - Iconoclastic Diaries"},
		shared.Album{"Abba - I Do, I Do, I Do, I Do, I Do"},
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Throbbing Gristle - Discipline"},
	}
	updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

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
	currentScan := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
	}
	updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 0)

	file, lines := knownAlbumsLines()
	defer file.Close()

	test.ExpectInt(t, "number of known albums", 2, len(lines))
	test.Expect(t, "known album 1", "Daniel Menche - Vent", lines[len(lines)-1])
	test.Expect(t, "known album 2", "The Krankies - It's Fan-dabi-dozi!", lines[len(lines)-2])
}

func TestReturnValueWhenNoNewAlbumsAndKnownAlbumsAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	latestAdditions := updateKnownAlbums([]shared.Album{}, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[0].Text)
	test.Expect(t, "latest addition 2", "Throbbing Gristle - Discipline", latestAdditions[1].Text)
}

func TestReturnValueWhenNewAlbumsAboveLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []shared.Album{
		shared.Album{"Shake - Iconoclastic Diaries"},
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Lieutenant Pigeon - Mouldy Old Dough"},
	}
	latestAdditions := updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0].Text)
	test.Expect(t, "latest addition 2", "Lieutenant Pigeon - Mouldy Old Dough", latestAdditions[1].Text)
}

func TestReturnValueWhenNewAlbumsBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
	}
	latestAdditions := updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 2)

	test.ExpectInt(t, "number of latest additions", 2, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0].Text)
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1].Text)
}

func TestReturnValueWhenNewAndKnownAlbumsCombinedAreBelowLimit(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	currentScan := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
	}
	latestAdditions := updateKnownAlbums(currentScan, knownAlbumsFile(), knownAlbumsBackup(), 5)

	test.ExpectInt(t, "number of latest additions", 4, len(latestAdditions))
	test.Expect(t, "latest addition 1", "Daniel Menche - Vent", latestAdditions[0].Text)
	test.Expect(t, "latest addition 2", "The Krankies - It's Fan-dabi-dozi!", latestAdditions[1].Text)
	test.Expect(t, "latest addition 3", "Throbbing Gristle - Discipline", latestAdditions[2].Text)
	test.Expect(t, "latest addition 4", "Abba - I Do, I Do, I Do, I Do, I Do", latestAdditions[3].Text)
}

func TestItBacksUpKnownAlbumsWhenChanged(t *testing.T) {
	setUp()
	defer tearDown()
	test.CopyFile(knownAlbumsFile(), path.Join(testData(), "3knownalbums"))
	updateKnownAlbums([]shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
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
	updateKnownAlbums([]shared.Album{}, knownAlbumsFile(), knownAlbumsBackup(), 0)

	_, backup := knownAlbumsBackupLines()
	if len(backup) > 0 {
		t.Error("Expected UpdateKnownAlbums to not create backup but it did")
	}
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

func updateKnownAlbums(folderScanList []shared.Album, knownAlbumsPath string, knownAlbumsBackupPath string, limit int) []shared.Album {
	folderScanner := &DummyFolderScanner{folderScanList}
	knownAlbums := &KnownAlbumsWithBackup{knownAlbumsPath, knownAlbumsBackupPath}
	additions := NewAdditions(folderScanner, knownAlbums, limit)
	return additions.FetchLatestAdditions()
}

type DummyFolderScanner struct {
	folderScanList []shared.Album
}

func (d *DummyFolderScanner) ScanFolders() []shared.Album {
	return d.folderScanList
}
