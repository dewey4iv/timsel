package timsel

import (
	"sort"
	"time"

	"github.com/dewey4iv/number"
)

func New(input map[time.Time]*number.N) TimeSelector {
	return &TimSel{
		data: input,
	}
}

type TimSel struct {
	from time.Time
	to   time.Time
	data map[time.Time]*number.N
}

func (ts *TimSel) From(from time.Time) TimeSelector {
	data := make(map[time.Time]*number.N)

	for t, n := range ts.data {
		if t.Unix() >= from.Unix() {
			data[t] = n
		}
	}

	result := TimSel{
		from: from,
		to:   ts.to,
		data: data,
	}

	return &result
}

func (ts *TimSel) To(to time.Time) TimeSelector {

	data := make(map[time.Time]*number.N)

	for t, n := range ts.data {
		if t.Unix() <= to.Unix() {
			data[t] = n
		}
	}

	result := TimSel{
		from: ts.from,
		to:   to,
		data: data,
	}

	return &result
}

func (ts *TimSel) Compress(groupBy time.Duration) TimeSelector {
	if ts.from.IsZero() {
		first := true
		for t, _ := range ts.data {
			if first {
				ts.from = t
			}

			if ts.from.Unix() > t.Unix() {
				ts.from = t
			}
		}
	}

	data := make(map[time.Time]*number.N)
	current := ts.from
	for t, n := range ts.data {
		if t.Unix() >= current.Unix() && t.Unix() < current.Add(groupBy).Unix() {
			data[current] = number.New(data[current].Float() + n.Float())
		}

		current = current.Add(groupBy)
	}

	result := TimSel{from: ts.from, to: ts.to, data: data}

	return &result
}

func (ts *TimSel) GroupBy(groupBy time.Duration) []TimeSelector {
	if ts.from.IsZero() {
		first := true
		for t, _ := range ts.data {
			if first {
				ts.from = t
			}

			if ts.from.Unix() > t.Unix() {
				ts.from = t
			}
		}
	}

	times := make([]int, len(ts.data))
	for t, _ := range ts.data {
		times = append(times, int(t.Unix()))
	}

	sort.Ints(times)

	var timeSelectors []TimeSelector

	data := make(map[time.Time]*number.N)
	current := ts.from
	for j := range times {
		i := int64(times[j])
		t := time.Unix(i, 0)

		if i >= current.Unix() && i < current.Add(groupBy).Unix() {
			data[t] = ts.data[t]
		} else {
			timeSelectors = append(timeSelectors, &TimSel{data: data})
			current = current.Add(groupBy)
			data = make(map[time.Time]*number.N)
		}
	}

	return timeSelectors
}

func (ts *TimSel) Average() *number.N {
	var total float64
	var len float64

	for _, n := range ts.MapTimeNumber() {
		total += n.Float()
		len++
	}

	return number.New(total / len)
}

func (ts *TimSel) Total() *number.N {
	var total float64

	for _, n := range ts.MapTimeNumber() {
		total += n.Float()
	}

	return number.New(total)
}

func (ts *TimSel) MapTimeNumber() map[time.Time]*number.N {
	return ts.data
}

func (ts *TimSel) MapTimeFloat() map[time.Time]float64 {
	data := make(map[time.Time]float64)

	for t, n := range ts.data {
		data[t] = n.Float()
	}

	return data
}

func (ts *TimSel) Slice() []TimeNumber {
	var results []TimeNumber

	for t, n := range ts.data {
		results = append(results, &timeNumber{t, n})
	}

	return results
}

type timeNumber struct {
	t time.Time
	n *number.N
}

func (tn *timeNumber) Time() time.Time {
	return tn.t
}

func (tn *timeNumber) Number() *number.N {
	return tn.n
}
