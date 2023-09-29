package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.GET("/:id", GetCalendarByID)
	router.POST("/", CreateCalendar)
	router.Run("localhost:8080")
}
