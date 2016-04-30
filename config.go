package telly

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(AppConfig.LoggingLevel)
}

// Config is the struct holding various configurations for the application
//
type Config struct {
	// Log Level
	LoggingLevel log.Level
	// HDHomeRun host
	HDHomeRunHost string
	// ffmpeg application path
	FFMpegPath string
	// Port to run the application on
	Port string
}

// Load config from environment variables. Load some defaults.
func loadConfig() Config {
	level, _ := log.ParseLevel(getEnvString("LOG", "info"))
	config := Config{
		HDHomeRunHost: getEnvString("HDHOMERUN_HOST", "10.0.1.31"),
		FFMpegPath:    getEnvString("FFMPEG_PATH", "/usr/local/bin/ffmpeg"),
		Port:          getEnvString("PORT", "8888"),
		LoggingLevel:  level,
	}
	log.WithField("config", config).Debug("Loading config")
	return config
}

func TranscodingPath() string {
	path := os.TempDir()
	return path
}

var AppConfig = loadConfig()
