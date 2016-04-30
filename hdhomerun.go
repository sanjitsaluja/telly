package telly

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// URL Factory for line up
func HdhomerunLineupUrl() string {
	return fmt.Sprintf("http://%s/lineup.json", AppConfig.HDHomeRunHost)
}

// URL Factory for raw stream from HDHomeRun
func RawStreamUrlForChannel(channel string) string {
	return fmt.Sprintf("http://%s:5004/auto/v%s", AppConfig.HDHomeRunHost, channel)
}

// Get the channel lineup
func GetChannelLineUp() ([]Channel, error) {
	lineupUrl := HdhomerunLineupUrl()
	log.WithField("url", lineupUrl).Info("Getting line up url")
	channels := []Channel{}
	err := getJson(lineupUrl, &channels)
	if err != nil {
		log.WithField("url", lineupUrl).WithError(err).Error("Error getting line up url")
	}
	return channels, err
}
