package main

import "os"

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesFile string, limit int) []string {
	_, err := os.Stat(knownReleasesFile)
	if err != nil {
		file, _ := os.Create(knownReleasesFile)
		defer file.Close()
	}

	// Build a map of known releases
	// Build a map of the current scan
	// Build a list of current scan entries not present in known releases (new releases)
	// Sort new releases by name
	// Append new releases to known releases file

	// ?How should latest n releases be determined?

	if len(folderScanList) >= 3 {
		return []string{folderScanList[0].String(), folderScanList[1].String(), folderScanList[2].String()}
	}
	return []string{"not enough folder infos"}
}