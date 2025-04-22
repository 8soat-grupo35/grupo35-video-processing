package repository

//go:generate mockgen -source=s3.go -destination=mock/s3.go
type S3 interface {
	SetBucketName(bucketName string)
	UploadFile(key string, fileData []byte, contentType string) error
}
