package report

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestMakeReport(t *testing.T) {
	tests := []struct {
		name       string
		request    int
		concurrent int
		qps        int
		total      time.Duration
		durations  []time.Duration
		wants      []time.Duration
	}{
		{
			name:       "test0",
			request:    10,
			concurrent: 2,
			qps:        100,
			total:      10 * time.Second,
			durations:  []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 100 * time.Millisecond, 200 * time.Millisecond},
			wants:      []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 1 * time.Second, 2 * time.Second, 3 * time.Second},
		},
	}

	for _, test := range tests {
		var sum time.Duration
		results := make([]Record, 0, len(test.durations))
		for _, d := range test.durations {
			sum += d
			results = append(results, Record{Elapse: d})
		}

		r := makeReport(test.name, test.request, test.concurrent, test.qps, test.total, results)
		if r.Name != test.name {
			t.Errorf("name: %s != %s", r.Name, test.name)
		}
		if r.Request != test.request {
			t.Errorf("request: %d != %d", r.Request, test.request)
		}
		if r.Concurrent != test.concurrent {
			t.Errorf("concurrent: %d != %d", r.Concurrent, test.concurrent)
		}
		if r.QPS != test.qps {
			t.Errorf("qps: %d != %d", r.QPS, test.qps)
		}
		if r.Total != test.total {
			t.Errorf("total: %s != %s", r.Total, test.total)
		}
		average := sum / time.Duration(len(test.durations))
		if r.Average != average {
			t.Errorf("average: %s != %s", r.Average, average)
		}
		if !reflect.DeepEqual(r.lats, test.wants) {
			t.Errorf("lats: %s != %s", r.lats, test.wants)
		}
	}
}

func TestReportPrint(t *testing.T) {
	r := report{
		Name:    "TestReportPrint",
		Total:   10 * time.Millisecond,
		Average: 2 * time.Millisecond,
		lats:    []time.Duration{time.Millisecond, time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond, 3 * time.Millisecond},
		Errs:    map[string]int{"network error": 1, "io timeout": 2},
	}
	r.print(ioutil.Discard)
}

func TestPrintLatencies(t *testing.T) {
	n := 100
	lats := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		lats[i] = time.Duration(i+1) * time.Millisecond
	}
	r := report{Latencies: latencyDistribution(lats)}
	//r.printLatencies(os.Stdout)
	r.printLatencies(ioutil.Discard)
}

func TestReport(t *testing.T) {
	r := Report{
		Name:       "TestReport",
		Total:      10 * time.Second,
		Concurrent: 2,
		Request:    10,
		QPS:        100,
		Records: []Record{
			{Elapse: 5 * time.Second},
			{Elapse: 1 * time.Second},
			{Elapse: 2 * time.Second},
			{Elapse: 3 * time.Second},
			{Elapse: 4 * time.Second},
		},
	}
	//r.Print(os.Stdout)
	r.Print(ioutil.Discard)
}
