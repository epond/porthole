package music

import (
	"path"
	"strconv"
	"strings"
)

type FileBasedAdditions struct {
	musicFolder          string
	knownReleasesFile    string
	knownReleasesBackup  string
	foldersToScan        string
	latestAdditionsLimit int
}

func NewFileBasedAdditions(
	musicFolder string,
	knownReleasesFile string,
	knownReleasesBackup string,
	foldersToScan string,
	latestAdditionsLimit int) *FileBasedAdditions {
	return &FileBasedAdditions{musicFolder, knownReleasesFile, knownReleasesBackup, foldersToScan, latestAdditionsLimit}
}

func (f *FileBasedAdditions) FetchLatestAdditions() []string {
	foldersToScan := parseFoldersToScan(f.musicFolder, f.foldersToScan)
	scannedReleases := ScanFolders(foldersToScan)
	latestAdditions := UpdateKnownReleases(
		scannedReleases,
		f.knownReleasesFile,
		f.knownReleasesBackup,
		f.latestAdditionsLimit)

	return latestAdditions
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
