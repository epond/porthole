package foldermusic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/epond/porthole/status"
)

const present = 1

type sortableAlbums []status.Album

// KnownAlbumsWithBackup makes a backup each time the known albums are updated
type KnownAlbumsWithBackup struct {
	knownAlbumsPath       string
	knownAlbumsBackupPath string
	limit                 int
}

func (k *KnownAlbumsWithBackup) readKnownAlbums() (albums []status.Album, lineMap map[string]int) {
	if _, err := os.Stat(k.knownAlbumsPath); os.IsNotExist(err) {
		file, errCreate := os.Create(k.knownAlbumsPath)
		if errCreate != nil {
			log.Printf("Could not create known albums at %v", k.knownAlbumsPath)
			panic(errCreate)
		}
		file.Close()
	}

	// Read knownalbums into an array of its lines and a map that conveys if a line is present
	albums, lineMap = readFile(k.knownAlbumsPath)
	return albums, lineMap
}

func (k *KnownAlbumsWithBackup) appendNewAlbums(knownAlbums []status.Album, newAlbums []status.Album) {
	ensureFileEndsInNewline(k.knownAlbumsPath)
	var knownAlbumsFile *os.File
	knownAlbumsFile, err := os.OpenFile(k.knownAlbumsPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("Could not open known albums file for appending: %v", k.knownAlbumsPath)
		panic(err)
	}
	defer knownAlbumsFile.Close()
	kaWriter := bufio.NewWriter(knownAlbumsFile)
	for _, newAlbum := range newAlbums {
		if _, err := kaWriter.WriteString(fmt.Sprintf("%v\n", newAlbum.Text)); err != nil {
			log.Printf("Could not write new album to %v", k.knownAlbumsPath)
			panic(err)
		}
	}
	if err := kaWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", k.knownAlbumsPath)
		panic(err)
	}

	if len(newAlbums) > 0 {
		backupKnownAlbums(k.knownAlbumsBackupPath, append(knownAlbums, newAlbums...))
	}
}

// UpdateKnownAlbums updates the known albums file and returns an array
// of new albums, based upon the array of folder passed in as the
// folderScanList argument.
func (k *KnownAlbumsWithBackup) UpdateKnownAlbums(folderScanList []status.Album) []status.Album {
	// Read knownalbums into an array of its lines and a map that conveys if a line is present
	knownAlbums, knownAlbumsMap := k.readKnownAlbums()

	// Build a list of current scan entries not present in known albums (new albums)
	var newAlbums []status.Album
	for _, scanItem := range folderScanList {
		if knownAlbumsMap[scanItem.Text] != present {
			newAlbums = append(newAlbums, scanItem)
		}
	}

	reverseSortByName(newAlbums)
	log.Printf("Found %v known and %v new albums", len(knownAlbumsMap), len(newAlbums))

	missingAlbums := findMissingAlbums(folderScanList, knownAlbums)

	if len(missingAlbums) > 0 {
		log.Printf("Found %v missing albums", len(missingAlbums))
		for i, missing := range missingAlbums {
			log.Printf("Missing #%v: %v", i+1, missing)
		}
	}

	// Append new albums to known albums file
	k.appendNewAlbums(knownAlbums, newAlbums)

	// Return sorted new albums then knownalbums from the end, up to a total of limit
	sortByName(newAlbums)
	var latestAdditions []status.Album
	i := 0
	for i < min(len(newAlbums), k.limit) {
		latestAdditions = append(latestAdditions, newAlbums[i])
		i++
	}

	i = 0
	for i < (k.limit - len(newAlbums)) {
		if i < len(knownAlbums) {
			latestAdditions = append(latestAdditions, knownAlbums[len(knownAlbums)-i-1])
		}
		i++
	}

	return latestAdditions
}

func findMissingAlbums(scanned []status.Album, known []status.Album) []status.Album {
	scannedMap := make(map[status.Album]int)
	for _, album := range scanned {
		scannedMap[album] = present
	}

	// Build a list of known albums not present in current scan
	var missingList []status.Album
	for _, album := range known {
		if scannedMap[album] != present {
			missingList = append(missingList, album)
		}
	}

	return missingList
}

func backupKnownAlbums(knownAlbumsBackupPath string, albums []status.Album) {
	os.Remove(knownAlbumsBackupPath)
	if _, err := os.Stat(knownAlbumsBackupPath); os.IsNotExist(err) {
		file, errCreate := os.Create(knownAlbumsBackupPath)
		if errCreate != nil {
			log.Printf("Could not create known albums backup at %v", knownAlbumsBackupPath)
			panic(errCreate)
		}
		file.Close()
	}

	knownAlbumsBackupFile, _ := os.OpenFile(knownAlbumsBackupPath, os.O_RDWR|os.O_APPEND, 0660)
	defer knownAlbumsBackupFile.Close()
	kaWriter := bufio.NewWriter(knownAlbumsBackupFile)
	for _, album := range albums {
		if _, err := kaWriter.WriteString(fmt.Sprintf("%v\n", album.Text)); err != nil {
			log.Printf("Could not write album to %v", knownAlbumsBackupPath)
			panic(err)
		}
	}
	if err := kaWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", knownAlbumsBackupPath)
		panic(err)
	}
}

func ensureFileEndsInNewline(fileLocation string) {
	file, _ := os.OpenFile(fileLocation, os.O_RDWR|os.O_APPEND, 0660)
	defer file.Close()
	fileInfo, _ := file.Stat()
	buf := []byte{' '}
	file.ReadAt(buf, fileInfo.Size()-1)
	if buf[0] != '\n' && fileInfo.Size() > 0 {
		file.Write([]byte{'\n'})
	}
}

func readFile(fileLocation string) (albums []status.Album, lineMap map[string]int) {
	file, _ := os.Open(fileLocation)
	scanner := bufio.NewScanner(file)
	lineMap = make(map[string]int)
	for scanner.Scan() {
		albums = append(albums, status.Album{scanner.Text()})
		lineMap[scanner.Text()] = present
	}
	file.Close()
	return albums, lineMap
}

func (slice sortableAlbums) Len() int {
	return len(slice)
}

func (slice sortableAlbums) Less(i, j int) bool {
	return slice[i].Text < slice[j].Text
}

func (slice sortableAlbums) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByName(albums sortableAlbums) sortableAlbums {
	sort.Sort(albums)
	return albums
}

func reverseSortByName(albums sortableAlbums) sortableAlbums {
	sort.Sort(sort.Reverse(albums))
	return albums
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
