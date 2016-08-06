package main

import (
	"os"
	"path"
	"fmt"
	"log"
	"unicode/utf8"
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
	children, err := rootFile.Readdir(0)
	if err != nil {
		log.Fatalf("Could not read folder info. Cause: %v", err)
	}

	if targetDepth == 1 {
		fileInfos := make([]os.FileInfo, 0)
		for _, child := range children {
			firstChar, _ := utf8.DecodeRuneInString(child.Name())
			if child.IsDir() && child.Name() != "@eaDir" && firstChar != '.' {
				fileInfos = append(fileInfos, child)
			}
		}
		return fileInfos
	}

	for _, _ = range children {
		// TODO concatenate results from children at targetDepth - 1
	}

	return children
}