package report

import "time"

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
