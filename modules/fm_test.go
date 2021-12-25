package modules

import (
	"log"
	"testing"

	"missevanbot/config"
)

func init() {
	config.LoadConfig()
}

func TestBaseCookie(t *testing.T) {
	cookie, err := BaseCookie()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cookie)
}

func TestChangeAttention(t *testing.T) {
	t.Run("follow", func(t *testing.T) {
		ret, err := ChangeAttention(cookie(0), 11111, Follow)
		if err != nil {
			log.Println(err)
			return
		}
		t.Log(string(ret))
	})
	t.Run("unfollow", func(t *testing.T) {
		ret, err := ChangeAttention(cookie(0), 11111, Unfollow)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(ret))
	})
}

func TestQueryUsername(t *testing.T) {
	tests := []struct {
		name    string
		uid     int
		want    string
		wantErr bool
	}{
		{"小熊家的软糖i", 18395018, "小熊家的软糖i", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryUsername(tt.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryUsername() got = %v, want %v", got, tt.want)
			}
		})
	}
}
