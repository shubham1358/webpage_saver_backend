package main

import (
	"log"
	"os"
)

// GetDBInstance returns a singleton DuckDB connection with a persistent database file
func InitDir() {
	// Check if the directory exists
	if _, err := os.Stat("database"); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir("database", 0755)
		if err != nil {
			log.Fatal("Failed to create directory:", err)
		}
		log.Println("Directory created successfully")
	} else {
		log.Println("Directory already exists")
	}
	if _, err := os.Stat("page_store"); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir("page_store", 0755)
		if err != nil {
			log.Fatal("Failed to create page_store:", err)
		}
		log.Println("page_store created successfully")
	} else {
		log.Println("page_store already exists")
	}
	// Check if mhtml-html exists and has the executable permission
	fileInfo, err := os.Stat("mhtml-to-html")
	if os.IsNotExist(err) {
		log.Fatal("mhtml-html does not exist")
	} else if err != nil {
		log.Fatal("Failed to check mhtml-html:", err)
	} else {
		// Check if the file has executable permission
		if fileInfo.Mode()&0111 == 0 {
			// Add executable permission
			err := os.Chmod("mhtml-to-html", fileInfo.Mode()|0755)
			if err != nil {
				log.Fatal("Failed to add executable permission to mhtml-html:", err)
			}
			log.Println("Executable permission added to mhtml-html")
		} else {
			log.Println("mhtml-html already has executable permission")
		}
	}
}
