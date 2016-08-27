package main

import (
	"os"
	"bufio"
	"fmt"
	"sort"
	"log"
)

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesPath string, limit int) []string {
	if _, err := os.Stat(knownReleasesPath); os.IsNotExist(err) {
		file, errCreate := os.Create(knownReleasesPath)
		if (errCreate != nil) {
			log.Printf("Could not create known releases at %v", knownReleasesPath)
			panic(errCreate)
		}
		file.Close()
	}

	ensureFileEndsInNewline(knownReleasesPath)

	// Read knownreleases into an array of its lines and a map
	knownReleasesLines, knownReleasesMap := readFile(knownReleasesPath)
	log.Printf("Found %v known releases", len(knownReleasesMap))

	// Build a list of current scan entries not present in known releases (new releases)
	var newReleases []string
	for _, scanItem := range folderScanList {
		if knownReleasesMap[scanItem.String()] != present {
			newReleases = append(newReleases, scanItem.String())
		}
	}

	reverseSortByName(newReleases)
	log.Printf("Found %v new releases", len(newReleases))

	// Append new releases to known releases file
	knownReleasesFile, _ := os.OpenFile(knownReleasesPath, os.O_RDWR|os.O_APPEND, 0660)
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

	// Return sorted new releases then knownreleases from the end, up to a total of limit
	sortByName(newReleases)
	var latestAdditions []string
	i := 0
	for i < min(len(newReleases), limit) {
		latestAdditions = append(latestAdditions, newReleases[i])
		i += 1
	}

	i = 0
	for i < (limit - len(newReleases)) {
		if i < len(knownReleasesLines) {
			latestAdditions = append(latestAdditions, knownReleasesLines[len(knownReleasesLines)-i-1])
		}
		i += 1
	}

	return latestAdditions
}

const present = 1

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

type SortableStrings []string

func (slice SortableStrings) Len() int {
	return len(slice)
}

func (slice SortableStrings) Less(i, j int) bool {
	return slice[i] < slice[j];
}

func (slice SortableStrings) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortByName(strings SortableStrings) SortableStrings {
	sort.Sort(strings)
	return strings
}

func reverseSortByName(strings SortableStrings) SortableStrings {
	sort.Sort(sort.Reverse(strings))
	return strings
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}