package music

import (
	"os"
	"path"
	"fmt"
	"log"
	"unicode/utf8"
	"strings"
)

type FolderToScan struct {
	rootFolderPath string
	targetDepth int
}

type FolderInfo struct {
	fileInfo os.FileInfo
	parent os.FileInfo
}

func ScanFolders(foldersToScan []FolderToScan) []FolderInfo {
	var folderScanList []FolderInfo
	for _, folder := range foldersToScan {
		folderScanList = append(folderScanList, FolderInfoAtDepth(folder)...)
	}
	return folderScanList
}

func FolderInfoAtDepthIter(folderToScan FolderToScan, parent os.FileInfo) []FolderInfo {
	if folderToScan.targetDepth <= 0 {
		return []FolderInfo{}
	}

	rootFile, err := os.Open(folderToScan.rootFolderPath)
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
			if folderToScan.targetDepth == 1 {
				folderInfos = append(folderInfos, FolderInfo{child, parent})
			} else {
				nextFolder := FolderToScan{
					path.Join(folderToScan.rootFolderPath, child.Name()),
					folderToScan.targetDepth - 1,
				}
				childFolderInfos := FolderInfoAtDepthIter(nextFolder, child)
				folderInfos = append(folderInfos, childFolderInfos...)
			}
		}
	}

	return folderInfos
}

func FolderInfoAtDepth(folderToScan FolderToScan) []FolderInfo {
	return FolderInfoAtDepthIter(folderToScan, nil)
}

func (f *FolderInfo) String() string {
	if f.parent == nil {
		return f.fileInfo.Name()
	}
	return fmt.Sprintf("%v - %v", capitalise(f.parent.Name()), capitalise(f.fileInfo.Name()))
}

func capitalise(word string) string {
	inputSplit := strings.Split(word, " ")
	for i, s := range inputSplit {
		var firstLetter string
		var remainder string
		if len(s) > 0 {
			firstLetter = strings.ToUpper(s[0:1])
			remainder = strings.ToLower(s[1:])
		} else {
			firstLetter = ""
			remainder = ""
		}
		inputSplit[i] = fmt.Sprintf("%v%v", firstLetter, remainder)
	}
	return strings.Join(inputSplit, " ")
}