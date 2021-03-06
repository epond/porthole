package foldermusic

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/epond/porthole/shared"
)

// FolderToScan represents a folder in the filesystem
type FolderToScan struct {
	rootFolderPath string
	targetDepth    int
}

// FolderInfo gives information about a folder
type FolderInfo struct {
	fileInfo os.FileInfo
	parent   os.FileInfo
}

// DepthAwareFolderScanner knows how to use the filesystem to scan for folders
type DepthAwareFolderScanner struct {
	MusicFolder string
	Folders     string
}

func (f *FolderInfo) String() string {
	if f.parent == nil {
		return f.fileInfo.Name()
	}
	return fmt.Sprintf("%v - %v", capitalise(f.parent.Name()), capitalise(f.fileInfo.Name()))
}

// ScanFolders scans the filesystem for folders
func (f *DepthAwareFolderScanner) ScanFolders() []shared.Album {
	foldersToScan := parseFoldersToScan(f.MusicFolder, f.Folders)
	var folderScanList []FolderInfo
	for _, folder := range foldersToScan {
		folderScanList = append(folderScanList, folderInfoAtDepth(folder)...)
	}
	albums := make([]shared.Album, 0)
	for _, folderInfo := range folderScanList {
		albums = append(albums, shared.Album{folderInfo.String()})
	}
	return albums
}

func parseFoldersToScan(musicFolder string, folders string) []FolderToScan {
	var foldersToScan []FolderToScan
	folderPairs := strings.Split(folders, ",")
	for _, pair := range folderPairs {
		pairArray := strings.Split(pair, ":")
		if len(pairArray) < 2 {
			log.Fatalf("Could not read depth of folder to scan from configuration. Expected folder:depth but got %v", pair)
		}
		depth, _ := strconv.Atoi(pairArray[1])
		foldersToScan = append(foldersToScan, FolderToScan{
			path.Join(musicFolder, pairArray[0]),
			depth})
	}
	return foldersToScan
}

func folderInfoAtDepthIter(folderToScan FolderToScan, parent os.FileInfo) []FolderInfo {
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
				childFolderInfos := folderInfoAtDepthIter(nextFolder, child)
				folderInfos = append(folderInfos, childFolderInfos...)
			}
		}
	}

	return folderInfos
}

func folderInfoAtDepth(folderToScan FolderToScan) []FolderInfo {
	return folderInfoAtDepthIter(folderToScan, nil)
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
