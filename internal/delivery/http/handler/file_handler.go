package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model/mapper"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
)

type FileHandler struct {
	useCase usecase.FileUploadUseCase
}

func NewFileHandler(useCase usecase.FileUploadUseCase) *FileHandler {
	return &FileHandler{
		useCase: useCase,
	}
}

func (h *FileHandler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/files")
	group.Post("/upload", h.Upload)
	group.Delete("/", h.Delete)
	group.Get("/presigned", h.GetPresignedURL)
}

func (h *FileHandler) Upload(c fiber.Ctx) error {
	var req model.UploadFileRequest
	if err := c.Bind().Form(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid form data", err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "No file provided or invalid multipart form", err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open file", err.Error())
	}
	defer src.Close()

	uploadInput := &entity.UploadInput{
		EntityType:  req.EntityType,
		EntityID:    req.EntityID,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		File:        src,
		Size:        file.Size,
	}

	result, err := h.useCase.Upload(c.Context(), uploadInput)
	if err != nil {
		return err
	}

	responseData := mapper.UploadResultToResponse(result)

	return response.Send(c, fiber.StatusCreated, "File uploaded successfully", responseData)
}

func (h *FileHandler) Delete(c fiber.Ctx) error {
	var req model.DeleteFileRequest
	if err := c.Bind().Query(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	err := h.useCase.Delete(c.Context(), req.Key)
	if err != nil {
		return err
	}

	return response.Send[any](c, fiber.StatusOK, "File deleted successfully", nil)
}

func (h *FileHandler) GetPresignedURL(c fiber.Ctx) error {
	var req model.PresignedURLRequest
	if err := c.Bind().Query(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	if req.Operation == "" {
		req.Operation = "GET"
	}

	url, err := h.useCase.GetPresignedURL(c.Context(), req.Key, req.Operation)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Presigned URL generated successfully", model.PresignedURLResponse{
		URL:       url,
		ExpiresAt: 0, // TODO: Calculate actual expiry if needed in response
	})
}
