package handlers

import (
	"fmt"
	"time"

	"missevanbot/models"
	"missevanbot/modules"
)

type spellsHandler func(room *models.Room, user models.FmUser, output chan<- string)

var _spellsMap = map[int]spellsHandler{
	models.SpellsInvisible: spellsInvisible,
	models.SpellsVisible:   spellsVisible,
}

func spellsInvisible(room *models.Room, user models.FmUser, output chan<- string) {
	modules.InvisibleSet(room.ID, user.UserID, time.Hour)
	output <- fmt.Sprintf("@%s 恭喜您解锁了掩耳盗铃成就，隐身时长一小时～", user.Username)
}

func spellsVisible(room *models.Room, user models.FmUser, output chan<- string) {
	modules.InvisibleCancel(room.ID, user.UserID)
	output <- fmt.Sprintf("@%s 我现在又能看到你了", user.Username)
}
