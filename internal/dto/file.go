package dto

type (
	UploadFileData struct {
		ImageUrl string `json:"imageUrl"`
	}
	FileResponse struct {
		Message string         `json:"message"`
		Data    UploadFileData `json:"data"`
	}
)
