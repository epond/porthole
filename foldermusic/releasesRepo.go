package foldermusic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const present = 1

type sortableStrings []string

// UpdateKnownReleases updates the known releases file and returns an array
// of new releases, based upon the array of folder passed in as the
// folderScanList argument.
func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesPath string, knownReleasesBackupPath string, limit int) []string {
	if _, err := os.Stat(knownReleasesPath); os.IsNotExist(err) {
		file, errCreate := os.Create(knownReleasesPath)
		if errCreate != nil {
			log.Printf("Could not create known releases at %v", knownReleasesPath)
			panic(errCreate)
		}
		file.Close()
	}

	ensureFileEndsInNewline(knownReleasesPath)

	// Read knownreleases into an array of its lines and a map
	knownReleasesLines, knownReleasesMap := readFile(knownReleasesPath)

	// Build a list of current scan entries not present in known releases (new releases)
	var newReleases []string
	for _, scanItem := range folderScanList {
		if knownReleasesMap[scanItem.String()] != present {
			newReleases = append(newReleases, scanItem.String())
		}
	}

	reverseSortByName(newReleases)
	log.Printf("Found %v known and %v new releases", len(knownReleasesMap), len(newReleases))

	missingReleases := findMissingReleases(folderScanList, knownReleasesLines)

	if len(missingReleases) > 0 {
		log.Printf("Found %v missing releases", len(missingReleases))
		for i, missing := range missingReleases {
			log.Printf("Missing #%v: %v", i+1, missing)
		}
	}

	// Append new releases to known releases file
	var knownReleasesFile *os.File
	knownReleasesFile, err := os.OpenFile(knownReleasesPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("Could not open known releases file for appending: %v", knownReleasesPath)
		panic(err)
	}
	defer knownReleasesFile.Close()
	krWriter := bufio.NewWriter(knownReleasesFile)
	for _, newRelease := range newReleases {
		if _, err := krWriter.WriteString(fmt.Sprintf("%v\n", newRelease)); err != nil {
			log.Printf("Could not write new release to %v", knownReleasesPath)
			panic(err)
		}
	}
	if err := krWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", knownReleasesPath)
		panic(err)
	}

	if len(newReleases) > 0 {
		backupKnownReleases(knownReleasesBackupPath, append(knownReleasesLines, newReleases...))
	}

	// Return sorted new releases then knownreleases from the end, up to a total of limit
	sortByName(newReleases)
	var latestAdditions []string
	i := 0
	for i < min(len(newReleases), limit) {
		latestAdditions = append(latestAdditions, newReleases[i])
		i++
	}

	i = 0
	for i < (limit - len(newReleases)) {
		if i < len(knownReleasesLines) {
			latestAdditions = append(latestAdditions, knownReleasesLines[len(knownReleasesLines)-i-1])
		}
		i++
	}

	return latestAdditions
}

func findMissingReleases(scanned []FolderInfo, known []string) []string {
	scannedMap := make(map[string]int)
	for _, release := range scanned {
		scannedMap[release.String()] = present
	}

	// Build a list of known releases not present in current scan
	var missingList []string
	for _, release := range known {
		if scannedMap[release] != present {
			missingList = append(missingList, release)
		}
	}

	return missingList
}

func backupKnownReleases(knownReleasesBackupPath string, releases []string) {
	os.Remove(knownReleasesBackupPath)
	if _, err := os.Stat(knownReleasesBackupPath); os.IsNotExist(err) {
		file, errCreate := os.Create(knownReleasesBackupPath)
		if errCreate != nil {
			log.Printf("Could not create known releases backup at %v", knownReleasesBackupPath)
			panic(errCreate)
		}
		file.Close()
	}

	knownReleasesBackupFile, _ := os.OpenFile(knownReleasesBackupPath, os.O_RDWR|os.O_APPEND, 0660)
	defer knownReleasesBackupFile.Close()
	krWriter := bufio.NewWriter(knownReleasesBackupFile)
	for _, release := range releases {
		if _, err := krWriter.WriteString(fmt.Sprintf("%v\n", release)); err != nil {
			log.Printf("Could not write release to %v", knownReleasesBackupPath)
			panic(err)
		}
	}
	if err := krWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", knownReleasesBackupPath)
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
