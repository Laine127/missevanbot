package modules

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules/thirdparty"
)

const RedisPrefix = "missevan:"

var ctx = context.Background()

// Sign 用户签到
func Sign(roomID, uid int, uname string) (string, error) {
	rdb := config.RDB
	prefix := RedisPrefix + strconv.Itoa(roomID)

	key := fmt.Sprintf("%s:sign:%d", prefix, uid)
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
		luck = luckString()
		rdb.HSet(ctx, key, "luck", luck)
	}
	// 判断是否重复签到
	if day == time.Now().Format("2006-01-02") {
		return fmt.Sprintf(models.TplSignDuplicate, count, luck), nil
	}
	// 判断是否连续签到
	if day != time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		rdb.HSet(ctx, key, "count", 0)
	}
	// 生成今天的运势
	luck = luckString()
	rdb.HMSet(ctx, key, "day", time.Now().Format("2006-01-02"), "luck", luck)
	// 增加签到天数
	countCMD := rdb.HIncrBy(ctx, key, "count", 1)
	// 放入排行榜
	rdb.RPush(ctx, fmt.Sprintf("%s:rank:%s:id", prefix, time.Now().Format("2006-01-02")), uid)
	rdb.RPush(ctx, fmt.Sprintf("%s:rank:%s:name", prefix, time.Now().Format("2006-01-02")), uname)

	poem, err := thirdparty.PoemText()
	if err != nil {
		poem = models.TplDefaultPoem
	}

	return fmt.Sprintf(models.TplSignSuccess, countCMD.Val(), luck, poem), nil
}

// Rank return the rank of sign task today.
func Rank(roomID int) string {
	rdb := config.RDB
	prefix := RedisPrefix + strconv.Itoa(roomID)

	key := fmt.Sprintf("%s:rank:%s:id", prefix, time.Now().Format("2006-01-02"))
	nKey := fmt.Sprintf("%s:rank:%s:name", prefix, time.Now().Format("2006-01-02"))
	listID := rdb.LRange(ctx, key, 0, -1)
	listName := rdb.LRange(ctx, nKey, 0, -1)
	result := strings.Builder{}
	result.WriteString("\n")
	for k, v := range listID.Val() {
		count := rdb.HGet(ctx, fmt.Sprintf("%s:sign:%s", prefix, v), "count").Val()
		result.WriteString(fmt.Sprintf("Rank %d. [%s] 连续签到%s天", k+1, listName.Val()[k], count))
		if k < len(listID.Val())-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

// luckString 返回运势字符串
func luckString() string {
	result := "今日运势："
	pool := models.LuckPool
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return result + pool[r.Intn(len(pool))]
}
