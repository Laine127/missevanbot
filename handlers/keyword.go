package handlers

import (
	"fmt"

	"missevanbot/models"
)

// keyEmotional handle the keyword `emo`.
func keyEmotional(output chan<- string, user models.FmUser) {
	output <- fmt.Sprintf("@%s %s", user.Username, models.ComfortString())
}
