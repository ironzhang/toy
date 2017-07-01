package report

import (
	"sort"
	"time"
)

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

func (r *Result) merge(a Result) {
	r.QPS += a.QPS
	r.Request += a.Request
	r.Concurrent += a.Concurrent
	if r.Total < a.Total {
		r.Total = a.Total
	}
	r.Records = append(r.Records, a.Records...)
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
	x := d.Seconds()
	if x == 0 {
		panic("calcQPS: duration is zoro")
	}
	return int(float64(n) / x)
}

type Results []Result

func (p *Results) AddResult(a Result) {
	for i, r := range *p {
		if r.Name == a.Name {
			(*p)[i].merge(a)
			return
		}
	}
	*p = append(*p, a)
}
