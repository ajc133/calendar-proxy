package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	// TODO: create table at server start
	router.StaticFile("/", "./static/index.html")
	router.GET("/:id", getCalendarByID)
	router.POST("/", createCalendar)
	router.Run("localhost:8080")
}

type CalendarParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

func createCalendar(c *gin.Context) {
	json := CalendarParams{}
	err := c.ShouldBind(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: validate url is a url

	writeRecord(json)
	c.JSON(http.StatusOK, gin.H{"status": "successfully added entry"})
}

func getCalendarByID(c *gin.Context) {
	id := c.Param("id")
	url, err := readRecord(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "url": url})
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

	stmt, err := db.Prepare("select url from calendars where id =?")
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

	stmt := "CREATE TABLE IF NOT EXISTS calendars(id INTEGER PRIMARY KEY, url TEXT, replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	stmt = "INSERT INTO calendars(url, replacementSummary) VALUES(?, ?);"
	_, err = db.Exec(stmt, params.Url, params.ReplacementSummary)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: return id of created record
}
