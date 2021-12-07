package game

import (
	"context"
	"fmt"
	"strconv"

	"missevan-fm/config"
)

type RankMember struct {
	UID   int
	Score int
}

var ctx = context.Background()

// UpdateScore update the score of a player which specified by UID.
func UpdateScore(roomID, uid, score int) {
	rdb := config.RDB

	prefix := config.RedisPrefix + strconv.Itoa(roomID) // Redis namespace prefix, `missevan:[roomID]`

	key := fmt.Sprintf("%s:game", prefix)

	rdb.ZIncrBy(ctx, key, float64(score), strconv.Itoa(uid))
}

// ScoreRank return list of the members which in rank list.
func ScoreRank(roomID int) ([]RankMember, error) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(roomID)

	key := fmt.Sprintf("%s:game", prefix)

	slice := rdb.ZRangeWithScores(ctx, key, 0, -1)

	rankM := make([]RankMember, len(slice.Val()))
	for k, v := range slice.Val() {
		uid, err := strconv.Atoi(v.Member.(string))
		if err != nil {
			return nil, err
		}
		rankM[k] = RankMember{
			UID:   uid,
			Score: int(v.Score),
		}
	}
	return rankM, nil
}
