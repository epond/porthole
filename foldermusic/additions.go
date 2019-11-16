package foldermusic

import (
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/epond/porthole/status"
)

// Additions treats folders on the filesystem as albums
type Additions struct {
	knownAlbumsFile      string
	knownAlbumsBackup    string
	latestAdditionsLimit int
	knownAlbums          KnownAlbums
	folderScanner        FolderScanner
}

// FolderScanner can scan for folder information
type FolderScanner interface {
	ScanFolders() []status.Album
}

// KnownAlbums provides access to the persisted known albums
type KnownAlbums interface {
	UpdateKnownAlbums(folderScanList []status.Album) []status.Album
}

// NewAdditions constructs a new Additions
func NewAdditions(
	folderScanner FolderScanner,
	knownAlbumsFile string,
	knownAlbumsBackup string,
	latestAdditionsLimit int) *Additions {
	return &Additions{
		knownAlbumsFile,
		knownAlbumsBackup,
		latestAdditionsLimit,
		&KnownAlbumsWithBackup{
			knownAlbumsFile,
			knownAlbumsBackup,
			latestAdditionsLimit,
		},
		folderScanner,
	}
}

// FetchLatestAdditions finds the most recently added albums
func (f *Additions) FetchLatestAdditions() []status.Album {
	scannedAlbums := f.folderScanner.ScanFolders()
	latestAdditions := f.knownAlbums.UpdateKnownAlbums(scannedAlbums)

	return latestAdditions
}

func ParseFoldersToScan(musicFolder string, folders string) []FolderToScan {
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
