package timsel

import (
	"time"

	"github.com/dewey4iv/number"
)

type TimeSelector interface {
	From(time.Time) TimeSelector
	To(time.Time) TimeSelector
	Compress(time.Duration) TimeSelector
	GroupBy(time.Duration) []TimeSelector
	MapTimeFloat() map[time.Time]float64
	MapTimeNumber() map[time.Time]*number.N
	Slice() []TimeNumber
}

type TimeNumber interface {
	Time() time.Time
	Number() *number.N
}
