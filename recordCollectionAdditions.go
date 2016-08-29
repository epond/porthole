package main

import (
	"strings"
	"strconv"
	"path"
)

type FileBasedAdditions struct {
	musicFolder string
	knownReleasesFile string
	foldersToScan string
	latestAdditionsLimit int
}

func (f *FileBasedAdditions) latestAdditions() []string {
	return UpdateKnownReleases(
		ScanFolders(parseFoldersToScan(f.musicFolder, f.foldersToScan)),
		f.knownReleasesFile,
		f.latestAdditionsLimit)
}

func parseFoldersToScan(musicFolder string, folders string) []FolderToScan {
	var foldersToScan []FolderToScan
	folderPairs := strings.Split(folders, ",")
	for _, pair := range folderPairs {
		pairArray := strings.Split(pair, ":")
		depth, _ := strconv.Atoi(pairArray[1])
		foldersToScan = append(foldersToScan, FolderToScan{
			path.Join(musicFolder, pairArray[0]),
			depth})
	}
	return foldersToScan
}