package services

import (
	"belimang/internal/config"
	"bytes"
	"image"
	"image/jpeg"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/image/draw"
)

func MinioClientConnection(cfg config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func CompressThumbnail(fileHeader *multipart.FileHeader, img image.Image, maxKB int) ([]byte, string, error) {
	var buf bytes.Buffer
	ext := strings.ToLower(fileHeader.Filename)

	// --- Resize to thumbnail (height = 256, width scaled proportionally) ---
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	newH := 256
	newW := w * newH / h

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	thumb := dst

	// --- If JPEG/JPG ---
	if strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".jpeg") {
		quality := 85
		for quality > 20 {
			buf.Reset()
			_ = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: quality})
			if buf.Len() <= maxKB*1024 {
				break
			}
			quality -= 5
		}
		return buf.Bytes(), "image/jpeg", nil
	}

	return buf.Bytes(), "image/png", nil
}

func IsAllowedFileType(fileName, fileType string) bool {
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
	}

	if !allowedMimeTypes[fileType] {
		return false
	}

	if fileType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(fileName))
		if !allowedExtensions[ext] {
			return false
		}
	}

	return true
}
