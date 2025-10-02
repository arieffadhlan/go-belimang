package services

import (
	"belimang/internal/config"
	"belimang/internal/dto"
	"context"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileService struct {
	minio *minio.Client
	cfg   config.Config
}

func NewFileService(client *minio.Client, config config.Config) FileService {
	return FileService{
		minio: client,
		cfg:   config,
	}
}

func (s FileService) UploadImage(ctx context.Context, src multipart.File, file *multipart.FileHeader, objectName string) (dto.FileResponse, error) {
	bucketName := s.cfg.BucketName

	_, err := s.minio.PutObject(
		ctx,
		bucketName,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		 return dto.FileResponse{}, err
	}

	url, err := s.minio.PresignedGetObject(ctx, bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		 return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		Message: "File uploaded sucessfully",
		Data: dto.UploadFileData{
			ImageUrl: url.String(),
		},
	}, nil
}
