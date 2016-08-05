package main

import (
	"testing"
	"path"
	"os"
	"fmt"
	"strconv"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	expect(t, "number of fileinfos", "0", strconv.Itoa(len(FileInfoAtDepth("anything", 0))))
}

func TestGivenDepthOfOneThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := FileInfoAtDepth(folderPath, 1)
	expect(t, "number of fileinfos", "3", strconv.Itoa(len(fileInfos)))

	foldernames := fmt.Sprintf("%v%v%v", fileInfos[0].Name(), fileInfos[1].Name(), fileInfos[2].Name())
	expect(t, "foldernames", "a2b2c2", foldernames)
}

func TestGivenDepthOfTwoThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := FileInfoAtDepth(folderPath, 2)
	expect(t, "number of fileinfos", "4", strconv.Itoa(len(fileInfos)))
	foldernames := fmt.Sprintf("%v%v%v%v", fileInfos[0].Name(), fileInfos[1].Name(), fileInfos[2].Name(), fileInfos[3].Name())
	expect(t, "foldernames", "a3b2b3b2a3c2b3c2", foldernames)
}

func expect(t *testing.T, valueName string, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %v to be %v but was %v", valueName, expected, actual)
	}
}
