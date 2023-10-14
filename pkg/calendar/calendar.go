package calendar

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

const DatabaseFileName string = "data/calendars.db"

func FetchAndTransformCalendar(url string, replacementSummary string) (string, error) {
	cal, err := fetchCalendar(url)
	if err != nil {
		log.Printf("Error fetching url (%s): %s", url, err)
		return "", err
	}

	newCal, err := transformCalendar(cal, replacementSummary)
	if err != nil {
		log.Printf("Error transforming calendar from url (%s): %s", url, err)
		return "", err
	}
	return newCal, nil
}

func isCalendarContentType(header string) bool {
	return strings.Contains(header, "text/calendar") ||
		strings.Contains(header, "application/octet-stream")
}

func fetchCalendar(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if !isCalendarContentType(resp.Header.Get("Content-Type")) {
		log.Print("Invalid content-type: Got ", resp.Header.Get("Content-Type"))
		return "", fmt.Errorf("URL is not a calendar")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

func transformCalendar(body string, replacementSummary string) (string, error) {
	newCal := ics.NewCalendar()
	cal, err := ics.ParseCalendar(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	for _, event := range cal.Events() {
		newEvent := copyBarebonesEvent(event)
		newEvent.SetSummary(replacementSummary)

		newCal.AddVEvent(&newEvent)

	}

	return newCal.Serialize(), nil
}

func copyBarebonesEvent(event *ics.VEvent) ics.VEvent {
	id := uuid.New().String()
	newEvent := ics.NewEvent(id)

	componentPropertiesToCopy := []ics.ComponentProperty{
		ics.ComponentPropertyDtStart,
		ics.ComponentPropertyDtEnd,
		ics.ComponentPropertyRrule,
		ics.ComponentPropertyRdate,
	}

	for _, prop := range componentPropertiesToCopy {
		toCopy := event.GetProperty(prop)
		if toCopy != nil {
			newEvent.Properties = append(newEvent.Properties, *toCopy)
		}
	}
	return *newEvent
}
