package handlers

import (
	"fmt"

	"missevanbot/models"
	"missevanbot/modules"
)

// keyEmotional handle the keyword `emo`.
func keyEmotional(output chan<- string, user models.FmUser) {
	output <- fmt.Sprintf("@%s %s", user.Username, modules.Word(modules.WordComfort))
}
