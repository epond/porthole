package main

import (
	"os"
	"bufio"
	"log"
	"fmt"
)

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesPath string, limit int) []string {
	_, err := os.Stat(knownReleasesPath)
	if err != nil {
		file, _ := os.Create(knownReleasesPath)
		file.Close()
	}

	// Read knownreleases into an array of its lines and a map
	_, knownReleasesMap := readFile(knownReleasesPath)

	// Build a list of current scan entries not present in known releases (new releases)
	var newReleases []string
	for _, scanItem := range folderScanList {
		log.Printf("scanItem: %v", scanItem.String())
		if knownReleasesMap[scanItem.String()] != present {
			log.Print("not present")
			newReleases = append(newReleases, scanItem.String())
		}
	}

	// Sort new releases by name
	// TODO

	// Append new releases to known releases file
	knownReleasesFile, _ := os.OpenFile(knownReleasesPath, os.O_RDWR|os.O_APPEND, 0660)
	krWriter := bufio.NewWriter(knownReleasesFile)
	for _, newRelease := range newReleases {
		krWriter.WriteString(fmt.Sprintf("%v\n", newRelease))
	}
	if err = krWriter.Flush(); err != nil {
		panic(err)
	}
	knownReleasesFile.Close()

	// Return sorted new releases then knownreleases from the end, up to a total of limit

	if len(folderScanList) >= 3 {
		return []string{folderScanList[0].String(), folderScanList[1].String(), folderScanList[2].String()}
	}
	return []string{"not enough folder infos"}
}

const present = 1

func readFile(fileLocation string) (lines []string, lineMap map[string]int) {
	file, _ := os.Open(fileLocation)
	scanner := bufio.NewScanner(file)
	lineMap = make(map[string]int)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineMap[scanner.Text()] = present
		log.Printf("Found known release: %v", scanner.Text())
	}
	file.Close()
	return lines, lineMap
}