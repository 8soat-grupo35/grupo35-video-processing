package entities

const (
	VideoStatusInProcessing = "IN_PROCESSING"
	VideoStatusSuccess      = "SUCCESS"
	VideoStatusError        = "ERROR"
)

type Video struct {
	UserID    string `dynamo:"user_id,hash"`
	Path      string `dynamo:"path,range"`
	Status    string `dynamo:"status"`
	CreatedAt string `dynamo:"created_at"`
	UpdatedAt string `dynamo:"updated_at"`
}
