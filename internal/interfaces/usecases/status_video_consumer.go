package usecases

//go:generate mockgen -source=status_video_consumer.go -destination=mock/video_status_consumer.go
type StatusVideoConsumer interface {
	ConsumeMessages()
}
