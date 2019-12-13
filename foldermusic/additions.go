package foldermusic

import (
	"log"

	"github.com/epond/porthole/shared"
)

// Additions treats folders on the filesystem as albums
type Additions struct {
	scanner     shared.Scanning
	persistence shared.Persistence
	limit       int
}

// NewAdditions constructs a new Additions
func NewAdditions(
	scanner shared.Scanning,
	persistence shared.Persistence,
	limit int) *Additions {
	return &Additions{
		scanner,
		persistence,
		limit,
	}
}

// FetchLatestAdditions finds the most recently added albums
func (a *Additions) FetchLatestAdditions() []shared.Album {
	scannedAlbums := a.scanner.ScanFolders()
	// Read known albums file into an array of its lines and a map that conveys if a line is present
	knownAlbums, knownAlbumsMap := a.persistence.ReadKnownAlbums()

	// Build a list of current scan entries not present in known albums (new albums)
	var newAlbums []shared.Album
	for _, scanItem := range scannedAlbums {
		if knownAlbumsMap[scanItem.Text] != present {
			newAlbums = append(newAlbums, scanItem)
		}
	}

	reverseSortByName(newAlbums)
	log.Printf("Found %v known and %v new albums", len(knownAlbumsMap), len(newAlbums))

	missingAlbums := findMissingAlbums(scannedAlbums, knownAlbums)

	if len(missingAlbums) > 0 {
		log.Printf("Found %v missing albums", len(missingAlbums))
		for i, missing := range missingAlbums {
			log.Printf("Missing #%v: %v", i+1, missing)
		}
	}

	// Append new albums to known albums file
	a.persistence.AppendNewAlbums(knownAlbums, newAlbums)

	// Return sorted new albums then knownalbums from the end, up to a total of limit
	sortByName(newAlbums)
	var latestAdditions []shared.Album
	i := 0
	for i < min(len(newAlbums), a.limit) {
		latestAdditions = append(latestAdditions, newAlbums[i])
		i++
	}

	i = 0
	for i < (a.limit - len(newAlbums)) {
		if i < len(knownAlbums) {
			latestAdditions = append(latestAdditions, knownAlbums[len(knownAlbums)-i-1])
		}
		i++
	}

	return latestAdditions
}

func findMissingAlbums(scanned []shared.Album, known []shared.Album) []shared.Album {
	scannedMap := make(map[shared.Album]int)
	for _, album := range scanned {
		scannedMap[album] = present
	}

	// Build a list of known albums not present in current scan
	var missingList []shared.Album
	for _, album := range known {
		if scannedMap[album] != present {
			missingList = append(missingList, album)
		}
	}

	return missingList
}
