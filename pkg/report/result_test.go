package report

import (
	"reflect"
	"testing"
	"time"
)

func TestDoStats(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		records []Record
		fastest time.Duration
		slowest time.Duration
		average time.Duration
		lats    []time.Duration
		errs    map[string]int
	}{
		{
			records: nil,
			fastest: 0,
			slowest: 0,
			average: 0,
			lats:    []time.Duration{},
			errs:    map[string]int{},
		},

		{
			records: []Record{
				{Start: ts, Latency: time.Second},
			},
			fastest: time.Second,
			slowest: time.Second,
			average: time.Second,
			lats: []time.Duration{
				time.Second,
			},
			errs: map[string]int{},
		},

		{
			records: []Record{
				{Start: ts, Latency: 1 * time.Second},
				{Start: ts, Latency: 2 * time.Second},
			},
			fastest: 1 * time.Second,
			slowest: 2 * time.Second,
			average: (1 + 2) * time.Second / 2,
			lats: []time.Duration{
				1 * time.Second,
				2 * time.Second,
			},
			errs: map[string]int{},
		},

		{
			records: []Record{
				{Err: "error"},
			},
			fastest: 0,
			slowest: 0,
			average: 0,
			lats:    []time.Duration{},
			errs: map[string]int{
				"error": 1,
			},
		},

		{
			records: []Record{
				{Err: "error", Start: ts, Latency: time.Second},
				{Err: "error", Start: ts, Latency: time.Second},
			},
			fastest: 0,
			slowest: 0,
			average: 0,
			lats:    []time.Duration{},
			errs: map[string]int{
				"error": 2,
			},
		},

		{
			records: []Record{
				{Start: ts, Latency: 1 * time.Second},
				{Start: ts, Latency: 2 * time.Second},
				{Start: ts, Latency: 3 * time.Second},
				{Err: "error1"},
				{Err: "error2"},
				{Err: "error2"},
				{Err: "error3", Start: ts, Latency: time.Second},
				{Err: "error3", Start: ts, Latency: time.Second},
				{Err: "error3", Start: ts, Latency: time.Second},
			},
			fastest: 1 * time.Second,
			slowest: 3 * time.Second,
			average: (1 + 2 + 3) * time.Second / 3,
			lats: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				3 * time.Second,
			},
			errs: map[string]int{
				"error1": 1,
				"error2": 2,
				"error3": 3,
			},
		},
	}
	for i, tt := range tests {
		fastest, slowest, average, lats, errs := doStats(tt.records)
		if got, want := fastest, tt.fastest; got != want {
			t.Errorf("case%d: fastest: %v != %v", i, got, want)
		}
		if got, want := slowest, tt.slowest; got != want {
			t.Errorf("case%d: slowest: %v != %v", i, got, want)
		}
		if got, want := average, tt.average; got != want {
			t.Errorf("case%d: average: %v != %v", i, got, want)
		}
		if got, want := lats, tt.lats; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: lats: %v != %v", i, got, want)
		}
		if got, want := errs, tt.errs; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: errs: %v != %v", i, got, want)
		}
	}
}

func TestCalcQPS(t *testing.T) {
	tests := []struct {
		n   int
		d   time.Duration
		qps int
	}{
		//{n: 0, d: 0, qps: 0},
		{n: 0, d: time.Second, qps: 0},
		{n: 1, d: time.Second, qps: 1},
		{n: 100, d: time.Second, qps: 100},
		{n: 60, d: 60 * time.Second, qps: 1},
		{n: 600, d: 60 * time.Second, qps: 10},
		{n: 10001, d: 20 * time.Second, qps: 500},
		{n: 10019, d: 20 * time.Second, qps: 500},
		{n: 10020, d: 20 * time.Second, qps: 501},
	}
	for i, tt := range tests {
		if got, want := calcQPS(tt.n, tt.d), tt.qps; got != want {
			t.Errorf("case%d: %v != %v", i, got, want)
		}
	}
}

func TestResultMerge(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		a Result
		b Result
		c Result
	}{
		{
			a: Result{},
			b: Result{},
			c: Result{},
		},

		{
			a: Result{
				QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 1 * time.Second},
				},
			},
			b: Result{},
			c: Result{
				QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 1 * time.Second},
				},
			},
		},

		{
			a: Result{
				QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 1 * time.Second},
				},
			},
			b: Result{
				QPS: 2, Request: 4, Concurrent: 6, Total: 8 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 2 * time.Second},
					{Start: ts, Latency: 3 * time.Second},
				},
			},
			c: Result{
				QPS: 3, Request: 6, Concurrent: 9, Total: 8 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 1 * time.Second},
					{Start: ts, Latency: 2 * time.Second},
					{Start: ts, Latency: 3 * time.Second},
				},
			},
		},

		{
			a: Result{
				QPS: 2, Request: 4, Concurrent: 6, Total: 8 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 2 * time.Second},
					{Start: ts, Latency: 3 * time.Second},
				},
			},
			b: Result{
				QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 1 * time.Second},
				},
			},
			c: Result{
				QPS: 3, Request: 6, Concurrent: 9, Total: 8 * time.Second,
				Records: []Record{
					{Start: ts, Latency: 2 * time.Second},
					{Start: ts, Latency: 3 * time.Second},
					{Start: ts, Latency: 1 * time.Second},
				},
			},
		},
	}
	for i, tt := range tests {
		tt.a.merge(tt.b)
		if !reflect.DeepEqual(tt.a, tt.c) {
			t.Errorf("case%d: %v != %v", i, tt.a, tt.c)
		}
	}
}

func TestResults(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		list []Result
		want Results
	}{
		// case0
		{
			list: nil,
			want: nil,
		},

		// case1
		{
			list: []Result{
				{
					Name: "n1", QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
			},
			want: Results{
				{
					Name: "n1", QPS: 1, Request: 2, Concurrent: 3, Total: 4 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
			},
		},

		// case2
		{
			list: []Result{
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 1 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
				{
					Name: "n2", QPS: 2, Request: 2, Concurrent: 2, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 2 * time.Second},
					},
				},
			},
			want: Results{
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 1 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
				{
					Name: "n2", QPS: 2, Request: 2, Concurrent: 2, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 2 * time.Second},
					},
				},
			},
		},

		// case3
		{
			list: []Result{
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 1 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 2 * time.Second},
					},
				},
			},
			want: Results{
				{
					Name: "n1", QPS: 2, Request: 2, Concurrent: 2, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
						{Start: ts, Latency: 2 * time.Second},
					},
				},
			},
		},

		// case4
		{
			list: []Result{
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 1 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
					},
				},
				{
					Name: "n3", QPS: 3, Request: 3, Concurrent: 3, Total: 3 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 3 * time.Second},
					},
				},
				{
					Name: "n1", QPS: 1, Request: 1, Concurrent: 1, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 2 * time.Second},
					},
				},
			},
			want: Results{
				{
					Name: "n1", QPS: 2, Request: 2, Concurrent: 2, Total: 2 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 1 * time.Second},
						{Start: ts, Latency: 2 * time.Second},
					},
				},
				{
					Name: "n3", QPS: 3, Request: 3, Concurrent: 3, Total: 3 * time.Second,
					Records: []Record{
						{Start: ts, Latency: 3 * time.Second},
					},
				},
			},
		},
	}
	for i, tt := range tests {
		var results Results
		for _, r := range tt.list {
			results.AddResult(r)
		}
		if got, want := results, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("case%d: %v != %v", i, got, want)
		}
	}
}
