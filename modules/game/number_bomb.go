package game

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	KeyBombMin = "bomb_min"
	KeyBombMax = "bomb_max"
	KeyBombNum = "bomb_num"
)

// BombGenerate set the range by the number of players,
// generate a random number in the range,
// store values into m, return the range boundary.
func BombGenerate(m map[string]int, players int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	max := players * 30
	bomb := rand.Intn(max) + 1

	m[KeyBombMin] = 1
	m[KeyBombNum] = bomb
	m[KeyBombMax] = max

	return 1, max
}

// BombGuess check if bomb corrected,
// if not, return the new range boundary.
func BombGuess(m map[string]int, s string) (bool, int, int) {
	min := m[KeyBombMin]
	bomb := m[KeyBombNum]
	max := m[KeyBombMax]

	n, err := strconv.Atoi(s)
	if err != nil || n < min || n > max {
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
