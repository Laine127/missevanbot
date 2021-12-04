package utils

import "time"

// Today 获取格式化的当天日期
func Today() string {
	return time.Now().Format("2006-01-02")
}
