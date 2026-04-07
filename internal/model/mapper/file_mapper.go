package mapper

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
)

func UploadResultToResponse(result *entity.UploadResult) *model.FileUploadResponse {
	if result == nil {
		return nil
	}
	return &model.FileUploadResponse{
		Key:        result.Key,
		PublicURL:  result.PublicURL,
		Size:       result.Size,
		UploadedAt: result.UploadedAt,
	}
}
