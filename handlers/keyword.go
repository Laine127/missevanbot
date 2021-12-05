package handlers

import (
	"fmt"

	"missevan-fm/models"
)

// keyEmotional 处理emo关键词
func keyEmotional(outputMsg chan<- string, user models.FmUser) {
	outputMsg <- fmt.Sprintf("@%s %s", user.Username, models.ComfortString())
}
