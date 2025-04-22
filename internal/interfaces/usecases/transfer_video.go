package usecases

//go:generate mockgen -source=transfer_video.go -destination=mock/transfer_video.go
type TransferVideoUseCase interface {
	UploadVideo(key string, videoData []byte) (err error)
}
