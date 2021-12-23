package modules

import "missevanbot/config"

const (
	WordChat    = "chat"
	WordComfort = "comfort"
	WordGuess   = "guess"
	WordLuck    = "luck"
	WordReply   = "reply"
	WordWelcome = "welcome"
)

func Word(key string) string {
	rdb := config.RDB
	key = config.RedisPrefixWords + key
	c := rdb.SRandMember(ctx, key)
	return c.Val()
}
