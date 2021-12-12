package modules

import (
	"log"
	"testing"

	"missevan-fm/config"
)

func init() {
	config.LoadConfig()
}

func TestCookie(t *testing.T) {
	cookie, err := ConnCookie()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cookie)
}

func TestChangeAttention(t *testing.T) {
	t.Run("follow", func(t *testing.T) {
		ret, err := ChangeAttention(11111, Follow)
		if err != nil {
			log.Println(err)
			return
		}
		t.Log(string(ret))
	})
	t.Run("unfollow", func(t *testing.T) {
		ret, err := ChangeAttention(11111, Unfollow)
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
