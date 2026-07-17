package bootstrap

import (
	"time"
)

// getCurrentTimestamp returns the current time in ISO 8601 format.
func getCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// getCurrentTimestampShort returns the current date in YYYY-MM-DD format.
func getCurrentTimestampShort() string {
	return time.Now().UTC().Format("2006-01-02")
}
