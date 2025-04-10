package firestoredb

import (
	"context"
	"log"
	"time"
	"webpage_saver/constants"

	"cloud.google.com/go/firestore"
)

var Client *firestore.Client

func Init(projectID string) {
	ctx := context.Background()
	var err error
	Client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
}

func AddPage(data constants.WebSaver) error {
	_, _, err := Client.Collection("saved_pages").Add(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

// GetPageByDate retrieves a page by its URL and date string
// date should be in ISO format
func GetWebPageByDate(url string, date string) (constants.WebSaver, time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		parsedDate = time.Now().Truncate(24 * time.Hour)
	}

	iter := Client.Collection("saved_pages").
		Where("url", "==", url).
		Where("date_only", "==", parsedDate).
		Limit(1).
		Documents(context.Background())

	doc, err := iter.Next()
	if err != nil {
		return constants.WebSaver{}, parsedDate, err
	}

	var result constants.WebSaver
	if err := doc.DataTo(&result); err != nil {
		return constants.WebSaver{}, parsedDate, err
	}

	return result, parsedDate, nil
}

func GetAvailableDatesByMonth(url string, date string) ([]time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	iter := Client.Collection("saved_pages").
		Where("url", "==", url).
		Where("date_only", ">=", parsedDate).
		Where("date_only", "<=", parsedDate.AddDate(0, 1, 0)).
		Select("date_only").
		Documents(context.Background())

	var dates []time.Time
	for {
		doc, err := iter.Next()
		if err != nil {
			if err.Error() == "iterator done" {
				break
			}
			log.Printf("Failed to get next document: %v", err)
			if len(dates) > 0 {
				return dates, nil
			}
			return nil, err
		}

		var data struct {
			DateOnly time.Time `firestore:"date_only"`
		}
		if err := doc.DataTo(&data); err != nil {
			log.Printf("Failed to convert document data: %v", err)
			return nil, err
		}

		dates = append(dates, data.DateOnly)
	}

	return dates, nil
}
