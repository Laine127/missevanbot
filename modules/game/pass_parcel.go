package game

import (
	"strconv"

	"missevan-fm/utils"
)

const (
	keyParcel = "parcel"
)

// ParcelNumber 生成一个数字，五位数
func ParcelNumber(m map[string]int) string {
	number := utils.RandomNumber(10000, 99999)
	m[keyParcel] = number
	str := strconv.FormatInt(int64(number), 10)
	return str
}

// IsParcelCorrect 判断输入是否正确
func IsParcelCorrect(m map[string]int, s string) bool {
	number := m[keyParcel]
	str := strconv.FormatInt(int64(number), 10)
	return s == str
}
