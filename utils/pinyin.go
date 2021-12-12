package utils

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// Pinyin return the pinyin of s.
func Pinyin(s string) string {
	str := strings.Builder{}
	py := pinyin.NewArgs()
	py.Style = pinyin.Tone // enable tone.
	// py.Heteronym = true // enable heteronym.
	arr := pinyin.Pinyin(s, py)
	for k, v := range arr {
		if len(v) == 1 {
			str.WriteString(v[0])
		}
		// insert whitespace between each pinyin.
		if k < len(arr)-1 {
			str.WriteString(" ")
		}
	}
	return str.String()
}
