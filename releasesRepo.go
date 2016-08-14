package main

import "os"

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesFile string, limit int) []string {
	_, err := os.Stat(knownReleasesFile)
	if err != nil {
		file, _ := os.Create(knownReleasesFile)
		defer file.Close()
	}

	if len(folderScanList) >= 3 {
		return []string{folderScanList[0].String(), folderScanList[1].String(), folderScanList[2].String()}
	}
	return []string{"not enough folder infos"}
}