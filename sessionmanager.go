package telly

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"time"
)

const (
	// How often to check for sessions (dead code, waiting)
	CheckSessionsTimeInterval = 1
)

// Session Manager for connected sessions
//
var DefaultSessionManager = NewSessionManager()

// Manager for active sessions. A session represents - viewer watching streaming, trancoding process, etc
//
type SessionManager struct {
	// List of active sessions
	sessions []*Session

	// Ticker to check for dead sessions
	ticker *time.Ticker

	// List of Channels
	channels []Channel
}

// Factory Method to get a new session manager
//
func NewSessionManager() *SessionManager {
	return &SessionManager{
		ticker: time.NewTicker(CheckSessionsTimeInterval * time.Second),
	}
}

// Find the channel that matches the given channel number (`guideNumber`)
//
func (m *SessionManager) FindChannel(guideNumber string) (Channel, bool) {
	var channel Channel
	var found bool
	for _, c := range m.channels {
		if c.GuideNumber == guideNumber {
			channel = c
			found = true
			break
		}
	}
	return channel, found
}

// Find the active session that matches the given channel number (`guideNumber`)
func (m *SessionManager) FindSession(guideNumber string) (*Session, bool) {
	var session *Session
	for _, s := range m.sessions {
		if s.TunedChannel.GuideNumber == guideNumber {
			session = s
			break
		}
	}
	return session, session != nil
}

// Create a new session with the given `channelName`. If successful, session is added
// to list of managed sessions.
//
func (m *SessionManager) NewSession(guideNumber string) (*Session, error) {
	channel, found := m.FindChannel(guideNumber)
	if !found {
		return nil, fmt.Errorf("%s channel not found", guideNumber)
	}
	session, err := NewSession(channel)
	if err == nil {
		m.sessions = append(m.sessions, session)
	} else {
		session = nil
	}
	return session, err
}

// Start the session manager.
//
func (m *SessionManager) Start() {
	var err error
	if m.channels, err = GetChannelLineUp(); err != nil {
		log.WithError(err).Error("Cannot get channel lineup")
		return
	}

	go func() {
		for t := range m.ticker.C {
			m.SweepSessions(t)
		}
	}()
}

// Stop the session manager
//
func (m *SessionManager) Stop() {
	m.ticker.Stop()
	for i := len(m.sessions) - 1; i >= 0; i-- {
		s := m.sessions[i]
		s.Stop()
	}
}

// Sweep the active sesions for dead or useless sessions
//
func (m *SessionManager) SweepSessions(t time.Time) {
	for i := len(m.sessions) - 1; i >= 0; i-- {
		s := m.sessions[i]
		// TODO: Can this be in goroutine
		m.sweepSession(s, i)
	}
}

func (m *SessionManager) sweepSession(s *Session, i int) {
	log.WithField("name", s.TunedChannel.GuideNumber).Debug("Sweeping Session")

	// check the status of the stream and if its ready to stream yet
	// If it isn't ready, check to see if it is ready
	if !s.Ready {
		log.WithField("name", s.TunedChannel.GuideNumber).Info("Checking session ready status")
		path := s.TranscodedFilePath()
		log.WithField("name", s.TunedChannel.GuideNumber).WithField("path", path).Debug("File path for stream")
		if _, err := os.Stat(path); err == nil {
			log.WithField("name", s.TunedChannel.GuideNumber).Info("Session is ready")
			s.Ready = true
		}
	}

	destroy := false
	// check to see when the last time the stream was accessed
	// if it was longer than 60 seconds, kill the session
	if time.Since(s.LastRead) > 60*time.Second {
		destroy = true
	} else {
		// Check for dead process
		destroy = !s.IsValid()
	}

	if destroy {
		log.WithField("name", s.TunedChannel.GuideNumber).Info("Killing session")
		s.Stop()
		m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
	}
}
