package telly

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Session is a connected user to the application
//
type Session struct {
	// Is the stream ready to use
	Ready bool

	// PID of the ffmpeg transcoding process
	TranscodingCmd *exec.Cmd

	// Timestamp when the stream was last read
	LastRead time.Time

	// Channel name which is tuned
	TunedChannel Channel
}

func NewSession(channel Channel) (*Session, error) {
	session := &Session{
		TunedChannel: channel,
	}
	return session, session.Start()
}

func (s *Session) TranscodedFilePath() string {
	return filepath.Join(TranscodingPath(), fmt.Sprintf("%s.m3u8", s.TunedChannel.GuideNumber))
}

func (s *Session) Stop() {
	log.WithField("name", s.TunedChannel.GuideNumber).Debug("Session#Stop")
	s.LastRead = time.Time{}

	// Stop transcoding process
	if s.TranscodingCmd != nil {
		if err := s.TranscodingCmd.Process.Kill(); err != nil {
			log.WithError(err).Error("Error killing process")
		} else {
			log.WithField("name", s.TunedChannel.GuideNumber).Info("Process killed")
		}
	}

	// Clean up tmp files?
}

func (s *Session) Start() error {
	s.LastRead = time.Now()

	// Get command line for ffmpeg command
	args := ffmpegCommandArgs(s.TunedChannel.GuideNumber, s.TranscodedFilePath(), "", 0)

	// Create an *exec.Cmd
	s.TranscodingCmd = exec.Command(args[0], args[1:]...)

	//
	// Run
	log.WithField("cmd", args).Debug("Starting transcoding process")
	err := s.TranscodingCmd.Start()
	if err != nil {
		log.WithError(err).Error("Cannot start transcoding process")
		s.Stop()
	} else {
		log.WithField("PID", s.TranscodingCmd.Process.Pid).Info("Started transcoding process")
	}
	go func() {
		if err := s.TranscodingCmd.Wait(); err != nil {
			log.WithField("PID", s.TranscodingCmd.Process.Pid).WithError(err).Error("Failed to Wait for Process")
			s.Stop()
		}
	}()
	return err
}

func (s *Session) IsValid() bool {
	// Check if process is still running
	p, err := os.FindProcess(s.TranscodingCmd.Process.Pid)
	log.WithField("p", p).Debug("Found process")
	return p != nil && err == nil
}
