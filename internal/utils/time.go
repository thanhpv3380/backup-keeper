package utils

import "time"

func GenerateTimestamp() string {
	return time.Now().Format("20060102_150405")
}
