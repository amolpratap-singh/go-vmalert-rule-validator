package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/amolpratap-singh/vmalert-rule-validator/minioClient"
)

var (
	genLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// Custom errors
	ErrUploadFile   = errors.New("file upload error")
	ErrDownloadFile = errors.New("file download error")
	ErrListFile     = errors.New("list file error")
)

func main() {
	genLogger.Println("Go VM alert rule Validator")
	genLogger.Println("vmalert-init-container")

	minioClient, err := minioClient.NewMinioClient()

	if err != nil {
		log.Fatalf("Failed to create Minio client: %v", err)
	}

	files, err := minioClient.ListFiles()
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	log.Println("Files in bucket:", strings.Join(files, ", "))
}
