package modules

import (
	"testing"

	"missevanbot/config"
)

const TestRID = 1000
const TestUID = 5163862

func init() {
	config.LoadConfig()
	conf := config.Config()
	config.InitRDBClient(conf.Redis)
}

func TestUpdateScore(t *testing.T) {
	UpdateScore(TestRID, TestUID, 1)
}

func TestScoreRank(t *testing.T) {
	res, err := ScoreRank(TestRID)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range res {
		t.Log(v)
	}
}
