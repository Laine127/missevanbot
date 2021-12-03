package thirdparty

import (
	"log"
	"testing"
)

func TestFetch(t *testing.T) {
	text, err := Fetch(88921)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range text {
		log.Println(v)
	}
}
