package game

import (
	"fmt"
	"log"
	"testing"

	"missevan-fm/config"
)

func init() {
	config.LoadConfig()
	conf := config.Config()
	config.InitRDBClient(conf.Redis)
}

func TestUpdateScore(t *testing.T) {
	UpdateScore(461808808, 5163862, 1)
}

func TestScoreRank(t *testing.T) {
	res, err := ScoreRank(461808808)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
