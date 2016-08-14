package main

import (
	"os"
	"path"
	"fmt"
	"log"
	"unicode/utf8"
)

func FolderInfoAtDepthIter(rootFolderPath string, targetDepth int, parent os.FileInfo) []FolderInfo {
	if targetDepth <= 0 {
		return []FolderInfo{}
	}

	rootFile, err := os.Open(rootFolderPath)
	if err != nil {
		log.Fatalf("Could not open root folder. Cause: %v", err)
	}
	children, err := rootFile.Readdir(0)
	if err != nil {
		log.Fatalf("Could not read root folder info. Cause: %v", err)
	}
	if parent == nil {
		rootFileInfo, err := rootFile.Stat()
		if err != nil {
			log.Fatalf("Could not read Stat on root folder. Cause: %v", err)
		}
		parent = rootFileInfo
	}

	folderInfos := make([]FolderInfo, 0)

	for _, child := range children {
		firstChar, _ := utf8.DecodeRuneInString(child.Name())
		if child.IsDir() && child.Name() != "@eaDir" && firstChar != '.' {
			if targetDepth == 1 {
				folderInfos = append(folderInfos, FolderInfo{child, parent})
			} else {
				childFolderInfos := FolderInfoAtDepthIter(path.Join(rootFolderPath, child.Name()), targetDepth - 1, child)
				folderInfos = append(folderInfos, childFolderInfos...)
			}
		}
	}

	return folderInfos
}

func FolderInfoAtDepth(rootFolderPath string, targetDepth int) []FolderInfo {
	return FolderInfoAtDepthIter(rootFolderPath, targetDepth, nil)
}

type FolderInfo struct {
	fileInfo os.FileInfo
	parent os.FileInfo
}

func (f *FolderInfo) String() string {
	if f.parent == nil {
		return f.fileInfo.Name()
	}
	return fmt.Sprintf("%v - %v", f.parent.Name(), f.fileInfo.Name())
}