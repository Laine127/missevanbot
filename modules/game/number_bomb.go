package game

import (
	"math/rand"
	"time"
)

// BombGenerate 返回一个值为 1 - max 的随机数
func BombGenerate(m map[string]int, players int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	max := players * 30
	bomb := rand.Intn(max) + 1

	m["bomb_min"] = 1
	m["bomb_num"] = bomb
	m["bomb_max"] = max

	return 1, max
}

// BombGuess 猜数字
func BombGuess(m map[string]int, n int) (bool, int, int) {
	min := m["bomb_min"]
	bomb := m["bomb_num"]
	max := m["bomb_max"]

	if n < min || n > max {
		// 越界
		return false, -1, -1
	}

	if n == bomb {
		return true, 0, 0
	} else if n < bomb {
		m["bomb_min"] = n
		return false, n + 1, max
	} else {
		m["bomb_max"] = n
		return false, min, n - 1
	}
}
