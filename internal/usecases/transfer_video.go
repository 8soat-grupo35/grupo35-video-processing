package usecases

import (
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases"
)

type TransferVideoUseCase struct {
	S3 repository.S3
}

func NewTransferVideo(s3 repository.S3) usecases.TransferVideoUseCase {
	return TransferVideoUseCase{
		S3: s3,
	}
}

func (t TransferVideoUseCase) UploadVideo(key string, videoData []byte) (err error) {
	fmt.Println("Uploading video to S3")
	t.S3.SetBucketName("grupo35-video-uploaded")
	err = t.S3.UploadFile(key, videoData)

	if err != nil {
		return fmt.Errorf("upload video failed: %w", err)
	}

	fmt.Println("Video uploaded successfully to S3")

	return nil
}
