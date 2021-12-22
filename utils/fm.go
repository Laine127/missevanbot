package utils

import "github.com/google/uuid"

// MessageID return a special UUID, only for missevan.
func MessageID() string {
	u := uuid.NewString()
	return "3" + u[1:]
}
