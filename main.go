package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	httpInterface = os.Getenv("HTTP_INTERFACE")
	httpPort      = os.Getenv("HTTP_PORT")
)

func init() {
	if httpInterface == "" {
		httpInterface = "0.0.0.0"
	}
	if httpPort == "" {
		httpPort = "8080"
	}
}

const ContentType string = "Content-Type"
const CalendarContent string = "text/calendar; charset=utf-8"

func GetCalendarStateless(c *gin.Context) {
	url := c.Query("url")
	fmt.Printf("url: %s\n", url)
	replacement := c.Query("replacementSummary")

	newCal, err := FetchAndTransformCalendar(url, replacement)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when fetching and parsing the given URL"})
	}
	c.Header(ContentType, CalendarContent)
	c.String(http.StatusOK, newCal)
}

func main() {

	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.GET("/cal", GetCalendarStateless)
	socket := fmt.Sprintf("%s:%s", httpInterface, httpPort)
	log.Printf("Starting server on %s\n", socket)
	router.Run(socket)
}
