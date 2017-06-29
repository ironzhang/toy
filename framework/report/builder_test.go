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

func ErrorRecords(n int) []Record {
	errs := []string{"timeout", "eof"}
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i] = Record{Err: errs[rand.Int()%len(errs)]}
	}
	return records
}

func TestBuilder(t *testing.T) {
	reports := []Report{
		{
			Name:       "Connect",
			Total:      10 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    append(RandomRecords(200), ErrorRecords(100)...),
		},
		{
			Name:       "Connect",
			Total:      12 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    RandomRecords(200),
		},
		{
			Name:       "Publish",
			Total:      10 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    RandomRecords(5000),
		},
		{
			Name:       "Disconnect",
			Total:      10 * time.Minute,
			Concurrent: 2,
			Request:    500,
			QPS:        1000,
			Records:    ErrorRecords(100),
		},
	}
	b := Builder{
		Template:   "./templates/report.template",
		OutputDir:  "./output",
		SampleSize: 500,
	}
	if err := b.Build(reports); err != nil {
		t.Fatal(err)
	}
}
