package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ContentType string = "Content-Type"
)

const (
	CalendarContent string = "text/calendar; charset=utf-8"
)

type CalendarParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

func CreateCalendar(c *gin.Context) {
	json := CalendarParams{}
	err := c.Bind(&json)
	if err != nil {
		// TODO: consider simplifying this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cal, err := FetchICS(json.Url)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Going to replace '%s' SUMMARY with '%s',", json.Url, json.ReplacementSummary)
	newCal, err := TransformCalendar(cal, json.ReplacementSummary)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := WriteRecord(json, newCal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header(ContentType, CalendarContent)
	// TODO: http status for 'created'
	c.JSON(http.StatusOK, gin.H{"id": id})

}

func GetCalendarByID(c *gin.Context) {
	id := c.Param("id")
	url, err := ReadRecord(id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "url": url})
}
