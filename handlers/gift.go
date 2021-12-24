package handlers

import (
	"fmt"

	"missevanbot/models"
)

// TODO: the number of gifts may be incorrect.
func eventSend(output chan<- string, textMsg models.FmTextMessage) {
	if name := textMsg.User.Username; name != "" {
		gift := textMsg.Gift
		output <- fmt.Sprintf(models.TplThankGift, name, gift.Number, gift.Name)
	}
}
