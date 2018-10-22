package middleware

func CutMessage(chunkSize int) func(message string) string {
	return func(message string) string {
		if len(message) > chunkSize {
			return message[:chunkSize]
		}
		return message
	}
}
