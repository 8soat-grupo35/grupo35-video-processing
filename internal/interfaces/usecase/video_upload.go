package usecase

import "mime/multipart"

type VideoUploaderUseCase interface {
	UploadVideos(videoFiles []*multipart.FileHeader) error
}
