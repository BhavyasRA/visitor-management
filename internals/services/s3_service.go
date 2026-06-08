package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"entry-system/internals/config"

	"github.com/minio/minio-go/v7"
)

type S3Service struct {
	client *minio.Client
}

func NewS3Service() *S3Service {
	return &S3Service{
		client: config.NewS3Client(),
	}
}

func (s *S3Service) UploadFile(file *multipart.FileHeader,folder string,) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	bucket := config.GetEnv(
		"S3_BUCKET",
		"visitor-documents",
	)

	exists, err := s.client.BucketExists(
		context.Background(),
		bucket,
	)
	if err != nil {
		return "", err
	}

	if !exists {
		err := s.client.MakeBucket(
			context.Background(),
			bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return "", err
		}
	}

	folder = strings.Trim(folder, "/")

	ext := filepath.Ext(file.Filename)

	fileKey := fmt.Sprintf(
		"%s/document_%d%s",
		folder,
		time.Now().UnixNano(),
		ext,
	)

	_, err = s.client.PutObject(
		context.Background(),
		bucket,
		fileKey,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}

	publicBaseURL := strings.TrimRight(
		config.GetEnv("S3_PUBLIC_URL", "http://localhost:9000"),
		"/",
	)

	documentURL := fmt.Sprintf(
		"%s/%s/%s",
		publicBaseURL,
		bucket,
		fileKey,
	)

	return documentURL, nil
}
