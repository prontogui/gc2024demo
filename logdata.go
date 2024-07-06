package main

import (
	"strings"
	"time"
)

type LogItem struct {
	Time     time.Time
	Message  string
	Severity int
}

const (
	RecentMinute = "Last Minute"
	RecentHour   = "Last Hour"
	RecentDay    = "Last Day"
	AllTime      = "All"
)

const (
	// The severity levels to filter by.  These also server as array indices [0, n-1]
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelPanic
	LevelAllocation
)

var logData []LogItem

func init() {
	logData = initialLogData()
}

func QueryLogItems(timeInterval string, debugYes bool, infoYes bool, warningYes bool, errorYes bool, panicYes bool, messageContains string) []LogItem {

	// The filter criteria expressed in time.Time
	var afterTime time.Time

	// Determine the (inclusive) time frame to filter by
	switch timeInterval {
	case RecentMinute:
		afterTime = time.Now().Add(-time.Minute).Add(-time.Nanosecond)
	case RecentHour:
		afterTime = time.Now().Add(-time.Hour).Add(-time.Nanosecond)
	case RecentDay:
		afterTime = time.Now().AddDate(0, 0, -1).Add(-time.Nanosecond)
	case AllTime:
		afterTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	default:
		afterTime = time.Now().Add(time.Nanosecond)
	}

	// The results that are built from filtering
	results := []LogItem{}

	// Go through all the log items while filtering the right ones
	for _, item := range logData {

		// Inside the time interval?
		if !item.Time.After(afterTime) {
			continue
		}

		// One of the filtered severities?
		if item.Severity == LevelDebug && !debugYes {
			continue
		}
		if item.Severity == LevelInfo && !infoYes {
			continue
		}
		if item.Severity == LevelWarning && !warningYes {
			continue
		}
		if item.Severity == LevelError && !errorYes {
			continue
		}
		if item.Severity == LevelPanic && !panicYes {
			continue
		}

		// Message contains text?
		if !strings.Contains(item.Message, messageContains) {
			continue
		}

		// All filter criteria were met - include this item
		results = append(results, item)
	}

	return results
}

func initialLogData() []LogItem {
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
