package model

import "time"

type UploadFileRequest struct {
	EntityType string `form:"entity_type" validate:"required"`
	EntityID   string `form:"entity_id" validate:"required"`
}

type FileUploadResponse struct {
	Key        string    `json:"key"`
	PublicURL  string    `json:"public_url"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type DeleteFileRequest struct {
	Key string `query:"key" validate:"required"`
}

type PresignedURLRequest struct {
	Key       string `query:"key" validate:"required"`
	Operation string `query:"operation" validate:"required,oneof=GET PUT"`
}

type PresignedURLResponse struct {
	URL       string `json:"url"`
	ExpiresAt int    `json:"expires_at"`
}
