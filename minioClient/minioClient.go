package minioClient

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	genLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// Custom errors
	ErrUploadFile   = errors.New("file upload error")
	ErrDownloadFile = errors.New("file download error")
	ErrListFile     = errors.New("list file error")
)

type MinioClient struct {
	client      *minio.Client
	minioHost   string
	minioPort   string
	minioBucket string
	minioUser   string
	minioPass   string
	sslCaCert   string
}

func NewMinioClient() (*MinioClient, error) {
	host := getEnv("MINIO_HOST", "data-object-stroage-minio")
	port := getEnv("MINIO_POR", "9000")
	bucket := getEnv("MINIO_FILE_BUCKET", "file_storage_bucket")
	user := getEnv("MINIO_ACCESS_KEY", "data-object-storage-minio-creds")
	pass := getEnv("MINIO_SECRET_KEY", "data-object-storage-minio-creds")
	sslCaCert := getEnv("MINIO_CA_CERT", "")

	client, err := minio.New(host+":"+port, &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: true,
	})

	if err != nil {
		return nil, err
	}

	minioClient := &MinioClient{
		client:      client,
		minioHost:   host,
		minioPort:   port,
		minioBucket: bucket,
		minioUser:   user,
		minioPass:   pass,
		sslCaCert:   sslCaCert,
	}

	return minioClient, nil
}

func (mc *MinioClient) ListFiles() ([]string, error) {
	ctx := context.Background()
	if exists, err := mc.checkBucketExists(mc.minioBucket); err != nil || !exists {
		genLogger.Printf("Bucket name: %s does not exist, cannot list files\n", mc.minioBucket)
		return nil, ErrListFile
	}
	var fileNames []string
	for object := range mc.client.ListObjects(ctx, mc.minioBucket, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			genLogger.Printf("Error in listing files from Minio - Reason: %v\n", object.Err)
			return nil, ErrListFile
		}
		fileNames = append(fileNames, object.Key)
	}
	return fileNames, nil
}

func (mc *MinioClient) checkBucketExists(bucketName string) (bool, error) {
	ctx := context.Background()
	exists, err := mc.client.BucketExists(ctx, bucketName)
	if err != nil {
		genLogger.Printf("Bucket name '%s' does not exist. Error: %v\n", bucketName, err)
		return false, err
	}
	return exists, nil
}

func (mc *MinioClient) GetFile(fileName, filePath string) error {
	ctx := context.Background()
	err := mc.client.FGetObject(ctx, mc.minioBucket, fileName, filePath, minio.GetObjectOptions{})
	if err != nil {
		genLogger.Printf("Error in getting the file %s from Minio - Reason: %v\n", fileName, err)
		return ErrDownloadFile
	}
	return nil
}

func (mc *MinioClient) IsFileExists(fileName string) (bool, error) {
	ctx := context.Background()
	for object := range mc.client.ListObjects(ctx, mc.minioBucket, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			genLogger.Printf("Error: %v\n", object.Err)
			return false, object.Err
		}
		if object.Key == fileName {
			return true, nil
		}
	}
	return false, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
