package main

import (
	"os"
	"path"
	"fmt"
)

func LatestAdditions(musicFolder string) string {
	var allLeaves = append(
		FileInfoAtDepth(path.Join(musicFolder, "flac"), 3),
		FileInfoAtDepth(path.Join(musicFolder, "flac-cd"), 3)...)
	return fmt.Sprintf("%v", len(allLeaves))
}

func FileInfoAtDepth(rootFolderPath string, targetDepth int) []os.FileInfo {
	var leaves []os.FileInfo
	return leaves
}