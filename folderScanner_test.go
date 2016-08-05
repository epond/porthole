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
	if len(FileInfoAtDepth(folderPath, 1)) != 3 {
		t.Error("File info array at depth 1 did not have length 3")
	}
}