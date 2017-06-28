package report

import (
	"os"
	"testing"
	"time"
)

func TestHTML(t *testing.T) {
	h1 := html{
		Name:        "Connect",
		Total:       100 * time.Second,
		Slowest:     2 * time.Second,
		Fastest:     time.Millisecond,
		Average:     5 * time.Millisecond,
		Concurrent:  10,
		Request:     10000,
		RealRequest: 10000,
		QPS:         1000,
		RealQPS:     999,
		LatencyImg:  "output.png",
		Latencys: []latency{
			{10, time.Millisecond},
			{20, 2 * time.Millisecond},
			{99, 5 * time.Millisecond},
		},
	}
	h2 := html{
		Name:        "Subscribe",
		Total:       40 * time.Second,
		Slowest:     time.Second,
		Fastest:     time.Millisecond,
		Average:     5 * time.Millisecond,
		Concurrent:  10,
		Request:     10000,
		RealRequest: 10000,
		QPS:         1000,
		RealQPS:     999,
		LatencyImg:  "output.png",
		Latencys: []latency{
			{10, time.Millisecond},
			{20, 2 * time.Millisecond},
			{30, 3 * time.Millisecond},
			{99, 5 * time.Millisecond},
		},
		Errs: map[string]int{
			"EOF":     2,
			"Timeout": 10,
		},
	}

	if err := executeTemplate("./report.template", "./report.html", []html{h1, h2}); err != nil {
		t.Fatal(err)
	}

	os.Remove("./report.html")
}
