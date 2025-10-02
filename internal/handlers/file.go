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
	const (
		MinFileSize = 1024 * 10       // 10 KB
		MaxFileSize = 1024 * 1024 * 2 // 2 MB
	)

	if r.ContentLength > 0 && r.ContentLength > (MaxFileSize+1024) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file size is upper maximum file size")
		return
	}

	// Limit the maximum bytes that can be read from the body.1024*10 is the overhead tolerance
	r.Body = http.MaxBytesReader(w, r.Body, int64(MaxFileSize+1024*10))

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	if handler.Size < MinFileSize {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file size is too small (minimum 10KB required)")
		return
	}

	if handler.Size > MaxFileSize {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file size exceeds the maximum limit of 2MB")
		return
	}

	if !services.IsAllowedFileType(handler.Filename, handler.Header.Get("Content-Type")) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "file type is not allowed")
		return
	}

	uploadedFile, err := h.service.UploadImage(r.Context(), file, handler, handler.Filename)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, uploadedFile)
}
