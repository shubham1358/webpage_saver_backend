package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"webpage_saver/constants"
	"webpage_saver/constants/envKeys"
	"webpage_saver/firestoredb"
	"webpage_saver/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
)

func downloadPage(url string, browser *rod.Browser) bool {
	// Launch a browser with default options
	// Create a new page
	page := browser.MustPage(url)

	// Wait for the page to load completely
	page.WaitStable(time.Second * 2)

	// Get the Page CDP session

	// Use the Page.captureSnapshot method to save as MHTML
	snapshot, err := proto.PageCaptureSnapshot{}.Call(page)
	if err != nil {
		log.Printf("Failed to capture snapshot: %v", err)
		return false
	}
	pageUUID := uuid.New()
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, pageUUID.String()+".mhtml")
	htmlPath := filepath.Join(tempDir, pageUUID.String()+".html")
	objectPath := os.Getenv(string(envKeys.StoragePath)) + "/" + pageUUID.String() + ".html"
	// Save the MHTML content to a file
	err = os.WriteFile(filePath, []byte(snapshot.Data), 0644)
	if err != nil {
		log.Printf("Failed to save MHTML file: %v", err)
		return false
	}
	// Convert MHTML to HTML using the external executer
	cmd := exec.Command("./mhtml-to-html", filePath, "--output", htmlPath)

	// Run the command
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to convert MHTML to HTML: %v", err)
		return false
	}
	// Upload the HTML file to Google Cloud Storage
	err = storage.UploadFile(objectPath, htmlPath)
	if err != nil {
		log.Printf("Failed to upload file to Google Cloud Storage: %v", err)
		return false
	}
	// Insert the page data into the database
	WebSaver := constants.WebSaver{
		Url:      url,
		Path:     objectPath,
		Date:     time.Now(),
		DateOnly: time.Now().Truncate(24 * time.Hour),
	}

	err = firestoredb.AddPage(WebSaver)
	if err != nil {
		log.Printf("Failed to save page data to Firestore: %v", err)
		return false
	}

	log.Printf("Successfully saved the page %s as MHTML", pageUUID.String())
	return true
}

func SaveHandler(c *gin.Context) {
	var requestBody struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	url := requestBody.URL
	if url == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}

	// Download the page as MHTML
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	response := downloadPage(url, browser)
	if response {
		c.JSON(200, gin.H{"message": "Page saved successfully"})
	} else {
		c.JSON(400, gin.H{"error": "Failed to save page"})
	}
}

func GetPageHander(c *gin.Context) {
	url := c.Query("url")
	date := c.Query("date")
	if url == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}
	// db := GetDBInstance()
	// Query the database for the page
	// Use a pointer to scan the value
	pageData, parsedDate, err := firestoredb.GetWebPageByDate(url, date)
	if err != nil {
		c.JSON(404, gin.H{"error": "Page not found"})
		return
	}
	if pageData.Url == "" {
		c.JSON(404, gin.H{"error": "Page not found"})
		return
	}
	downloadedData, err := storage.DownloadFileIntoMemory(pageData.Path)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate URL"})
		return
	}
	c.JSON(200, gin.H{
		"html": string(downloadedData),
		"date": parsedDate,
	})
}

func GetDatesHandler(c *gin.Context) {
	date := c.Query("date")
	url := c.Query("url")
	if url == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}
	if date == "" {
		c.JSON(400, gin.H{"error": "date parameter is required"})
		return
	}
	// Query the database for the page
	dates, err := firestoredb.GetAvailableDatesByMonth(url, date)
	if err != nil {
		log.Printf("Failed to get available dates: %v", err)
		c.JSON(404, gin.H{"error": "Page not found"})
		return
	}
	c.JSON(200, gin.H{"dates": dates})
}
