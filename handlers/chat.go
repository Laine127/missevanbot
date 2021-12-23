package handlers

import (
	"fmt"

	"missevanbot/modules"
)

// Chat returns a simple response.
func Chat(username string) string {
	return fmt.Sprintf("@%s %s", username, modules.Word(modules.WordChat))
}
