package usecases

import (
	"github.com/8soat-grupo35/fastfood-order-production/internal/interfaces/usecase"
	"mime/multipart"
)

type VideoUploaderUseCase struct {
}

func NewVideoUploaderUseCase() usecase.VideoUploaderUseCase {
	return &VideoUploaderUseCase{}
}

func (v *VideoUploaderUseCase) UploadVideos(videoFiles []*multipart.FileHeader) error {
	return nil
}
