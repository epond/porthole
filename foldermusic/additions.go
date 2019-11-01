package foldermusic

import (
	"log"
	"path"
	"strconv"
	"strings"
)

// Additions treats folders on the filesystem as releases
type Additions struct {
	musicFolder          string
	knownReleasesFile    string
	knownReleasesBackup  string
	foldersToScan        string
	latestAdditionsLimit int
}

// NewAdditions constructs a new Additions
func NewAdditions(
	musicFolder string,
	knownReleasesFile string,
	knownReleasesBackup string,
	foldersToScan string,
	latestAdditionsLimit int) *Additions {
	return &Additions{musicFolder, knownReleasesFile, knownReleasesBackup, foldersToScan, latestAdditionsLimit}
}

// FetchLatestAdditions finds latest releases
func (f *Additions) FetchLatestAdditions() []string {
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
