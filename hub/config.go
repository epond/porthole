package hub

import (
	"os"
	"strconv"
)

// Config is configuration for Porthole
type Config struct {
	MusicFolder              string
	KnownAlbumsFile          string
	KnownAlbumsBackup        string
	FoldersToScan            string
	LogFile                  string
	GitCommit                string
	FetchInterval            int
	DashboardRefreshInterval int
	SleepAfter               int
	LatestAdditionsLimit     int
}

// ConfigFromEnv creates a new Config from the Environment
func ConfigFromEnv() *Config {
	fetchInterval, _ := strconv.Atoi(os.Getenv("FETCH_INTERVAL"))
	dashboardRefreshInterval, _ := strconv.Atoi(os.Getenv("DASHBOARD_REFRESH_INTERVAL"))
	sleepAfter, _ := strconv.Atoi(os.Getenv("SLEEP_AFTER"))
	latestAdditionsLimit, _ := strconv.Atoi(os.Getenv("LATEST_ADDITIONS_LIMIT"))
	return &Config{
		MusicFolder:              os.Getenv("MUSIC_FOLDER"),
		KnownAlbumsFile:          os.Getenv("KNOWN_ALBUMS_FILE"),
		KnownAlbumsBackup:        os.Getenv("KNOWN_ALBUMS_BACKUP"),
		GitCommit:                os.Getenv("GIT_COMMIT"),
		LogFile:                  os.Getenv("LOG_FILE"),
		FetchInterval:            fetchInterval,
		DashboardRefreshInterval: dashboardRefreshInterval,
		SleepAfter:               sleepAfter,
		FoldersToScan:            os.Getenv("FOLDERS_TO_SCAN"),
		LatestAdditionsLimit:     latestAdditionsLimit,
	}
}
