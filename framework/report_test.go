package framework

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestMakeReport(t *testing.T) {
	tests := []struct {
		name      string
		request   int
		qps       int
		total     time.Duration
		durations []time.Duration
		wants     []time.Duration
	}{
		{
			name:      "test0",
			request:   10,
			qps:       100,
			total:     10 * time.Second,
			durations: []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 100 * time.Millisecond, 200 * time.Millisecond},
			wants:     []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 1 * time.Second, 2 * time.Second, 3 * time.Second},
		},
	}

	for _, test := range tests {
		var sum time.Duration
		results := make([]result, 0, len(test.durations))
		for _, d := range test.durations {
			sum += d
			results = append(results, result{duration: d})
		}

		r := makeReport(test.name, test.request, test.qps, test.total, results)
		if r.name != test.name {
			t.Errorf("name: %s != %s", r.name, test.name)
		}
		if r.request != test.request {
			t.Errorf("request: %d != %d", r.request, test.request)
		}
		if r.qps != test.qps {
			t.Errorf("qps: %d != %d", r.qps, test.qps)
		}
		if r.total != test.total {
			t.Errorf("total: %s != %s", r.total, test.total)
		}
		average := sum / time.Duration(len(test.durations))
		if r.average != average {
			t.Errorf("average: %s != %s", r.average, average)
		}
		if !reflect.DeepEqual(r.lats, test.wants) {
			t.Errorf("lats: %s != %s", r.lats, test.wants)
		}
	}
}

func TestReportPrint(t *testing.T) {
	r := report{
		name:    "TestReportPrint",
		total:   10 * time.Millisecond,
		average: 2 * time.Millisecond,
		lats:    []time.Duration{time.Millisecond, time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond, 3 * time.Millisecond},
		errs:    map[string]int{"network error": 1, "io timeout": 2},
	}
	r.print(ioutil.Discard)
}

func TestPrintLatencies(t *testing.T) {
	n := 100
	lats := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		lats[i] = time.Duration(i+1) * time.Millisecond
	}
	r := report{lats: lats}
	//r.printLatencies(os.Stdout)
	r.printLatencies(ioutil.Discard)
}
