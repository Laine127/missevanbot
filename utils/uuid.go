package utils

import "github.com/google/uuid"

// MessageID 用于生成消息唯一ID
func MessageID() string {
	u := uuid.NewString()
	return "3" + u[1:]
}
