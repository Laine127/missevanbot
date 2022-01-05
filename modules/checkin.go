package modules

import (
	"fmt"
	"strings"
	"time"

	"missevanbot/models"
	"missevanbot/modules/thirdparty"
	"missevanbot/utils"
)

// Checkin do the action of user checkin.
// TODO: enhance this function.
func Checkin(rid, uid int, name string) (string, error) {
	prefix := prefixRoom(rid)
	key := prefix + fmt.Sprintf("checkin:%d", uid) // missevan:[RoomID]:checkin:[UID]

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
		luck = "今日运势：" + Word(WordLuck)
		rdb.HSet(ctx, key, "luck", luck)
	}
	// check if already checkin that day.
	if day == utils.Today() {
		return fmt.Sprintf(models.TplSignDuplicate, name, count, luck), nil
	}
	// check if checkin for consecutive days,
	// if not, set the `count` to 0.
	if day != time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		rdb.HSet(ctx, key, "count", 0)
	}
	// generate the string of `luck` that day.
	luck = "今日运势：" + Word(WordLuck)
	rdb.HMSet(ctx, key, "day", utils.Today(), "luck", luck)
	// increase the count of consecutive checkin days.
	countCMD := rdb.HIncrBy(ctx, key, "count", 1)
	// push the UID and Username into the rank list.
	rdb.RPush(ctx, prefix+fmt.Sprintf("rank:%s:id", utils.Today()), uid)
	rdb.RPush(ctx, prefix+fmt.Sprintf("rank:%s:name", utils.Today()), name)

	poem, err := thirdparty.PoemText()
	if err != nil {
		poem = models.TplDefaultPoem
	}

	return fmt.Sprintf(models.TplSignSuccess, name, countCMD.Val(), luck, poem), nil
}

// CheckinRank return the rank of checkin task today.
func CheckinRank(rid int) string {
	prefix := prefixRoom(rid)
	key := prefix + fmt.Sprintf("rank:%s:id", utils.Today())
	nKey := prefix + fmt.Sprintf("rank:%s:name", utils.Today())

	listID := rdb.LRange(ctx, key, 0, -1)
	listName := rdb.LRange(ctx, nKey, 0, -1)
	result := strings.Builder{}
	result.WriteString("\n")
	for k, v := range listID.Val() {
		count := rdb.HGet(ctx, prefix+fmt.Sprintf("checkin:%s", v), "count").Val()
		result.WriteString(fmt.Sprintf("Rank %d. [%s] 连续签到%s天", k+1, listName.Val()[k], count))
		if k < len(listID.Val())-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
