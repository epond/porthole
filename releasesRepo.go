package main

func UpdateKnownReleases(folderScanList []FolderInfo, knownReleasesFile string, limit int) []string {
	return []string{folderScanList[0].String(), folderScanList[1].String(), folderScanList[2].String()}
}