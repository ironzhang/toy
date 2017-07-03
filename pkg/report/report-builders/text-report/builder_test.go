package text_report

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ironzhang/toy/pkg/report"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func makeTestRecord(n int) []report.Record {
	ts := time.Now()
	//rand.Seed(ts.Unix())
	records := make([]report.Record, n)
	for i := 0; i < n; i++ {
		records[i].Start = ts.Add(200 * time.Millisecond)
		records[i].Latency = time.Duration(random(1, 1200)) * time.Millisecond
		if records[i].Latency > 1*time.Second {
			records[i].Err = "timeout"
		}
	}
	return records
}

func TestBuilder(t *testing.T) {
	r := report.Result{
		Name:       "Test",
		QPS:        1000,
		Request:    10000,
		Concurrent: 10,
		Total:      10 * time.Second,
		Records:    makeTestRecord(10000),
	}
	(&Builder{W: os.Stdout}).Build(r)
}
