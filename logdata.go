package main

import (
	"log/slog"
	"slices"
	"strings"
	"time"
)

type LogItem struct {
	Time     time.Time
	Message  string
	Severity int
}

const (
	RecentMinute = iota
	RecentHour
	RecentDay
	AllTime

	LevelDebug   = int(slog.LevelDebug)
	LevelInfo    = int(slog.LevelInfo)
	LevelWarning = int(slog.LevelWarn)
	LevelError   = int(slog.LevelError)
	LevelPanic   = int(slog.LevelError + 1)
)

func QueryLogItems(timeInterval int, messageContains string, severities []int) []LogItem {

	all := allLogData()

	var afterTime time.Time

	switch timeInterval {
	case RecentMinute:
		afterTime = time.Now().Add(-time.Minute).Add(-time.Nanosecond)
	case RecentHour:
		afterTime = time.Now().Add(-time.Hour).Add(-time.Nanosecond)
	case RecentDay:
		afterTime = time.Now().AddDate(0, 0, -1).Add(-time.Nanosecond)
	case AllTime:
	default:
		afterTime = time.Date(0, 0, 0, 0, 0, 0, 0, nil)
	}

	results := []LogItem{}

	for _, item := range all {
		if item.Time.After(afterTime) {
			if slices.Contains(severities, item.Severity) {
				if len(messageContains) == 0 || strings.Contains(item.Message, messageContains) {
					results = append(results, item)
				}
			}

		}
	}

	return results
}

func allLogData() []LogItem {
	// Base all log times from 7 days ago
	baseTime := time.Now().AddDate(0, 0, -5)

	makeTime := func(relativeDay int, hour int, minute int) time.Time {

		hoursAndMinutes := int64(minute)*int64(time.Minute) + int64(hour)*int64(time.Hour)

		return baseTime.AddDate(0, 0, relativeDay).Add(time.Duration(hoursAndMinutes))
	}

	return []LogItem{
		{makeTime(0, 0, 0), "program quit unexpectadly", LevelPanic},
		{makeTime(1, 1, 1), "flushStream() function was called", LevelDebug},
		{makeTime(1, 3, 9), "incoming connection from 192.168.1.234 was established", LevelInfo},
		{makeTime(1, 3, 40), "incoming connection from was 192.168.1.234 closed normally", LevelInfo},
		{makeTime(1, 7, 0), "buffer channels are more than 80 %% full", LevelWarning},
		{makeTime(2, 10, 3), "file /var/log/stream_92034023.log could not be opened", LevelError},
		{makeTime(2, 10, 29), "flushStream() function was called", LevelDebug},
		{makeTime(2, 10, 41), "incoming connection from 192.168.1.92 was established", LevelInfo},
		{makeTime(2, 11, 42), "buffer channels are more than 80 %% full", LevelWarning},
		{makeTime(5, -1, 3), "incoming connection from was 192.168.1.92 closed normally", LevelInfo},
		{makeTime(5, 0, -2), "file /var/log/stream_563366.log could not be opened", LevelError},
	}
}
