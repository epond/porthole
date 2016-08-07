package main

import (
	"testing"
	"path"
	"os"
	"strconv"
	"bytes"
	"sort"
	"time"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	expectInt(t, "number of fileinfos", 0, len(FileInfoAtDepth("anything", 0)))
}

func TestGivenDepthOfOneThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := sortByName(FileInfoAtDepth(folderPath, 1))
	expectInt(t, "number of fileinfos", 3, len(fileInfos))
	expect(t, "foldernames", "a2|b2|c2", fileInfoNames(fileInfos))
}

func TestGivenDepthOfTwoThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := sortByName(FileInfoAtDepth(folderPath, 2))
	expectInt(t, "number of fileinfos", 4, len(fileInfos))
	expect(t, "folder names", "a3b2|a3c2|b3b2|b3c2", fileInfoNames(fileInfos))
}

func TestLatestFileInfosLimitsNumberOfResults(t *testing.T) {
	fileInfos := []os.FileInfo{
		DummyFileInfo{"file1", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file2", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file3", 1, os.ModeDir, time.Now(), true},
	}
	expectInt(t, "length of latest fileInfos", 2, len(LatestFileInfos(fileInfos, 2)))
	expectInt(t, "length of latest fileInfos", 3, len(LatestFileInfos(fileInfos, 3)))
	expectInt(t, "length of latest fileInfos", 3, len(LatestFileInfos(fileInfos, 4)))
}

func TestLatestFileInfosOrdersByModifiedTime(t *testing.T) {
	now := time.Now()
	rawFileInfos := []os.FileInfo{
		DummyFileInfo{"file1", 1, os.ModeDir, now.Add(12 * time.Hour), true},
		DummyFileInfo{"file2", 1, os.ModeDir, now, true},
		DummyFileInfo{"file3", 1, os.ModeDir, now.Add(3 * time.Hour), true},
		DummyFileInfo{"file4", 1, os.ModeDir, now.Add(9 * time.Hour), true},
		DummyFileInfo{"file5", 1, os.ModeDir, now.Add(2 * time.Hour), true},
	}
	expect(t, "folder names", "file1|file4|file3", fileInfoNames(LatestFileInfos(rawFileInfos, 3)))
}

func expect(t *testing.T, valueName string, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %v to be %v but was %v", valueName, expected, actual)
	}
}

func expectInt(t *testing.T, valueName string, expected int, actual int) {
	expect(t, valueName, strconv.Itoa(expected), strconv.Itoa(actual))
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

type FileInfosSortedByName []os.FileInfo

func (slice FileInfosSortedByName) Len() int {
	return len(slice)
}

func (slice FileInfosSortedByName) Less(i, j int) bool {
	return slice[i].Name() < slice[j].Name();
}

func (slice FileInfosSortedByName) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByName(fileInfos FileInfosSortedByName) FileInfosSortedByName {
	sort.Sort(fileInfos)
	return fileInfos
}

type DummyFileInfo struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
	isDir bool
}

func (d DummyFileInfo) Name() string {
	return d.name
}

func (d DummyFileInfo) Size() int64 {
	return d.size
}

func (d DummyFileInfo) Mode() os.FileMode {
	return d.mode
}

func (d DummyFileInfo) ModTime() time.Time {
	return d.modTime
}

func (d DummyFileInfo) IsDir() bool {
	return d.isDir
}

func (d DummyFileInfo) Sys() interface{} {
	return nil
}