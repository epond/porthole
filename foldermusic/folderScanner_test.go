package foldermusic

import (
	"bytes"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/epond/porthole/shared"
	"github.com/epond/porthole/test"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	albums := scanFolders("anything:0")
	test.ExpectInt(t, "number of albums", 0, len(albums))
}

// A folder such as a/b needs to have entry "a - b" where a is artist and b is name of album
func TestScanListEntriesContainTwoFolderLevels(t *testing.T) {
	albums := scanFolders("a1:1")
	test.ExpectInt(t, "number of albums", 3, len(albums))
	test.Expect(t, "album strings", "A1 - A2|A1 - B2|A1 - C2", pipeDelimitedString(albums))
}

func TestScanListAtGreaterDepth(t *testing.T) {
	albums := scanFolders("a1:2")
	test.ExpectInt(t, "number of albums", 4, len(albums))
	test.Expect(t, "folderInfo strings", "B2 - A3b2|B2 - B3b2|C2 - A3c2|C2 - B3c2", pipeDelimitedString(albums))
}

func TestMultipleFoldersToScan(t *testing.T) {
	albums := scanFolders("a1:1,b1:1")
	test.ExpectInt(t, "number of albums", 4, len(albums))
	test.Expect(t, "album strings", "A1 - A2|A1 - B2|A1 - C2|B1 - D2", pipeDelimitedString(albums))
}

func TestFolderInfoStringCapitalisesFirstLettersOnly(t *testing.T) {
	folderInfo := FolderInfo{test.DummyFileInfo{"Ef Gh", true}, test.DummyFileInfo{"Ab Cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"ef Gh", true}, test.DummyFileInfo{"Ab Cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"Ef gh", true}, test.DummyFileInfo{"Ab Cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"Ef Gh", true}, test.DummyFileInfo{"ab Cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"Ef Gh", true}, test.DummyFileInfo{"Ab cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"EF GH", true}, test.DummyFileInfo{"AB CD", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef Gh", folderInfo.String())
}

func TestFolderInfoStringCapitalisesEdgeCases(t *testing.T) {
	folderInfo := FolderInfo{test.DummyFileInfo{"q", true}, test.DummyFileInfo{"m", true}}
	test.Expect(t, "FolderInfo String()", "M - Q", folderInfo.String())
	folderInfo = FolderInfo{test.DummyFileInfo{"Ef  Gh", true}, test.DummyFileInfo{"Ab Cd", true}}
	test.Expect(t, "FolderInfo String()", "Ab Cd - Ef  Gh", folderInfo.String())
}

func pipeDelimitedString(list []shared.Album) string {
	var all bytes.Buffer
	for i, element := range list {
		all.WriteString(element.Text)
		if i < len(list)-1 {
			all.WriteString("|")
		}
	}
	return all.String()
}

type SortedAlbums []shared.Album

func (slice SortedAlbums) Len() int {
	return len(slice)
}

func (slice SortedAlbums) Less(i, j int) bool {
	return slice[i].Text < slice[j].Text
}

func (slice SortedAlbums) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortAlbums(albums SortedAlbums) SortedAlbums {
	sort.Sort(albums)
	return albums
}

func scanFolders(foldersToScan string) []shared.Album {
	fs := &DepthAwareFolderScanner{path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata"), foldersToScan}
	return sortAlbums(fs.ScanFolders())
}
