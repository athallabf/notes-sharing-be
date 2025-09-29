// src/service/s3_service.go

package service

import (
	"app/src/config"
	"app/src/utils"
	"context"
	"mime/multipart"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type S3Service interface {
	UploadFile(file multipart.File, objectKey string) error
	GeneratePresignedURL(objectKey string) (string, error)
}

type s3Service struct {
	Log           *logrus.Logger
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	BucketName    string
}

func NewS3Service() S3Service {
	if config.S3AccessKey == "" || config.S3SecretKey == "" || config.S3BucketName == "" || config.S3Region == "" {
		utils.Log.Fatalf("S3 is not configured. Please set S3_ACCESS_KEY, S3_SECRET_KEY, S3_BUCKET_NAME, and S3_REGION.")
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(config.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.S3AccessKey,
			config.S3SecretKey,
			"",
		)),
	)
	if err != nil {
		utils.Log.Fatalf("Failed to load S3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &s3Service{
		Log:           utils.Log,
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
		BucketName:    config.S3BucketName,
	}
}

func (s *s3Service) UploadFile(file multipart.File, objectKey string) error {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.BucketName,
		Key:    &objectKey,
		Body:   file,
	})
	return err
}

func (s *s3Service) GeneratePresignedURL(objectKey string) (string, error) {
	request, err := s.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.BucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return "", err
	}
	return request.URL, nil
}
