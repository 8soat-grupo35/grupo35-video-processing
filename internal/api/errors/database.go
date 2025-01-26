package custom_errors

type DatabaseError struct {
	Message string
}

func (b *DatabaseError) Error() string {
	return b.Message
}
