package main

import (
	"os"
	"bufio"
)

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesPath string, limit int) []string {
	_, err := os.Stat(knownReleasesPath)
	if err != nil {
		file, _ := os.Create(knownReleasesPath)
		defer file.Close()
	}

	// Read knownreleases into an array of its lines and a map
	knownReleasesFile, _, _ := readFile(knownReleasesPath)
	defer knownReleasesFile.Close()

	// Build a map of the current scan
	folderScanMap := make(map[string]int)
	for _, item := range folderScanList {
		folderScanMap[item.String()] = present
	}

	// Build a list of current scan entries not present in known releases (new releases)
	// Sort new releases by name
	// Append new releases to known releases file
	// Return sorted new releases then knownreleases from the end, up to a total of limit

	if len(folderScanList) >= 3 {
		return []string{folderScanList[0].String(), folderScanList[1].String(), folderScanList[2].String()}
	}
	return []string{"not enough folder infos"}
}

const present = 1

func readFile(fileLocation string) (file *os.File, lines []string, lineMap map[string]int) {
	file, _ = os.Open(fileLocation)
	scanner := bufio.NewScanner(file)
	lineMap = make(map[string]int)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineMap[scanner.Text()] = present
	}
	return file, lines, lineMap
}