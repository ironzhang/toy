package report

import (
	"math"
	"time"
)

type DataPoint struct {
	Timestamp  int64
	Throughput int64
	MaxLatency time.Duration
	MinLatency time.Duration
	AvgLatency time.Duration
}

type TimeSeries []DataPoint

func makeTimeSeries(records []Record) TimeSeries {
	type point struct {
		max   time.Duration
		min   time.Duration
		sum   time.Duration
		count int64
	}

	maxTS := int64(math.MinInt64)
	minTS := int64(math.MaxInt64)
	points := make(map[int64]point)
	for _, r := range records {
		ts := r.Start.Unix()
		if maxTS < ts {
			maxTS = ts
		}
		if minTS > ts {
			minTS = ts
		}
		if v, ok := points[ts]; !ok {
			points[ts] = point{max: r.Latency, min: r.Latency, sum: r.Latency, count: 1}
		} else {
			v.max = maxDuration(v.max, r.Latency)
			v.min = minDuration(v.min, r.Latency)
			v.sum += r.Latency
			v.count++
			points[ts] = v
		}
	}

	if maxTS >= minTS {
		series := make(TimeSeries, maxTS-minTS+1)
		for i := minTS; i <= maxTS; i++ {
			if v, ok := points[i]; ok {
				series[i-minTS] = DataPoint{
					Timestamp:  i,
					Throughput: v.count,
					MaxLatency: v.max,
					MinLatency: v.min,
					AvgLatency: v.sum / time.Duration(v.count),
				}
			} else {
				series[i-minTS] = DataPoint{Timestamp: i}
			}
		}
		return series
	}
	return nil
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func (s TimeSeries) Aggregation() DataPoint {
	n := len(s)
	if n <= 0 {
		return DataPoint{}
	} else if n == 1 {
		return s[0]
	}

	p := s[0]
	for i := 1; i < n; i++ {
		p.Throughput += s[i].Throughput
		p.AvgLatency += s[i].AvgLatency
		if s[i].MaxLatency > p.MaxLatency {
			p.MaxLatency = s[i].MaxLatency
		}
		if s[i].MinLatency < p.MinLatency {
			p.MinLatency = s[i].MinLatency
		}
	}
	p.Throughput /= int64(n)
	p.AvgLatency /= time.Duration(n)
	return p
}

func (s TimeSeries) Sampling(size int) TimeSeries {
	n := len(s)
	if n <= size || size <= 0 {
		return s
	}

	freq := n / size
	samples := make(TimeSeries, size)
	for i := 0; i < size; i++ {
		start := i * freq
		end := (i + 1) * freq
		samples[i] = s[start:end].Aggregation()
	}
	return samples
}
