package foldermusic

import (
	"bytes"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/epond/porthole/test"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	test.ExpectInt(t, "number of folderInfos", 0, len(folderInfoAtDepth(FolderToScan{"anything", 0})))
}

// A folder such as a/b needs to have entry "a - b" where a is artist and b is name of album
func TestScanListEntriesContainTwoFolderLevels(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	folderInfos := sortFolderInfoByString(ScanFolders([]FolderToScan{{folderPath, 1}}))
	test.ExpectInt(t, "number of folderInfos", 3, len(folderInfos))
	test.Expect(t, "folderInfo strings", "A1 - A2|A1 - B2|A1 - C2", pipeDelimitedString(folderInfos))
}

func TestScanListAtGreaterDepth(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	folderInfos := sortFolderInfoByString(ScanFolders([]FolderToScan{{folderPath, 2}}))
	test.ExpectInt(t, "number of folderInfos", 4, len(folderInfos))
	test.Expect(t, "folderInfo strings", "B2 - A3b2|B2 - B3b2|C2 - A3c2|C2 - B3c2", pipeDelimitedString(folderInfos))
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

func pipeDelimitedString(list []FolderInfo) string {
	var all bytes.Buffer
	for i, element := range list {
		all.WriteString(element.String())
		if i < len(list)-1 {
			all.WriteString("|")
		}
	}
	return all.String()
}

type FolderInfosSortedByString []FolderInfo

func (slice FolderInfosSortedByString) Len() int {
	return len(slice)
}

func (slice FolderInfosSortedByString) Less(i, j int) bool {
	return slice[i].String() < slice[j].String()
}

func (slice FolderInfosSortedByString) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortFolderInfoByString(folderInfos FolderInfosSortedByString) FolderInfosSortedByString {
	sort.Sort(folderInfos)
	return folderInfos
}
