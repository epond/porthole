package main

import (
	"os"
	"path"
	"fmt"
	"log"
	"unicode/utf8"
	"sort"
)

func LatestAdditions(musicFolder string) string {
	allFileInfos := append(
		FileInfoAtDepth(path.Join(musicFolder, "flac-add"), 2),
		FileInfoAtDepth(path.Join(musicFolder, "flac-vorbis320"), 2)...)
	latestFileInfos := LatestFileInfos(allFileInfos, 3)
	return fmt.Sprintf("%v, %v, %v", latestFileInfos[0].Name(), latestFileInfos[1].Name(), latestFileInfos[2].Name())
}

func FileInfoAtDepth(rootFolderPath string, targetDepth int) []os.FileInfo {
	if targetDepth <= 0 {
		return []os.FileInfo{}
	}

	rootFile, err := os.Open(rootFolderPath)
	if err != nil {
		log.Fatalf("Could not open root folder. Cause: %v", err)
	}
	children, err := rootFile.Readdir(0)
	if err != nil {
		log.Fatalf("Could not read root folder info. Cause: %v", err)
	}

	fileInfos := make([]os.FileInfo, 0)

	for _, child := range children {
		firstChar, _ := utf8.DecodeRuneInString(child.Name())
		if child.IsDir() && child.Name() != "@eaDir" && firstChar != '.' {
			if targetDepth == 1 {
				fileInfos = append(fileInfos, child)
			} else {
				childFileInfos := FileInfoAtDepth(path.Join(rootFolderPath, child.Name()), targetDepth - 1)
				fileInfos = append(fileInfos, childFileInfos...)
			}
		}
	}

	return fileInfos
}

func LatestFileInfos(fileInfos []os.FileInfo, limit int) []os.FileInfo {
	sortedFileInfos := sortByModTime(fileInfos)
	if len(sortedFileInfos) > limit {
		return sortedFileInfos[:limit]
	}
	return sortedFileInfos
}

type FileInfosSortedByModifiedTime []os.FileInfo

func (slice FileInfosSortedByModifiedTime) Len() int {
	return len(slice)
}

func (slice FileInfosSortedByModifiedTime) Less(i, j int) bool {
	return slice[i].ModTime().After(slice[j].ModTime())
}

func (slice FileInfosSortedByModifiedTime) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByModTime(fileInfos FileInfosSortedByModifiedTime) FileInfosSortedByModifiedTime {
	sort.Sort(fileInfos)
	return fileInfos
}
