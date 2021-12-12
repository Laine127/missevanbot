package handlers

import (
	"fmt"

	"missevan-fm/models"
)

// keyEmotional handle the keyword `emo`.
func keyEmotional(outputMsg chan<- string, user models.FmUser) {
	outputMsg <- fmt.Sprintf("@%s %s", user.Username, models.ComfortString())
}
