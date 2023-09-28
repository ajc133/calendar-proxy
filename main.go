package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type CalendarParams struct {
	Url                string `form:"url" json:"url" binding:"required"`
	ReplacementSummary string `form:"replacementSummary" json:"replacementSummary" binding:"required"`
}

func main() {
	CreateDB()
	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.GET("/:id", GetCalendarByID)
	router.POST("/", CreateCalendar)
	router.Run("localhost:8080")
}

func CreateCalendar(c *gin.Context) {
	json := CalendarParams{}
	err := c.Bind(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	WriteRecord(json)
	c.JSON(http.StatusOK, gin.H{"status": "successfully added entry"})
}

func GetCalendarByID(c *gin.Context) {
	id := c.Param("id")
	url, err := ReadRecord(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "url": url})
}
func CreateDB() error {
	db, err := sql.Open("sqlite3", "./calendars.db")
	stmt := "CREATE TABLE IF NOT EXISTS calendars(id INTEGER PRIMARY KEY, url TEXT, replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ReadRecord(id string) (string, error) {
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

func WriteRecord(params CalendarParams) {
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
