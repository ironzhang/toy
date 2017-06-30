package tbr

import (
	"sort"
	"time"
)

type Stats struct {
	Name        string
	QPS         int
	RealQPS     int
	Request     int
	RealRequest int
	Concurrent  int
	Total       time.Duration
	Fastest     time.Duration
	Slowest     time.Duration
	Average     time.Duration
	Series      TimeSeries
	Lats        []time.Duration
	Errs        map[string]int
}

type Record struct {
	Err     string
	Start   time.Time
	Latency time.Duration
}

type Result struct {
	Name       string
	QPS        int
	Request    int
	Concurrent int
	Total      time.Duration
	Records    []Record
}

func (r *Result) Stats() Stats {
	fastest, slowest, average, lats, errs := doStats(r.Records)
	return Stats{
		Name:        r.Name,
		QPS:         r.QPS,
		RealQPS:     calcQPS(len(lats), r.Total),
		Request:     r.Request,
		RealRequest: len(lats),
		Concurrent:  r.Concurrent,
		Total:       r.Total,
		Fastest:     fastest,
		Slowest:     slowest,
		Average:     average,
		Series:      makeTimeSeries(r.Records),
		Lats:        lats,
		Errs:        errs,
	}
}

func doStats(records []Record) (fastest, slowest, average time.Duration, lats []time.Duration, errs map[string]int) {
	var sum time.Duration
	errs = make(map[string]int)
	lats = make([]time.Duration, 0, len(records))
	for _, r := range records {
		if r.Err != "" {
			errs[r.Err]++
		} else {
			sum += r.Latency
			lats = append(lats, r.Latency)
		}
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
	if len(lats) > 0 {
		n := len(lats)
		fastest = lats[0]
		slowest = lats[n-1]
		average = sum / time.Duration(n)
	}
	return
}

func calcQPS(n int, d time.Duration) int {
	return int(float64(n) / d.Seconds())
}

type Results []Result

func (p Results) AddResult(r Result) {
}
