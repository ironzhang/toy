package tbr

import (
	"reflect"
	"testing"
	"time"
)

func TestMakeTimeSeries(t *testing.T) {
	ts := time.Time{}
	tests := []struct {
		records []Record
		series  TimeSeries
	}{
		{
			records: nil,
			series:  nil,
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
			},
			series: TimeSeries{
				{Timestamp: ts.Unix(), ThroughPut: 1, MaxLatency: time.Second, MinLatency: time.Second, AvgLatency: time.Second},
			},
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
				{Start: ts, Latency: 2 * time.Second},
				{Start: ts, Latency: 3 * time.Second},
			},
			series: TimeSeries{
				{Timestamp: ts.Unix(), ThroughPut: 3, MaxLatency: 3 * time.Second, MinLatency: time.Second, AvgLatency: 2 * time.Second},
			},
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
				{Start: ts.Add(100 * time.Millisecond), Latency: 2 * time.Second},
				{Start: ts.Add(200 * time.Millisecond), Latency: 3 * time.Second},
			},
			series: TimeSeries{
				{Timestamp: ts.Unix(), ThroughPut: 3, MaxLatency: 3 * time.Second, MinLatency: time.Second, AvgLatency: 2 * time.Second},
			},
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
				{Start: ts.Add(100 * time.Millisecond), Latency: 2 * time.Second},
				{Start: ts.Add(200 * time.Millisecond), Latency: 3 * time.Second},
				{Start: ts.Add(1*time.Second + 100*time.Millisecond), Latency: 4 * time.Second},
				{Start: ts.Add(1*time.Second + 200*time.Millisecond), Latency: 5 * time.Second},
				{Start: ts.Add(1*time.Second + 300*time.Millisecond), Latency: 6 * time.Second},
				{Start: ts.Add(1*time.Second + 400*time.Millisecond), Latency: 7 * time.Second},
				{Start: ts.Add(1*time.Second + 500*time.Millisecond), Latency: 8 * time.Second},
			},
			series: TimeSeries{
				{Timestamp: ts.Unix(), ThroughPut: 3, MaxLatency: 3 * time.Second, MinLatency: time.Second, AvgLatency: 2 * time.Second},
				{Timestamp: ts.Add(1 * time.Second).Unix(), ThroughPut: 5, MaxLatency: 8 * time.Second, MinLatency: 4 * time.Second, AvgLatency: 6 * time.Second},
			},
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
				{Start: ts.Add(100 * time.Millisecond), Latency: 2 * time.Second},
				{Start: ts.Add(200 * time.Millisecond), Latency: 3 * time.Second},
				{Start: ts.Add(2*time.Second + 100*time.Millisecond), Latency: 4 * time.Second},
				{Start: ts.Add(2*time.Second + 200*time.Millisecond), Latency: 5 * time.Second},
				{Start: ts.Add(2*time.Second + 300*time.Millisecond), Latency: 6 * time.Second},
				{Start: ts.Add(2*time.Second + 400*time.Millisecond), Latency: 7 * time.Second},
				{Start: ts.Add(2*time.Second + 500*time.Millisecond), Latency: 8 * time.Second},
			},
			series: TimeSeries{
				{Timestamp: ts.Unix(), ThroughPut: 3, MaxLatency: 3 * time.Second, MinLatency: time.Second, AvgLatency: 2 * time.Second},
				{Timestamp: ts.Add(1 * time.Second).Unix(), ThroughPut: 0, MaxLatency: 0, MinLatency: 0, AvgLatency: 0},
				{Timestamp: ts.Add(2 * time.Second).Unix(), ThroughPut: 5, MaxLatency: 8 * time.Second, MinLatency: 4 * time.Second, AvgLatency: 6 * time.Second},
			},
		},
	}
	for i, tt := range tests {
		if got, want := makeTimeSeries(tt.records), tt.series; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: %v != %v", i, got, want)
		}
	}
}
