package thirdparty

import (
	"testing"
)

func TestFetch(t *testing.T) {
	roles, paras, err := DramaScript(95915)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range roles {
		t.Log(v)
	}

	for _, v := range paras {
		t.Log(v)
	}
}
