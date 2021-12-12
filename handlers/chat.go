package handlers

import (
	"fmt"

	"missevan-fm/models"
)

// Chat return a simple response.
func Chat(username string) string {
	return fmt.Sprintf("@%s %s", username, models.ChatString())
}
