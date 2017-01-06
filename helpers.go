package timsel

import "time"

func Ago(ago time.Duration) time.Time {
	return time.Now().Add(ago * -1)
}
