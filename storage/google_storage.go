package storage

// Package storage provides functions to interact with Google Cloud Storage.
// It includes functions to upload files to a specified bucket and object name.
// It uses the Google Cloud Storage client library to handle the upload process.
// The uploadFile function takes a bucket name, object name, and file path as arguments.
// It creates a new client, opens the file, and uploads it to the specified bucket and object name.

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
	"webpage_saver/constants/envKeys"

	"cloud.google.com/go/storage"
)

func UploadFile(objectName, filePath string) error {
	bucketName := os.Getenv(string(envKeys.BucketName)) // Replace with your bucket name
	ctx := context.Background()

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Println("File uploaded.")
	return nil
}

func DownloadFileIntoMemory(object string) ([]byte, error) {
	bucket := os.Getenv(string(envKeys.BucketName)) // Replace with your bucket name
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	return data, nil
}

func GenerateSignedURL(object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()
	url, err := storage.SignedURL(os.Getenv(string(envKeys.BucketName)), object, &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv(string(envKeys.GCPAccessKey)),
		PrivateKey:     []byte(os.Getenv(string(envKeys.GCPPrivateKey))),
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
