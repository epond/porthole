package foldermusic

import (
	"github.com/epond/porthole/status"
)

// Additions treats folders on the filesystem as albums
type Additions struct {
	folderScanner FolderScanner
	knownAlbums   KnownAlbums
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
	knownAlbums KnownAlbums) *Additions {
	return &Additions{
		folderScanner,
		knownAlbums,
	}
}

// FetchLatestAdditions finds the most recently added albums
func (f *Additions) FetchLatestAdditions() []status.Album {
	scannedAlbums := f.folderScanner.ScanFolders()
	latestAdditions := f.knownAlbums.UpdateKnownAlbums(scannedAlbums)

	return latestAdditions
}
