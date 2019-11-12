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

type sortableStrings []string

type KnownAlbumsWithBackup struct{}

// UpdateKnownAlbums updates the known albums file and returns an array
// of new albums, based upon the array of folder passed in as the
// folderScanList argument.
func (k *KnownAlbumsWithBackup) UpdateKnownAlbums(folderScanList []status.Album, knownAlbumsPath string, knownAlbumsBackupPath string, limit int) []status.Album {
	if _, err := os.Stat(knownAlbumsPath); os.IsNotExist(err) {
		file, errCreate := os.Create(knownAlbumsPath)
		if errCreate != nil {
			log.Printf("Could not create known albums at %v", knownAlbumsPath)
			panic(errCreate)
		}
		file.Close()
	}

	ensureFileEndsInNewline(knownAlbumsPath)

	// Read knownalbums into an array of its lines and a map that conveys if a line is present
	knownAlbumsLines, knownAlbumsMap := readFile(knownAlbumsPath)

	// Build a list of current scan entries not present in known albums (new albums)
	var newAlbums []string
	for _, scanItem := range folderScanList {
		if knownAlbumsMap[scanItem] != present {
			newAlbums = append(newAlbums, scanItem)
		}
	}

	reverseSortByName(newAlbums)
	log.Printf("Found %v known and %v new albums", len(knownAlbumsMap), len(newAlbums))

	missingAlbums := findMissingAlbums(folderScanList, knownAlbumsLines)

	if len(missingAlbums) > 0 {
		log.Printf("Found %v missing albums", len(missingAlbums))
		for i, missing := range missingAlbums {
			log.Printf("Missing #%v: %v", i+1, missing)
		}
	}

	// Append new albums to known albums file
	var knownAlbumsFile *os.File
	knownAlbumsFile, err := os.OpenFile(knownAlbumsPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("Could not open known albums file for appending: %v", knownAlbumsPath)
		panic(err)
	}
	defer knownAlbumsFile.Close()
	krWriter := bufio.NewWriter(knownAlbumsFile)
	for _, newAlbum := range newAlbums {
		if _, err := krWriter.WriteString(fmt.Sprintf("%v\n", newAlbum)); err != nil {
			log.Printf("Could not write new album to %v", knownAlbumsPath)
			panic(err)
		}
	}
	if err := krWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", knownAlbumsPath)
		panic(err)
	}

	if len(newAlbums) > 0 {
		backupKnownAlbums(knownAlbumsBackupPath, append(knownAlbumsLines, newAlbums...))
	}

	// Return sorted new albums then knownalbums from the end, up to a total of limit
	sortByName(newAlbums)
	var latestAdditions []string
	i := 0
	for i < min(len(newAlbums), limit) {
		latestAdditions = append(latestAdditions, newAlbums[i])
		i++
	}

	i = 0
	for i < (limit - len(newAlbums)) {
		if i < len(knownAlbumsLines) {
			latestAdditions = append(latestAdditions, knownAlbumsLines[len(knownAlbumsLines)-i-1])
		}
		i++
	}

	return latestAdditions
}

func findMissingAlbums(scanned []status.Album, known []string) []string {
	scannedMap := make(map[string]int)
	for _, album := range scanned {
		scannedMap[album] = present
	}

	// Build a list of known albums not present in current scan
	var missingList []string
	for _, album := range known {
		if scannedMap[album] != present {
			missingList = append(missingList, album)
		}
	}

	return missingList
}

func backupKnownAlbums(knownAlbumsBackupPath string, albums []string) {
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
	krWriter := bufio.NewWriter(knownAlbumsBackupFile)
	for _, album := range albums {
		if _, err := krWriter.WriteString(fmt.Sprintf("%v\n", album)); err != nil {
			log.Printf("Could not write album to %v", knownAlbumsBackupPath)
			panic(err)
		}
	}
	if err := krWriter.Flush(); err != nil {
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

func readFile(fileLocation string) (lines []string, lineMap map[string]int) {
	file, _ := os.Open(fileLocation)
	scanner := bufio.NewScanner(file)
	lineMap = make(map[string]int)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineMap[scanner.Text()] = present
	}
	file.Close()
	return lines, lineMap
}

func (slice sortableStrings) Len() int {
	return len(slice)
}

func (slice sortableStrings) Less(i, j int) bool {
	return slice[i] < slice[j]
}

func (slice sortableStrings) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByName(strings sortableStrings) sortableStrings {
	sort.Sort(strings)
	return strings
}

func reverseSortByName(strings sortableStrings) sortableStrings {
	sort.Sort(sort.Reverse(strings))
	return strings
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
