package constants

import (
	"time"
)

type WebSaver struct {
	Url      string    `firestore:"url"`
	Date     time.Time `firestore:"date"`
	Path     string    `firestore:"path"`
	DateOnly time.Time `firestore:"date_only"`
}
