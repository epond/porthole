package main

import (
	"testing"
	"path"
	"os"
	"bytes"
	"sort"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	expectInt(t, "number of folderInfos", 0, len(FolderInfoAtDepth(FolderToScan{"anything", 0})))
}

// A folder such as a/b needs to have entry "a - b" where a is artist and b is name of release
func TestScanListEntriesContainTwoFolderLevels(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	folderInfos := sortFolderInfoByString(ScanFolders([]FolderToScan{{folderPath, 1}}))
	expectInt(t, "number of folderInfos", 3, len(folderInfos))
	expect(t, "folderInfo strings", "a1 - a2|a1 - b2|a1 - c2", pipeDelimitedString(folderInfos))
}

func TestScanListAtGreaterDepth(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	folderInfos := sortFolderInfoByString(ScanFolders([]FolderToScan{{folderPath, 2}}))
	expectInt(t, "number of folderInfos", 4, len(folderInfos))
	expect(t, "folderInfo strings", "b2 - a3b2|b2 - b3b2|c2 - a3c2|c2 - b3c2", pipeDelimitedString(folderInfos))
}

func pipeDelimitedString(list []FolderInfo) string {
	var all bytes.Buffer
	for i, element := range list {
		all.WriteString(element.String())
		if i < len(list) - 1 {
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
	return slice[i].String() < slice[j].String();
}

func (slice FolderInfosSortedByString) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortFolderInfoByString(folderInfos FolderInfosSortedByString) FolderInfosSortedByString {
	sort.Sort(folderInfos)
	return folderInfos
}