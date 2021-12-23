package utils

import "time"

// Today returns the format string of today's date.
func Today() string {
	return time.Now().Format("2006-01-02")
}
