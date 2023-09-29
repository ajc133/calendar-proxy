package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type CalendarParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

func main() {
	initDB()
	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.GET("/:id", getCalendarByID)
	router.POST("/", createCalendar)
	router.Run("localhost:8080")
}

func createCalendar(c *gin.Context) {
	json := CalendarParams{}
	err := c.Bind(&json)
	if err != nil {
		// TODO: consider simplifying this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ics, err := fetchICS(json.Url)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newIcs, err := transformCalendar(ics, json.ReplacementSummary)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	writeRecord(json)
	// c.JSON(http.StatusOK, gin.H{"status": "successfully added entry"})
	c.Header("Content-Type", "text-calendar; charset=utf-8")
	c.String(http.StatusOK, newIcs)

}

func transformCalendar(body string, replacementSummary string) (string, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(body))
	newCal := ics.NewCalendar()
	log.Println("Attempting to transform")

	if err != nil {
		return "", err
	}
	log.Println(cal.Serialize())

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

func fetchICS(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	isCalendar := strings.Contains(resp.Header.Get("Content-Type"), "text/calendar")
	if !isCalendar {
		log.Print("Invalid content-type: Got ", resp.Header.Get("Content-Type"))
		return "", fmt.Errorf("URL is not a calendar")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

func getCalendarByID(c *gin.Context) {
	id := c.Param("id")
	url, err := readRecord(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "url": url})
}

func initDB() error {
	db, err := sql.Open("sqlite3", "./calendars.db")
	stmt := "CREATE TABLE IF NOT EXISTS calendars(id INTEGER PRIMARY KEY, url TEXT, replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func readRecord(id string) (string, error) {
	db, err := sql.Open("sqlite3", "./calendars.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	if err != nil {
		return "", err
	}

	stmt, err := db.Prepare("select url from calendars where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var url string
	err = stmt.QueryRow("1").Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil

}

func writeRecord(params CalendarParams) {
	db, err := sql.Open("sqlite3", "./calendars.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt := "INSERT INTO calendars(url, replacementSummary) VALUES(?, ?);"
	result, err := db.Exec(stmt, params.Url, params.ReplacementSummary)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.LastInsertId())
}
