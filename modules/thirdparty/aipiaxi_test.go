package thirdparty

import (
	"fmt"
	"log"
	"testing"
)

func TestFetch(t *testing.T) {
	roles, text, err := Fetch(95915)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range roles {
		log.Println(v)
	}
	fmt.Print("\n")
	for _, v := range text {
		log.Println(v)
	}
}
