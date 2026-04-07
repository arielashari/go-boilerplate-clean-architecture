package storage

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3client "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type S3Storage struct {
	client *s3.Client
	cfg    *configs.S3Config
}

func NewS3Storage(cfg *configs.S3Config) (*S3Storage, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"",
	))

	config, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithCredentialsProvider(creds),
		awsconfig.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3client.NewFromConfig(config)
	slog.Info("S3 client initialized successfully", "region", cfg.Region, "bucket", cfg.Bucket)

	return &S3Storage{
		client: client,
		cfg:    cfg,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, input *entity.UploadInput) (*entity.UploadResult, error) {
	if input == nil {
		return nil, apperror.New(
			apperror.CodeValidation,
			"upload input cannot be nil",
		).WithOperation("S3Storage.Upload")
	}

	ext := filepath.Ext(input.FileName)
	fileID := uuid.New().String()
	key := fmt.Sprintf("uploads/%s/%s/%s%s", input.EntityType, input.EntityID, fileID, ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(key),
		Body:        input.File,
		ContentType: aws.String(input.ContentType),
		Metadata: map[string]string{
			"entity-type": input.EntityType,
			"entity-id":   input.EntityID,
		},
	})
	if err != nil {
		return nil, apperror.New(
			apperror.CodeInternal,
			"failed to upload file to storage",
		).WithInternal(err).WithOperation("S3Storage.Upload")
	}

	result := &entity.UploadResult{
		Key:        key,
		PublicURL:  s.GetPublicURL(key),
		Size:       input.Size,
		UploadedAt: time.Now(),
	}

	return result, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return apperror.New(
			apperror.CodeValidation,
			"file key cannot be empty",
		).WithOperation("S3Storage.Delete")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return apperror.New(
			apperror.CodeInternal,
			"failed to delete file from storage",
		).WithInternal(err).WithOperation("S3Storage.Delete")
	}

	return nil
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, key string, operation string) (string, error) {
	if key == "" {
		return "", apperror.New(
			apperror.CodeValidation,
			"file key cannot be empty",
		).WithOperation("S3Storage.GetPresignedURL")
	}

	if operation == "" {
		operation = "GET"
	}

	if operation != "GET" && operation != "PUT" {
		return "", apperror.New(
			apperror.CodeValidation,
			"operation must be GET or PUT",
		).WithOperation("S3Storage.GetPresignedURL")
	}

	publicURL := s.GetPublicURL(key)
	if publicURL == "" {
		return "", apperror.New(
			apperror.CodeInternal,
			"failed to generate URL - base URL not configured",
		).WithOperation("S3Storage.GetPresignedURL")
	}

	// TODO: Implement full presigning with SigV4 signature and expiration
	return publicURL, nil
}

func (s *S3Storage) GetPublicURL(key string) string {
	if s.cfg.BaseURL == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, key)
}

func (s *S3Storage) MakePublic(ctx context.Context, key string) error {
	_, err := s.client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return apperror.New(
			apperror.CodeInternal,
			"failed to make file public",
		).WithInternal(err).WithOperation("S3Storage.MakePublic")
	}
	return nil
}
