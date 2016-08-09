package main

import (
	"testing"
	"path"
	"os"
	"strconv"
	"bytes"
	"sort"
	"time"
	"github.com/djherbis/times"
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


func TestLatestFolderInfosLimitsNumberOfResults(t *testing.T) {
	fileInfo := DummyFileInfo{"file1", 1, os.ModeDir, time.Now(), true}
	fileInfos := []os.FileInfo{fileInfo, fileInfo, fileInfo}
	timespec := DummyTimespec{time.Now()}
	timespecs := []DummyTimespec{timespec, timespec, timespec, timespec}
	cannedConversions := CannedConversions{0, timespecs}
	expectInt(t, "length of latest folderInfos", 2, len(LatestFolderInfos(fileInfos, 2, cannedConversions.conversion)))
	cannedConversions = CannedConversions{0, timespecs}
	expectInt(t, "length of latest folderInfos", 3, len(LatestFolderInfos(fileInfos, 3, cannedConversions.conversion)))
	cannedConversions = CannedConversions{0, timespecs}
	expectInt(t, "length of latest folderInfos", 3, len(LatestFolderInfos(fileInfos, 4, cannedConversions.conversion)))
}

func TestLatestFolderInfosOrdersByChangeTime(t *testing.T) {
	fileInfos := []os.FileInfo{
		DummyFileInfo{"file1", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file2", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file3", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file4", 1, os.ModeDir, time.Now(), true},
		DummyFileInfo{"file5", 1, os.ModeDir, time.Now(), true},
	}
	now := time.Now()
	timespecs := []DummyTimespec{
		{now.Add(12 * time.Hour)},
		{now},
		{now.Add(3 * time.Hour)},
		{now.Add(9 * time.Hour)},
		{now.Add(2 * time.Hour)},
	}
	cannedConversions := CannedConversions{0, timespecs}
	expect(t, "folder names", "file1|file4|file3", folderInfoNames(LatestFolderInfos(fileInfos, 3, cannedConversions.conversion)))
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

func folderInfoNames(folderInfos []FolderInfo) string {
	var allnames bytes.Buffer
	for i, folderInfo := range folderInfos {
		allnames.WriteString(folderInfo.fileInfo.Name())
		if i < len(folderInfos) - 1 {
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

type DummyTimespec struct {
	changeTime time.Time
}
func (d DummyTimespec) ModTime() time.Time {
	return time.Now()
}
func (d DummyTimespec) AccessTime() time.Time {
	return time.Now()
}
func (d DummyTimespec) ChangeTime() time.Time {
	return d.changeTime
}
func (d DummyTimespec) BirthTime() time.Time {
	return time.Now()
}
func (d DummyTimespec) HasChangeTime() bool {
	return true
}
func (d DummyTimespec) HasBirthTime() bool {
	return false
}

type CannedConversions struct {
	cursor int
	timespecs []DummyTimespec
}

func (c *CannedConversions) conversion(fi os.FileInfo) times.Timespec {
	c.cursor = c.cursor + 1
	return c.timespecs[c.cursor - 1]
}