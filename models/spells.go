package models

const (
	SpellsReparo = iota
	SpellsIncendio
	SpellsAvis
	SpellsInvisible
	SpellsVisible
)

var _spellsMap = map[string]int{
	// "REPARO":    SpellsReparo,
	// "INCENDIO":  SpellsIncendio,  // 火焰熊熊
	// "AVIS":      SpellsAvis,      // 飞鸟群群
	"INVISIBLE": SpellsInvisible, // 隐身
	"VISIBLE":   SpellsInvisible, // 取消隐身
}

func Spells(key string) int {
	if v, ok := _spellsMap[key]; ok {
		return v
	}
	return -1
}
