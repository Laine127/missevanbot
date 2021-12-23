package handlers

import (
	"fmt"

	"missevanbot/models"
)

// TODO: the number of gifts may be incorrect.
func eventSend(output chan<- string, textMsg models.FmTextMessage) {
	if username := textMsg.User.Username; username != "" {
		gift := textMsg.Gift
		output <- fmt.Sprintf(models.TplThankGift, username, gift.Number, gift.Name)
	}
}
