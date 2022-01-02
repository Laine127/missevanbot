package utils

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// Pinyin returns the pinyin of s.
func Pinyin(s string) (string, string) {
	args := pinyin.NewArgs()
	args.Style = pinyin.Tone // enable tone
	args.Heteronym = true    // enable heteronym
	arr := pinyin.Pinyin(s, args)
	textS := strings.Builder{}
	for i, multi := range arr {
		if len(multi) > 0 {
			textS.WriteString(multi[0])
		}
		if i < len(arr)-1 {
			textS.WriteString(" ")
		}
	}

	textM := strings.Builder{}
	for i, multi := range arr {
		textM.WriteString("[")
		if len(multi) == 1 {
			textM.WriteString(multi[0])
		}
		if len(multi) > 1 {
			textM.WriteString(multi[0] + " " + multi[1])
		}
		textM.WriteString("]")
		// Insert a whitespace between each element.
		if i < len(arr)-1 {
			textM.WriteString(" ")
		}
	}
	return textS.String(), textM.String()
}
