package report

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func NewTestReport() *Report {
	n := 100
	now := time.Now()
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i] = Record{Start: now.Add(time.Duration(i) * time.Second), Elapse: time.Duration(random(10, 500)) * time.Millisecond}
	}
	return &Report{
		Name:       "Connect",
		Total:      100 * time.Second,
		Concurrent: 10,
		Request:    100,
		QPS:        0,
		Records:    records,
	}
}

func TestRender(t *testing.T) {
	r := NewTestReport()

	f, err := os.Create("output.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err = r.render(f); err != nil {
		t.Fatal(err)
	}

	os.Remove("output.png")
}
