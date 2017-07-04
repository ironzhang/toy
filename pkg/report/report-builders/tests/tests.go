package tests

import (
	"math/rand"
	"time"

	"github.com/ironzhang/toy/pkg/report"
)

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func MakeRandomRecords(n int) []report.Record {
	ts := time.Now()
	//rand.Seed(ts.Unix())
	records := make([]report.Record, n)
	for i := 0; i < n; i++ {
		ts = ts.Add(time.Duration(Random(10, 100)) * time.Millisecond)
		records[i].Start = ts
		records[i].Latency = time.Duration(Random(1, 1200)) * time.Millisecond
		if records[i].Latency > 1*time.Second {
			records[i].Err = "timeout"
		}
	}
	return records
}

func MakeTestResult(name string) report.Result {
	return report.Result{
		Name:       name,
		QPS:        1000,
		Request:    10000,
		Concurrent: 10,
		Total:      10 * time.Second,
		Records:    MakeRandomRecords(10000),
	}
}
