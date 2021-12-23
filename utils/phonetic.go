package utils

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// Pinyin returns the pinyin of s.
func Pinyin(s string) string {
	text := strings.Builder{}
	args := pinyin.NewArgs()
	args.Style = pinyin.Tone // enable tone
	args.Heteronym = true    // enable heteronym
	arr := pinyin.Pinyin(s, args)
	for i, multi := range arr {
		text.WriteString("[")
		if len(multi) == 1 {
			text.WriteString(multi[0])
		}
		if len(multi) > 1 {
			text.WriteString(multi[0] + " " + multi[1])
		}
		text.WriteString("]")
		// Insert a whitespace between each element.
		if i < len(arr)-1 {
			text.WriteString(" ")
		}
	}
	return text.String()
}
