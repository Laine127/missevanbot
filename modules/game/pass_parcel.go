package game

import (
	"strconv"

	"missevan-fm/utils"
)

const (
	keyParcel = "parcel"
)

// ParcelNumber generate a random five-digit number,
// and store it into m.
func ParcelNumber(m map[string]int) string {
	number := utils.RandomNumber(10000, 99999)
	m[keyParcel] = number
	str := strconv.FormatInt(int64(number), 10)
	return str
}

// IsParcelCorrect check if int(s) equals to the value stored in m.
func IsParcelCorrect(m map[string]int, s string) bool {
	number := m[keyParcel]
	str := strconv.FormatInt(int64(number), 10)
	return s == str
}
