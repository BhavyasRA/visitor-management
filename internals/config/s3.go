package config

import (
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3Client *minio.Client

func NewS3Client() *minio.Client {
	if S3Client != nil {
		return S3Client
	}

	endpoint := GetEnv("S3_ENDPOINT", "localhost:9000")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	accessKey := GetEnv("S3_ACCESS_KEY", "bhavya")
	secretKey := GetEnv("S3_SECRET_KEY", "Nanibhavya*979")
	useSSL := GetEnv("S3_USE_SSL", "true") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal("failed to connect MinIO: ", err)
	}

	S3Client = client
	return S3Client
}