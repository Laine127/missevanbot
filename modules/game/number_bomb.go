package game

import (
	"math/rand"
	"time"
)

const (
	KeyBombMin = "bomb_min"
	KeyBombMax = "bomb_max"
	KeyBombNum = "bomb_num"
)

// BombGenerate 返回一个值为 1 - max 的随机数
func BombGenerate(m map[string]int, players int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	max := players * 30
	bomb := rand.Intn(max) + 1

	m[KeyBombMin] = 1
	m[KeyBombNum] = bomb
	m[KeyBombMax] = max

	return 1, max
}

// BombGuess 猜数字
func BombGuess(m map[string]int, n int) (bool, int, int) {
	min := m[KeyBombMin]
	bomb := m[KeyBombNum]
	max := m[KeyBombMax]

	if n < min || n > max {
		// 越界
		return false, -1, -1
	}

	if n == bomb {
		return true, 0, 0
	} else if n < bomb {
		m[KeyBombMin] = n + 1
		return false, n + 1, max
	} else {
		m[KeyBombMax] = n - 1
		return false, min, n - 1
	}
}
