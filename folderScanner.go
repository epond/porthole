package main

import (
	"os"
	"path"
	"fmt"
	"log"
)

func LatestAdditions(musicFolder string) string {
	var allLeaves = append(
		FileInfoAtDepth(path.Join(musicFolder, "flac"), 3),
		FileInfoAtDepth(path.Join(musicFolder, "flac-cd"), 3)...)
	return fmt.Sprintf("%v", len(allLeaves))
}

func FileInfoAtDepth(rootFolderPath string, targetDepth int) []os.FileInfo {
	if targetDepth <= 0 {
		return []os.FileInfo{}
	}

	rootFile, err := os.Open(rootFolderPath)
	if err != nil {
		log.Fatalf("Could not open folder. Cause: %v", err)
	}
	leaves, err := rootFile.Readdir(0)
	if err != nil {
		log.Fatalf("Could not read folder info. Cause: %v", err)
	}

	if targetDepth == 1 {
		// TODO exclude apple and synology hidden folders
		return leaves
	}

	for _, _ = range leaves {
		// TODO concatenate results from children at targetDepth - 1
	}

	return leaves
}