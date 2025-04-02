package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
)

func downloadPage(url string, browser *rod.Browser) bool {
	// Launch a browser with default options
	// Create a new page
	db := GetDBInstance()
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
	filePath := "page_store/" + pageUUID.String() + ".mhtml"
	htmlPath := "page_store/" + pageUUID.String() + ".html"
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
	// Insert the page data into the database
	_, err = db.Exec("INSERT INTO webmap (id, url, path) VALUES (?, ?, ?)", pageUUID, url, htmlPath)
	if err != nil {
		log.Printf("Failed to insert data into database: %v", err)
		return false
	}

	log.Printf("Successfully saved the page %s as MHTML", pageUUID.String())
	return true
}

func saveEndpoint(c *gin.Context, browser *rod.Browser) {
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
	response := downloadPage(url, browser)
	if response {
		c.JSON(200, gin.H{"message": "Page saved successfully"})
	} else {
		c.JSON(400, gin.H{"error": "Failed to save page"})
	}
}

func getPageEndpoint(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}
	db := GetDBInstance()
	// Query the database for the page
	// Use a pointer to scan the value
	var path string
	err := db.QueryRow("SELECT path FROM webmap WHERE url = ? ORDER BY created_at DESC", url).Scan(&path)
	if err != nil {
		c.JSON(404, gin.H{"error": "Page not found"})
		return
	}

	// Serve the HTML file directly
	c.Header("Content-Type", "text/html")
	c.File(path)
}
