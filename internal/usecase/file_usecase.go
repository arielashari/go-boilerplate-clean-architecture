package usecase

import (
	"context"
	"fmt"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"
)

const (
	MaxFileSize = 10 * 1024 * 1024
)

var (
	AllowedMIMETypes = map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
		"application/pdf": true,
	}
)

type FileUploadUseCase interface {
	Upload(ctx context.Context, input *entity.UploadInput) (*entity.UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key, operation string) (string, error)
}

type fileUploadUseCase struct {
	storage entity.FileStorage
}

func NewFileUploadUseCase(storage entity.FileStorage) FileUploadUseCase {
	return &fileUploadUseCase{
		storage: storage,
	}
}

func (u *fileUploadUseCase) Upload(ctx context.Context, input *entity.UploadInput) (*entity.UploadResult, error) {
	if u.storage == nil {
		return nil, entity.ErrStorageUnavailable.WithOperation("FileUploadUseCase.Upload")
	}

	if input.Size > MaxFileSize {
		return nil, entity.ErrFileTooLarge.
			WithInternal(fmt.Errorf("file size %d exceeds maximum %d bytes", input.Size, MaxFileSize)).
			WithOperation("FileUploadUseCase.Upload.ValidateSize")
	}

	if !AllowedMIMETypes[input.ContentType] {
		return nil, entity.ErrInvalidFileType.
			WithInternal(fmt.Errorf("MIME type %q not supported", input.ContentType)).
			WithOperation("FileUploadUseCase.Upload.ValidateMIME")
	}

	result, err := u.storage.Upload(ctx, input)
	if err != nil {
		if appErr, ok := apperror.As(err); ok {
			return nil, appErr.WithOperation("FileUploadUseCase.Upload")
		}
		return nil, apperror.New(apperror.CodeInternal, "file upload failed").
			WithInternal(err).
			WithOperation("FileUploadUseCase.Upload")
	}

	return result, nil
}

func (u *fileUploadUseCase) Delete(ctx context.Context, key string) error {
	if u.storage == nil {
		return entity.ErrStorageUnavailable.WithOperation("FileUploadUseCase.Delete")
	}

	err := u.storage.Delete(ctx, key)
	if err != nil {
		if appErr, ok := apperror.As(err); ok {
			return appErr.WithOperation("FileUploadUseCase.Delete")
		}
		return apperror.New(apperror.CodeInternal, "file deletion failed").
			WithInternal(err).
			WithOperation("FileUploadUseCase.Delete")
	}

	return nil
}

func (u *fileUploadUseCase) GetPresignedURL(ctx context.Context, key, operation string) (string, error) {
	if u.storage == nil {
		return "", entity.ErrStorageUnavailable.WithOperation("FileUploadUseCase.GetPresignedURL")
	}

	if operation == "" {
		operation = "GET"
	}

	if operation != "GET" && operation != "PUT" {
		return "", apperror.New(apperror.CodeValidation, "operation must be GET or PUT").
			WithOperation("FileUploadUseCase.GetPresignedURL.ValidateOperation")
	}

	url, err := u.storage.GetPresignedURL(ctx, key, operation)
	if err != nil {
		if appErr, ok := apperror.As(err); ok {
			return "", appErr.WithOperation("FileUploadUseCase.GetPresignedURL")
		}
		return "", apperror.New(apperror.CodeInternal, "failed to generate presigned URL").
			WithInternal(err).
			WithOperation("FileUploadUseCase.GetPresignedURL")
	}

	return url, nil
}
