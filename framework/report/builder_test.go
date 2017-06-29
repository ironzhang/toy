package report

import (
	"math/rand"
	"testing"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomRecords(n int) []Record {
	now := time.Now()
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i] = Record{Start: now.Add(time.Duration(i) * time.Second).UTC(), Elapse: time.Duration(random(10, 500)) * time.Millisecond}
	}
	return records
}

func TestOutputHTML(t *testing.T) {
	reports := []Report{
		{
			Name:       "Connect",
			Total:      10 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    RandomRecords(5002),
		},
		{
			Name:       "Publish",
			Total:      10 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    RandomRecords(5000),
		},
	}

	b := Builder{
		Template:   "./templates/report.template",
		OutputDir:  "./output",
		SampleSize: 600,
	}
	if err := b.MakeHTML(reports); err != nil {
		t.Fatal(err)
	}
}
