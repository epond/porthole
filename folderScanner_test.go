package main

import (
	"testing"
	"path"
	"os"
)

func TestGivenZeroDepthThenReturnEmptyArray(t *testing.T) {
	if len(FileInfoAtDepth("anything", 0)) != 0 {
		t.Error("File info array at depth 0 was non-empty")
	}
}

func TestGivenDepthOfOneThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := FileInfoAtDepth(folderPath, 1)
	if len(fileInfos) != 3 {
		t.Error("File info array at depth 1 did not have length 3")
	}
	if fileInfos[0].Name() != "a2" || fileInfos[1].Name() != "b2" || fileInfos[2].Name() != "c2" {
		t.Error("Folder name was not as expected")
	}
}

func TestGivenDepthOfTwoThenReturnSubfolderInfo(t *testing.T) {
	folderPath := path.Join(os.Getenv("GOPATH"), "src/github.com/epond/porthole/testdata/a1")
	fileInfos := FileInfoAtDepth(folderPath, 2)
	if len(fileInfos) != 4 {
		t.Error("File info array at depth 2 did not have length 4")
	}
	if fileInfos[0].Name() != "a3b2" || fileInfos[1].Name() != "b3b2" || fileInfos[2].Name() != "a3c2" || fileInfos[3].Name() != "b3c2" {
		t.Error("Folder name was not as expected")
	}
}
