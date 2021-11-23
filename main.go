package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	rid := os.Args[1]
	// 455984011
	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println("Room ID错误。。。")
		return
	}
	connect(id)
}
