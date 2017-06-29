package report

import (
	"math"
	"os"
	"testing"
	"time"
)

func MakeTestRecords(n int) []Record {
	now := time.Now()
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i].Start = now.Add(time.Duration(i) * time.Second)
		sin := math.Sin(float64(i) / 10)
		records[i].Elapse = time.Duration(sin*float64(time.Second)) + time.Second
	}
	return records
}

func TestRenderLatencies(t *testing.T) {
	filename := "lat.png"

	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	records := MakeTestRecords(100)
	if err := renderLatencies(f, records); err != nil {
		t.Error(err)
	}

	os.Remove(filename)
}

func TestRenderTemplate(t *testing.T) {
	filename := "report.html"

	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	reports := []*report{
		makeReport("test", 200, 10, 1000, 10*time.Second, "", MakeTestRecords(100)),
	}
	if err := renderTemplate(f, "./templates/report.template", reports); err != nil {
		t.Error(err)
	}

	os.Remove(filename)
}