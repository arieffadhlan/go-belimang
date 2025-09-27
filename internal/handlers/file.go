package handlers

import "belimang/internal/services"

type FileHandler struct {
	service services.FileService
}

func NewFileHandler(service services.FileService) FileHandler {
	return FileHandler{
		service: service,
	}
}
