package thirdparty

import (
	"log"
	"testing"
)

func TestZodiac(t *testing.T) {
	ret, err := Zodiac("摩羯座", Level1)
	if err != nil {
		log.Println()
		return
	}
	log.Println(ret.Content)
}
