package main

import (
	"log"
	"os"
	"webpage_saver/firestoredb"
)

func init_db() {
	// Open the database
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is not set")
	}
	firestoredb.Init(projectID)
}
