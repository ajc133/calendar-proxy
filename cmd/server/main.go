package main

import (
	"fmt"
	"os"

	"github.com/ajc133/calendarproxy/pkg/db"
	"github.com/ajc133/calendarproxy/pkg/handlers"
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

func main() {
	db.InitDB("data/calendars.db")

	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.GET("/:id", handlers.GetCalendarByID)
	router.POST("/", handlers.CreateCalendar)
	socket := fmt.Sprintf("%s:%s", httpInterface, httpPort)
	router.Run(socket)
}
