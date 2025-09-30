package handlers

import (
	"belimang/internal/services"
	"belimang/internal/utils"
	"net/http"
)

type FileHandler struct {
	service services.FileService
}

func NewFileHandler(service services.FileService) FileHandler {
	return FileHandler{
		service: service,
	}
}

func (h FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, handler, err := r.FormFile("file")

	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// size limits
	const (
		minFileSize = 10 * 1024       // 10 KB
		maxFileSize = 2 * 1024 * 1024 // 2 MB
	)

	if handler.Size < minFileSize {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file size is too small (minimum 10KB required)")
		return
	}

	if handler.Size > maxFileSize {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file size exceeds the maximum limit of 2MB")
		return
	}

	if !services.IsAllowedFileType(handler.Filename, handler.Header.Get("Content-Type")) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file type is not allowed")
		return
	}

	uploadedFile, err := h.service.UploadImage(ctx, handler, file, handler.Filename)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusAccepted, uploadedFile)
}
