package handlers

import (
	"fmt"

	"missevanbot/modules"
)

// Chat returns a simple response.
func Chat(name string) string {
	return fmt.Sprintf("@%s %s", name, modules.Word(modules.WordChat))
}
