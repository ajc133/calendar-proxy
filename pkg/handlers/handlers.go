package handlers

import (
	"log"
	"net/http"

	"github.com/ajc133/calendarproxy/pkg/calendar"
	"github.com/ajc133/calendarproxy/pkg/db"
	"github.com/gin-gonic/gin"
)

const ContentType string = "Content-Type"
const CalendarContent string = "text/calendar; charset=utf-8"
const DatabaseFileName string = "calendars.db"

type CalendarParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

func CreateCalendar(c *gin.Context) {
	json := CalendarParams{}
	err := c.Bind(&json)
	if err != nil {
		log.Printf("Error: %s", err)

		// TODO: User-friendly error for form submission
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Parameters submitted were malformed"})
		return
	}

	cal, err := calendar.FetchICS(json.Url)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when fetching the given URL"})
		return
	}

	log.Printf("Going to replace '%s' SUMMARY with '%s',", json.Url, json.ReplacementSummary)
	newCal, err := calendar.TransformCalendar(cal, json.ReplacementSummary)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Failed to parse fetched calendar body"})
		return
	}

	record := db.Record{
		Url:                json.Url,
		ReplacementSummary: json.ReplacementSummary,
		CalendarBody:       newCal,
	}

	// TODO: schedule a cronjob to periodically refresh this entry
	id, err := db.WriteRecord(DatabaseFileName, record)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error storing calendar in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func GetCalendarByID(c *gin.Context) {
	id := c.Param("id")
	calendarBody, err := db.ReadRecord(DatabaseFileName, id)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error in database lookup"})
		return
	}

	if calendarBody == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Couldn't find calendar with that id"})
		return
	}

	c.Header(ContentType, CalendarContent)
	c.String(http.StatusOK, calendarBody)
}
