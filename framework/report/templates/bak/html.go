package report

import "time"

type latency struct {
	Percent  int
	Duration time.Duration
}

type html struct {
	Name        string
	Total       time.Duration
	Slowest     time.Duration
	Fastest     time.Duration
	Average     time.Duration
	Concurrent  int
	Request     int
	RealRequest int
	QPS         int
	RealQPS     float64
	LatencyImg  string
	Latencys    []latency
	Errs        map[string]int
}

func HTML(outdir string, reports []Report) error {
}
