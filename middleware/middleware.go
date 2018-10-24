package middleware

import (
	"encoding/json"
	"fmt"
)

type SentryMessage struct {
	Project     string
	ProjectName string `json:"project_name"`
	ProjectSlug string `json:"project_slug"`
	Url         string
	Message     string
}

func CutMessage(chunkSize int) func(message string) string {
	return func(message string) string {
		if len(message) > chunkSize {
			return message[:chunkSize]
		}
		return message
	}
}

func SentryFormatter(message string) string {
	var sentryMessage SentryMessage
	json.Unmarshal([]byte(message), &sentryMessage)
	return fmt.Sprintf(
		"Project: %v\nUrl: %v\nMessage: \n%v", sentryMessage.Project, sentryMessage.Url, sentryMessage.Message)
}
