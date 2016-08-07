package main

import (
	"testing"
	"path"
	"os"
	"strconv"
	"bytes"
	"sort"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	expect(t, "number of fileinfos", "0", strconv.Itoa(len(FileInfoAtDepth("anything", 0))))
}

func TestGivenDepthOfOneThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := sortByName(FileInfoAtDepth(folderPath, 1))
	expect(t, "number of fileinfos", "3", strconv.Itoa(len(fileInfos)))
	expect(t, "foldernames", "a2|b2|c2", fileInfoNames(fileInfos))
}

func TestGivenDepthOfTwoThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := sortByName(FileInfoAtDepth(folderPath, 2))
	expect(t, "number of fileinfos", "4", strconv.Itoa(len(fileInfos)))
	expect(t, "foldernames", "a3b2|a3c2|b3b2|b3c2", fileInfoNames(fileInfos))
}

func expect(t *testing.T, valueName string, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %v to be %v but was %v", valueName, expected, actual)
	}
}

func fileInfoNames(fileInfos []os.FileInfo) string {
	var allnames bytes.Buffer
	for i, fileInfo := range fileInfos {
		allnames.WriteString(fileInfo.Name())
		if i < len(fileInfos) - 1 {
			allnames.WriteString("|")
		}
	}
	return allnames.String()
}

type FileInfos []os.FileInfo

func (slice FileInfos) Len() int {
	return len(slice)
}

func (slice FileInfos) Less(i, j int) bool {
	return slice[i].Name() < slice[j].Name();
}

func (slice FileInfos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByName(fileInfos FileInfos) FileInfos {
	sort.Sort(fileInfos)
	return fileInfos
}