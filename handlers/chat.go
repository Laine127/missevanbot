package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"missevan-fm/models"
)

// Chat 简单回复
func Chat(username string) string {
	pool := models.ChatPool
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("@%s %s", username, pool[r.Intn(len(pool))])
}
