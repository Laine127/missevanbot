package modules

import (
	"testing"

	"missevanbot/config"
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
		t.Error(err)
		return
	}
	for _, v := range res {
		t.Log(v)
	}
}
