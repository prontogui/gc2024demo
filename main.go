package main

import (
	"fmt"
	"time"

	pg "github.com/prontogui/golib"
)

func main() {

	pgui := pg.NewProntoGUI()
	err := pgui.StartServing("127.0.0.1", 50053)

	if err != nil {
		fmt.Printf("Error trying to start server:  %s", err.Error())
		return
	}

	logItems := QueryLogItems(AllTime, "", []int{LevelDebug, LevelInfo, LevelWarning, LevelError, LevelPanic})

	// Convert messages into PG model
	rows := [][]pg.Primitive{}

	for _, item := range logItems {
		rows = append(rows, []pg.Primitive{
			timePrimitive(item.Time),
			severityPrimitive(item.Severity),
			messagePrimitive(item.Message),
		})
	}

	table := pg.TableWith{
		TemplateRow: []pg.Primitive{
			&pg.Text{}, &pg.Text{}, &pg.Text{},
		},
		Headings: []string{
			"Time", "Severity", "Message",
		},
		Rows: rows,
	}.Make()

	pgui.SetGUI(table)

	for {
		_, err := pgui.Wait()
		if err != nil {
			fmt.Printf("error from Wait() is:  %s\n", err.Error())
			break
		}
	}
}

func timePrimitive(t time.Time) pg.Primitive {
	timeString := t.Format(time.RFC1123)
	return pg.TextWith{Content: timeString}.Make()
}

func severityPrimitive(severity int) pg.Primitive {

	var severityString string

	switch severity {
	case LevelDebug:
		severityString = "DEBUG"
	case LevelInfo:
		severityString = "INFO"
	case LevelWarning:
		severityString = "WARNING"
	case LevelError:
		severityString = "ERROR"
	case LevelPanic:
		severityString = "PANIC"
	default:
		severityString = "UNKNOWN"
	}

	return pg.TextWith{Content: severityString}.Make()
}

func messagePrimitive(message string) pg.Primitive {
	return pg.TextWith{Content: message}.Make()
}
