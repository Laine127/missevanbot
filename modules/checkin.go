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

// Checkin 用户签到
func Checkin(roomID, uid int, uname string) (string, error) {
	rdb := config.RDB
	prefix := models.RedisPrefix + strconv.Itoa(roomID)

	key := fmt.Sprintf("%s:checkin:%d", prefix, uid)
	// 获取当前缓存中的签到信息
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
		// 还没有运势记录
		luck = models.LuckString()
		rdb.HSet(ctx, key, "luck", luck)
	}
	// 判断是否重复签到
	if day == utils.Today() {
		return fmt.Sprintf(models.TplSignDuplicate, count, luck), nil
	}
	// 判断是否连续签到
	if day != time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		rdb.HSet(ctx, key, "count", 0)
	}
	// 生成今天的运势
	luck = models.LuckString()
	rdb.HMSet(ctx, key, "day", utils.Today(), "luck", luck)
	// 增加签到天数
	countCMD := rdb.HIncrBy(ctx, key, "count", 1)
	// 放入排行榜
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
	prefix := models.RedisPrefix + strconv.Itoa(roomID)

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
