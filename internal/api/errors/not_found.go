package custom_errors

type NotFoundError struct {
	Message string
}

func (b *NotFoundError) Error() string {
	return b.Message
}
