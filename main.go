package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
)

func main() {
	// Creates a router without any middleware by default
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router := gin.New()
	// Set up a custom logger to write to the log directory
	logFile, err := os.OpenFile("log/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// If log directory doesn't exist, create it
		os.MkdirAll("log", 0755)
		logFile, err = os.Create("log/gin.log")
		if err != nil {
			log.Fatal("Could not create log file: ", err)
		}
	}
	log.SetOutput(logFile)
	gin.DefaultWriter = logFile
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.ClientIP,
			param.Method,
			param.StatusCode,
			param.Path,
		)
	}))
	InitDir()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	router.POST("/save", func(c *gin.Context) {
		saveEndpoint(c, browser) // Pass nil for rod.Browser or initialize it properly
	})
	router.GET("/page", getPageEndpoint)
	router.GET("/", func(c *gin.Context) {
		log.Println("Get Called")
		c.String(200, "Runnning at port 8080")
	})
	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
