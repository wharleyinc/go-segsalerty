package common

import "time"

// TimeNow wrapper returns the current time in UTC
func TimeNow() time.Time {
	return time.Now().UTC()
}
