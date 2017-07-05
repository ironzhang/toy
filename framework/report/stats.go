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

func (s *Stats) Histogram() (buckets []time.Duration, counts []int, max int) {
	bc := 10
	buckets = make([]time.Duration, bc+1)
	counts = make([]int, bc+1)
	bs := (s.Slowest - s.Fastest) / time.Duration(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = s.Fastest + time.Duration(i)*bs
	}
	buckets[bc] = s.Slowest

	bi := 0
	for i := 0; i < len(s.Lats); {
		if s.Lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	return buckets, counts, max
}

var pctls = []float64{10, 25, 50, 75, 90, 95, 99, 99.9}

func (s *Stats) Latencies() (pcs []float64, data []time.Duration) {
	if len(s.Lats) <= 0 {
		return
	}

	i := -1
	n := float64(len(s.Lats))
	pcs = make([]float64, 0, len(pctls))
	data = make([]time.Duration, 0, len(pctls))
	for _, p := range pctls {
		l := s.Lats[int(p/100*n)]
		if i >= 0 && data[i] == l {
			pcs[i] = p
		} else {
			pcs = append(pcs, p)
			data = append(data, l)
			i++
		}
	}
	return
}
