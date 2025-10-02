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

func NewFileService(cfg config.Config) FileService {

	miniClient, _ := MinioClientConnection(cfg)

	return FileService{
		minio: miniClient,
		cfg:   cfg,
	}
}

func (s FileService) UploadImage(ctx context.Context, file *multipart.FileHeader, src multipart.File, fileName string) (dto.FileResponse, error) {
	bucketName := s.cfg.BucketName

	_, err := s.minio.PutObject(
		ctx,
		bucketName,
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		return dto.FileResponse{}, err
	}

	url, err := s.minio.PresignedGetObject(ctx, bucketName, fileName, time.Hour*24*7, nil)
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
