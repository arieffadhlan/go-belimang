package services

import (
	"belimang/internal/config"
)

type FileService struct {
	cfg config.Config
}

func NewFileService(cfg config.Config) FileService {
	return FileService{
		cfg: cfg,
	}
}
