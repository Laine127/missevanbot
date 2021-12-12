package thirdparty

import (
	"testing"
)

func TestZodiac(t *testing.T) {
	ret, err := Zodiac("摩羯座", Level1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret.Content)
}
