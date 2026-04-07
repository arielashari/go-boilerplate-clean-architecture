package usecase_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase/mocks"
)

func setupFileUploadUseCase(t *testing.T) (usecase.FileUploadUseCase, *mocks.MockFileStorage) {
	mockStorage := new(mocks.MockFileStorage)
	uc := usecase.NewFileUploadUseCase(mockStorage)
	t.Cleanup(func() { mockStorage.AssertExpectations(t) })
	return uc, mockStorage
}

func TestFileUploadUseCase_Upload(t *testing.T) {
	t.Run("success with valid file", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		fileContent := []byte("test file content")
		input := &entity.UploadInput{
			EntityType:  "user",
			EntityID:    "user-123",
			FileName:    "test.jpg",
			ContentType: "image/jpeg",
			File:        bytes.NewReader(fileContent),
			Size:        int64(len(fileContent)),
		}

		expectedResult := &entity.UploadResult{
			Key:        "uploads/user/user-123/uuid.jpg",
			PublicURL:  "https://bucket.s3.amazonaws.com/uploads/user/user-123/uuid.jpg",
			Size:       int64(len(fileContent)),
			UploadedAt: time.Now(),
		}

		mockStorage.On("Upload", mock.Anything, input).Return(expectedResult, nil)

		result, err := uc.Upload(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult.Key, result.Key)
		assert.Equal(t, expectedResult.Size, result.Size)
	})

	t.Run("file too large", func(t *testing.T) {
		uc, _ := setupFileUploadUseCase(t)

		input := &entity.UploadInput{
			EntityType:  "user",
			EntityID:    "user-123",
			FileName:    "large.jpg",
			ContentType: "image/jpeg",
			File:        bytes.NewReader([]byte("")),
			Size:        11 * 1024 * 1024, // 11MB, exceeds 10MB limit
		}

		result, err := uc.Upload(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrFileTooLarge)
	})

	t.Run("invalid MIME type", func(t *testing.T) {
		uc, _ := setupFileUploadUseCase(t)

		input := &entity.UploadInput{
			EntityType:  "user",
			EntityID:    "user-123",
			FileName:    "test.exe",
			ContentType: "application/x-msdownload",
			File:        bytes.NewReader([]byte("test")),
			Size:        4,
		}

		result, err := uc.Upload(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrInvalidFileType)
	})

	t.Run("storage not initialized", func(t *testing.T) {
		uc := usecase.NewFileUploadUseCase(nil) // nil storage

		input := &entity.UploadInput{
			EntityType:  "user",
			EntityID:    "user-123",
			FileName:    "test.jpg",
			ContentType: "image/jpeg",
			File:        bytes.NewReader([]byte("test")),
			Size:        4,
		}

		result, err := uc.Upload(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrStorageUnavailable)
	})

	t.Run("all allowed MIME types", func(t *testing.T) {
		allowedTypes := []string{
			"image/jpeg",
			"image/png",
			"image/webp",
			"application/pdf",
		}

		for _, mimeType := range allowedTypes {
			t.Run(mimeType, func(t *testing.T) {
				uc, mockStorage := setupFileUploadUseCase(t)

				fileContent := []byte("test content")
				input := &entity.UploadInput{
					EntityType:  "user",
					EntityID:    "user-123",
					FileName:    "test.file",
					ContentType: mimeType,
					File:        bytes.NewReader(fileContent),
					Size:        int64(len(fileContent)),
				}

				expectedResult := &entity.UploadResult{
					Key:        "uploads/user/user-123/uuid.file",
					PublicURL:  "https://bucket.s3.amazonaws.com/uploads/user/user-123/uuid.file",
					Size:       int64(len(fileContent)),
					UploadedAt: time.Now(),
				}

				mockStorage.On("Upload", mock.Anything, input).Return(expectedResult, nil)

				result, err := uc.Upload(context.Background(), input)

				assert.NoError(t, err)
				assert.NotNil(t, result)
			})
		}
	})
}

func TestFileUploadUseCase_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		mockStorage.On("Delete", mock.Anything, key).Return(nil)

		err := uc.Delete(context.Background(), key)

		assert.NoError(t, err)
	})

	t.Run("storage not initialized", func(t *testing.T) {
		uc := usecase.NewFileUploadUseCase(nil) // nil storage

		err := uc.Delete(context.Background(), "uploads/user/user-123/uuid.jpg")

		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrStorageUnavailable)
	})

	t.Run("storage error", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		mockStorage.On("Delete", mock.Anything, key).Return(entity.ErrInternal)

		err := uc.Delete(context.Background(), key)

		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrInternal)
	})
}

func TestFileUploadUseCase_GetPresignedURL(t *testing.T) {
	t.Run("success GET operation", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		expectedURL := "https://presigned-url.example.com"

		mockStorage.On("GetPresignedURL", mock.Anything, key, "GET").Return(expectedURL, nil)

		url, err := uc.GetPresignedURL(context.Background(), key, "GET")

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
	})

	t.Run("success PUT operation", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		expectedURL := "https://presigned-url.example.com"

		mockStorage.On("GetPresignedURL", mock.Anything, key, "PUT").Return(expectedURL, nil)

		url, err := uc.GetPresignedURL(context.Background(), key, "PUT")

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
	})

	t.Run("default to GET when operation empty", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		expectedURL := "https://presigned-url.example.com"

		mockStorage.On("GetPresignedURL", mock.Anything, key, "GET").Return(expectedURL, nil)

		url, err := uc.GetPresignedURL(context.Background(), key, "")

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
	})

	t.Run("storage not initialized", func(t *testing.T) {
		uc := usecase.NewFileUploadUseCase(nil) // nil storage

		url, err := uc.GetPresignedURL(context.Background(), "uploads/user/user-123/uuid.jpg", "GET")

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.ErrorIs(t, err, entity.ErrStorageUnavailable)
	})

	t.Run("storage error", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		key := "uploads/user/user-123/uuid.jpg"
		mockStorage.On("GetPresignedURL", mock.Anything, key, "GET").Return("", entity.ErrInternal)

		url, err := uc.GetPresignedURL(context.Background(), key, "GET")

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.ErrorIs(t, err, entity.ErrInternal)
	})

	t.Run("invalid operation", func(t *testing.T) {
		uc, mockStorage := setupFileUploadUseCase(t)

		url, err := uc.GetPresignedURL(context.Background(), "key", "INVALID")

		assert.Error(t, err)
		assert.Empty(t, url)
		// Should not have called storage.GetPresignedURL due to validation
		mockStorage.AssertNotCalled(t, "GetPresignedURL")
	})
}
