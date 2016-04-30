package telly

import (
	"fmt"
)

func addOption(command []string, option string, value interface{}) []string {
	command = append(command, option)
	if value != nil {
		command = append(command, fmt.Sprintf("%v", value))
	}
	return command
}

func ffmpegCommandArgs(channel string, outputFile string, resolution string, bitrate int) []string {
	if resolution == "" {
		resolution = "1280x720"
	}
	if bitrate == 0 {
		bitrate = 3000
	}
	var command []string

	// command name - ffmpeg
	command = append(command, AppConfig.FFMpegPath)

	options := []struct {
		Option string
		Value  interface{}
	}{
		// Input stream
		{"-i", RawStreamUrlForChannel(channel)},

		// Force Format
		//{"-f", "mpegts"},

		// Overwrite files
		{"-y", nil},

		// Threads (say auto for auto detect, default is 1)
		{"-threads", "auto"},

		// specify how many microseconds are analyzed to probe the input (from 0 to INT_MAX) (default 5e+06)
		{"-analyzeduration", "2000000"},

		// set number of audio channels (from INT_MIN to INT_MAX) (default 0)
		//{"-ac", "2"},

		// force audio codec ('copy' to copy stream)
		{"-acodec", "libfaac"},

		// Video bitrate
		{"-b:v", fmt.Sprintf("%dk", bitrate)},

		// set ratecontrol buffer size (in bits) (from INT_MIN to INT_MAX) (default 0)
		{"-bufsize", fmt.Sprintf("%dk", 2*bitrate)},

		// Set minimum bitrate tolerance (in bits/s). Most useful in setting up a CBR encode. It is of little use otherwise. (from INT_MIN to INT_MAX) (default 0)
		{"-minrate", fmt.Sprintf("%dk", int(0.8*float32(bitrate)))},

		// Set maximum bitrate tolerance (in bits/s). Requires bufsize to be set. (from INT_MIN to INT_MAX) (default 0)
		{"-maxrate", fmt.Sprintf("%dk", bitrate)},

		// force video codec ('copy' to copy stream)
		{"-vcodec", "libx264"},

		// force audio codec ('copy' to copy stream)
		{"-acodec", "libfaac"},

		// Video bitrate
		{"-b:v", fmt.Sprintf("%dk", bitrate)},

		// set ratecontrol buffer size (in bits) (from INT_MIN to INT_MAX) (default 0)
		{"-bufsize", fmt.Sprintf("%dk", 2*bitrate)},

		// Set minimum bitrate tolerance (in bits/s). Most useful in setting up a CBR encode. It is of little use otherwise. (from INT_MIN to INT_MAX) (default 0)
		{"-minrate", fmt.Sprintf("%dk", int(0.8*float32(bitrate)))},

		// Set maximum bitrate tolerance (in bits/s). Requires bufsize to be set. (from INT_MIN to INT_MAX) (default 0)
		{"-maxrate", fmt.Sprintf("%dk", bitrate)},

		// force video codec ('copy' to copy stream)
		{"-vcodec", "libx264"},

		// Size
		//{"-s", resolution},

		// Set the encoding preset (cf. x264 --fullhelp) (default "medium")
		{"-preset", "superfast"},

		// set frame rate (Hz value, fraction or abbreviation)
		{"-r", "29.97"},

		// set segment length in seconds (from 0 to FLT_MAX) (default 2)
		{"-hls_time", "2"},

		// set number after which the index wraps (from 0 to INT_MAX) (default 0)
		{"-hls_wrap", "40"},

		// set logging level
		{"-loglevel", "warning"},

		// simplified 1 parameter audio timestamp matching, 0(disabled), 1(filling and trimming), >1(maximum stretch/squeeze in samples per second) (from INT_MIN to INT_MAX) (default 0)
		{"-async", "1"},

		// Tune the encoding params
		{"-tune", "zerolatency"},

		{"-flags", "-global_header"},
		{"-fflags", "+genpts"},
		{"-fflags", "+genpts"},

		// -map [-]input_file_id[:stream_specifier][,sync_file_id[:stream_s  set input stream mapping
		{"-map", "0:0"},
		{"-map", "0:1"},
	}

	for _, option := range options {
		command = addOption(command, option.Option, option.Value)
	}

	command = append(command, outputFile)
	return command
}
