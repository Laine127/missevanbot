package modules

import (
	"strconv"
)

type RankMember struct {
	UID   int
	Score int
}

// UpdateScore update the score of a player which specified by UID.
func UpdateScore(rid, uid, score int) {
	key := prefixRoom(rid) + "game" // "missevan:[RoomID]:game"

	rdb.ZIncrBy(ctx, key, float64(score), strconv.Itoa(uid))
}

// ScoreRank return list of the members which in rank list.
func ScoreRank(rid int) ([]RankMember, error) {
	key := prefixRoom(rid) + "game" // "missevan:[RoomID]:game"

	slice := rdb.ZRevRangeWithScores(ctx, key, 0, -1)

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
