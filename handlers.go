package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	ics, err := FetchICS(json.Url)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newIcs, err := TransformCalendar(ics, json.ReplacementSummary)
	if err != nil {
		// TODO: consider masking this error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	WriteRecord(json)
	// c.JSON(http.StatusOK, gin.H{"status": "successfully added entry"})
	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.String(http.StatusOK, newIcs)

}

func GetCalendarByID(c *gin.Context) {
	id := c.Param("id")
	url, err := ReadRecord(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "url": url})
}
