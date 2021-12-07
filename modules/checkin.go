package modules

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules/thirdparty"
	"missevan-fm/utils"
)

var ctx = context.Background()

// Checkin do the action of user checkin.
func Checkin(roomID, uid int, uname string) (string, error) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(roomID) // Redis namespace prefix, `missevan:[roomID]`

	key := fmt.Sprintf("%s:checkin:%d", prefix, uid) // `missevan:[roomID]:checkin:[UID]`

	// get the checkin data from the Redis cache.
	cmd := rdb.HMGet(ctx, key, "count", "day", "luck")
	var count, day, luck string
	if v, ok := cmd.Val()[0].(string); ok {
		count = v
	}
	if v, ok := cmd.Val()[1].(string); ok {
		day = v
	}
	if v, ok := cmd.Val()[2].(string); ok {
		luck = v
	} else {
		// there is no `luck` record.
		luck = models.LuckString()
		rdb.HSet(ctx, key, "luck", luck)
	}
	// check if already checkin that day.
	if day == utils.Today() {
		return fmt.Sprintf(models.TplSignDuplicate, count, luck), nil
	}
	// check if checkin for consecutive days,
	// if not, set the `count` to 0.
	if day != time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		rdb.HSet(ctx, key, "count", 0)
	}
	// generate the string of `luck` that day.
	luck = models.LuckString()
	rdb.HMSet(ctx, key, "day", utils.Today(), "luck", luck)
	// increase the count of consecutive checkin days.
	countCMD := rdb.HIncrBy(ctx, key, "count", 1)
	// push the UID and Username into the rank list.
	rdb.RPush(ctx, fmt.Sprintf("%s:rank:%s:id", prefix, utils.Today()), uid)
	rdb.RPush(ctx, fmt.Sprintf("%s:rank:%s:name", prefix, utils.Today()), uname)

	poem, err := thirdparty.PoemText()
	if err != nil {
		poem = models.TplDefaultPoem
	}

	return fmt.Sprintf(models.TplSignSuccess, countCMD.Val(), luck, poem), nil
}

// CheckinRank return the rank of checkin task today.
func CheckinRank(roomID int) string {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(roomID)

	key := fmt.Sprintf("%s:rank:%s:id", prefix, utils.Today())
	nKey := fmt.Sprintf("%s:rank:%s:name", prefix, utils.Today())
	listID := rdb.LRange(ctx, key, 0, -1)
	listName := rdb.LRange(ctx, nKey, 0, -1)
	result := strings.Builder{}
	result.WriteString("\n")
	for k, v := range listID.Val() {
		count := rdb.HGet(ctx, fmt.Sprintf("%s:checkin:%s", prefix, v), "count").Val()
		result.WriteString(fmt.Sprintf("Rank %d. [%s] 连续签到%s天", k+1, listName.Val()[k], count))
		if k < len(listID.Val())-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
