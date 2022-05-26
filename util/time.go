package util

import "time"

// Time2String Time to mm-dd time
func Time2String(t time.Time) string {
	return t.Format("01-02")
}
