package main

import "testing"

func TestGivenZeroDepthThenFileInfoAtDepthReturnsEmptyArray(t *testing.T) {
	if len(FileInfoAtDepth("anything", 0)) != 0 {
		t.Error("not implemented")
	}
}