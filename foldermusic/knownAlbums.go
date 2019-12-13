package foldermusic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/epond/porthole/shared"
)

const present = 1

type sortableAlbums []shared.Album

// KnownAlbumsWithBackup makes a backup each time the known albums are updated
type KnownAlbumsWithBackup struct {
	KnownAlbumsPath       string
	KnownAlbumsBackupPath string
}

func (k *KnownAlbumsWithBackup) ReadKnownAlbums() (albums []shared.Album, lineMap map[string]int) {
	if _, err := os.Stat(k.KnownAlbumsPath); os.IsNotExist(err) {
		file, errCreate := os.Create(k.KnownAlbumsPath)
		if errCreate != nil {
			log.Printf("Could not create known albums at %v", k.KnownAlbumsPath)
			panic(errCreate)
		}
		file.Close()
	}

	// Read knownalbums into an array of its lines and a map that conveys if a line is present
	albums, lineMap = readFile(k.KnownAlbumsPath)
	return albums, lineMap
}

func (k *KnownAlbumsWithBackup) AppendNewAlbums(knownAlbums []shared.Album, newAlbums []shared.Album) {
	ensureFileEndsInNewline(k.KnownAlbumsPath)
	var knownAlbumsFile *os.File
	knownAlbumsFile, err := os.OpenFile(k.KnownAlbumsPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("Could not open known albums file for appending: %v", k.KnownAlbumsPath)
		panic(err)
	}
	defer knownAlbumsFile.Close()
	kaWriter := bufio.NewWriter(knownAlbumsFile)
	for _, newAlbum := range newAlbums {
		if _, err := kaWriter.WriteString(fmt.Sprintf("%v\n", newAlbum.Text)); err != nil {
			log.Printf("Could not write new album to %v", k.KnownAlbumsPath)
			panic(err)
		}
	}
	if err := kaWriter.Flush(); err != nil {
		log.Printf("Could not flush %v", k.KnownAlbumsPath)
		panic(err)
	}

	if len(newAlbums) > 0 {
		backupKnownAlbums(k.KnownAlbumsBackupPath, append(knownAlbums, newAlbums...))
	}
}

func backupKnownAlbums(knownAlbumsBackupPath string, albums []shared.Album) {
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

func readFile(fileLocation string) (albums []shared.Album, lineMap map[string]int) {
	file, _ := os.Open(fileLocation)
	scanner := bufio.NewScanner(file)
	lineMap = make(map[string]int)
	for scanner.Scan() {
		albums = append(albums, shared.Album{scanner.Text()})
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
