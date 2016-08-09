package main

import (
	"os"
	"path"
	"fmt"
	"log"
	"unicode/utf8"
	"sort"
	"github.com/djherbis/times"
)

func LatestAdditions(musicFolder string) string {
	allFileInfos := append(
		FileInfoAtDepth(path.Join(musicFolder, "flac-add"), 2),
		FileInfoAtDepth(path.Join(musicFolder, "flac-vorbis320"), 2)...)
	latestFolderInfos := LatestFolderInfos(allFileInfos, 3, times.Get)
	return fmt.Sprintf("%v, %v, %v",
		latestFolderInfos[0].fileInfo.Name(),
		latestFolderInfos[1].fileInfo.Name(),
		latestFolderInfos[2].fileInfo.Name())
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

func LatestFolderInfos(fileInfos []os.FileInfo, limit int, conversion func(fi os.FileInfo) times.Timespec) []FolderInfo {
	folderInfos := BuildFolderInfos(fileInfos, conversion)
	sort.Sort(folderInfos)
	if len(fileInfos) > limit {
		return folderInfos[:limit]
	}
	return folderInfos
}

type FolderInfo struct {
	fileInfo os.FileInfo
	timespec times.Timespec
}

type FolderInfos []FolderInfo

func (slice FolderInfos) Len() int {
	return len(slice)
}

func (slice FolderInfos) Less(i, j int) bool {
	return slice[i].timespec.BirthTime().After(slice[j].timespec.BirthTime())
}

func (slice FolderInfos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func BuildFolderInfos(fileInfos []os.FileInfo, conversion func(fi os.FileInfo) times.Timespec) FolderInfos {
	folderInfos := make([]FolderInfo, len(fileInfos))
	for i, fileInfo := range fileInfos {
		folderInfos[i] = FolderInfo{fileInfo, conversion(fileInfo)}
	}
	return folderInfos
}