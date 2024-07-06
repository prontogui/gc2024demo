package main

import (
	"fmt"
	"time"

	pg "github.com/prontogui/golib"
)

func main() {

	// Initialize ProntoGUI
	pgui := pg.NewProntoGUI()
	err := pgui.StartServing("127.0.0.1", 50053)

	if err != nil {
		fmt.Printf("Error trying to start server:  %s", err.Error())
		return
	}

	// Initial filter settings (show everything)
	timeInterval := AllTime
	severities := [LevelAllocation]bool{true, true, true, true, true}
	messageContains := ""

	// Build our GUI

	// Big and bold heading for our GUI
	guiHeading := pg.TextWith{Content: "Simple Log Viewer"}.Make()

	// Time interval choice
	timeHeading := pg.TextWith{Content: "Time Interval"}.Make()
	timeChoice := pg.ChoiceWith{Choice: AllTime, Choices: []string{RecentMinute, RecentHour, RecentDay, AllTime}}.Make()
	timeGroup := pg.GroupWith{
		GroupItems: []pg.Primitive{timeHeading, timeChoice},
	}.Make()

	// Severities check boxes
	severitiesHeading := pg.TextWith{Content: "Severities"}.Make()
	debugCheck := pg.CheckWith{Label: "Debug", Checked: severities[LevelDebug]}.Make()
	infoCheck := pg.CheckWith{Label: "Info", Checked: severities[LevelInfo]}.Make()
	warningCheck := pg.CheckWith{Label: "Warning", Checked: severities[LevelWarning]}.Make()
	errorCheck := pg.CheckWith{Label: "Error", Checked: severities[LevelError]}.Make()
	panicCheck := pg.CheckWith{Label: "Panic", Checked: severities[LevelPanic]}.Make()
	severitiesGroup := pg.GroupWith{
		GroupItems: []pg.Primitive{severitiesHeading, debugCheck, infoCheck, warningCheck, errorCheck, panicCheck},
	}.Make()

	// Message text filter
	messageHeading := pg.TextWith{Content: "Message Contains"}.Make()
	messageTextField := pg.TextFieldWith{}.Make()
	messageGroup := pg.GroupWith{
		GroupItems: []pg.Primitive{messageHeading, messageTextField},
	}.Make()

	// A table to show log items
	table := pg.TableWith{
		TemplateRow: []pg.Primitive{
			&pg.Text{}, &pg.Text{}, &pg.Text{},
		},
		Headings: []string{
			"Time", "Severity", "Message",
		},
	}.Make()

	pgui.SetGUI(guiHeading, timeGroup, severitiesGroup, messageGroup, table)

	for {
		// Query for the log items as of this moment
		logItems := QueryLogItems(timeInterval, severities, messageContains)

		// Convert messages into a table row of ProntoGUI primitives
		rows := [][]pg.Primitive{}

		for _, item := range logItems {
			rows = append(rows, []pg.Primitive{
				timePrimitive(item.Time),
				severityPrimitive(item.Severity),
				messagePrimitive(item.Message),
			})
		}

		// Update the table contents
		table.SetRows(rows)

		// Wait for something to happen in the GUI
		_, err := pgui.Wait()
		if err != nil {
			fmt.Printf("error from Wait() is:  %s\n", err.Error())
			break
		}

		//
		timeInterval = timeChoice.Choice()
		severities[LevelDebug] = debugCheck.Checked()
		severities[LevelInfo] = infoCheck.Checked()
		severities[LevelWarning] = warningCheck.Checked()
		severities[LevelError] = errorCheck.Checked()
		severities[LevelPanic] = panicCheck.Checked()
		messageContains = messageTextField.TextEntry()
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
