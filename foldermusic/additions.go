package foldermusic

import (
	"log"
	"path"
	"strconv"
	"strings"
)

// Additions treats folders on the filesystem as albums
type Additions struct {
	musicFolder          string
	knownAlbumsFile      string
	knownAlbumsBackup    string
	foldersToScan        string
	latestAdditionsLimit int
	folderScanner        FolderInfoScanner
}

// FolderInfoScanner can scan for folder information
type FolderInfoScanner interface {
	ScanFolders(foldersToScan []FolderToScan) []FolderInfo
}

// NewAdditions constructs a new Additions
func NewAdditions(
	musicFolder string,
	knownAlbumsFile string,
	knownAlbumsBackup string,
	foldersToScan string,
	latestAdditionsLimit int) *Additions {
	return &Additions{
		musicFolder,
		knownAlbumsFile,
		knownAlbumsBackup,
		foldersToScan,
		latestAdditionsLimit,
		&FolderScanner{},
	}
}

// FetchLatestAdditions finds the most recently added albums
func (f *Additions) FetchLatestAdditions() []string {
	foldersToScan := parseFoldersToScan(f.musicFolder, f.foldersToScan)
	scannedAlbums := f.folderScanner.ScanFolders(foldersToScan)
	latestAdditions := UpdateKnownAlbums(
		scannedAlbums,
		f.knownAlbumsFile,
		f.knownAlbumsBackup,
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
