package modules

import (
	"log"
	"testing"

	"missevan-fm/config"
)

func TestCookie(t *testing.T) {
	cookie, err := ConnCookie()
	log.Println(err)
	log.Println(cookie)
}

func TestChangeAttention(t *testing.T) {
	config.LoadConfig()
	t.Run("follow", func(t *testing.T) {
		ret, err := ChangeAttention(11111, 1)
		log.Println(string(ret))
		log.Println(err)
	})
	t.Run("unfollow", func(t *testing.T) {
		ret, err := ChangeAttention(11111, 1)
		log.Println(string(ret))
		log.Println(err)
	})
}

func TestUnfollowAll(t *testing.T) {
	config.LoadConfig()
	InitBot()
	t.Run("unfollow_all", func(t *testing.T) {
		log.Println(UnfollowAll())
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
