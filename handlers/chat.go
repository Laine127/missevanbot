package handlers

import (
	"fmt"

	"missevan-fm/models"
)

// Chat 简单回复
func Chat(username string) string {
	return fmt.Sprintf("@%s %s", username, models.ChatString())
}
