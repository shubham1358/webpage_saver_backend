package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")
	// Creates a router without any middleware by default
	if os.Getenv("ENV") == "DEVELOPMENT" {
		// Set Gin to development mode
		fmt.Println("Setting Gin to development mode")
		gin.SetMode(gin.DebugMode)
		godotenv.Load(".env.local")
	} else {
		// Set Gin to production mode
		fmt.Println("Setting Gin to production mode")
	}
	port := os.Getenv("PORT")

	// Load environment variables from .env file
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	init_db()
	router := gin.New()
	// Set up a custom logger to write to the log directory
	if os.Getenv("ENV") == "DEVELOPMENT" {
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
	} else {
		gin.DefaultWriter = os.Stdout
	}
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CLIENT_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Per route middleware, you can add as many as you desire.
	router.POST("/save", SaveHandler)
	router.GET("/page", GetPageHander)
	router.GET("/dates", GetDatesHandler)
	router.GET("/", func(c *gin.Context) {
		log.Println("Get Called")
		c.String(200, "Runnning at port %s", port)
	})
	// Listen and serve on 0.0.0.0:8080
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}
	router.Run(":" + port)
}
