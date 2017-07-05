package html_report

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/report/builders/tests"
)

func MakeTestHistogram() (buckets []time.Duration, counts []int) {
	buckets = make([]time.Duration, 11)
	counts = make([]int, 11)

	rand.Seed(time.Now().Unix())
	for i := 0; i < len(buckets); i++ {
		buckets[i] = time.Duration(i+1) * time.Second
		counts[i] = rand.Int() % 10000
	}
	return
}

func TestDrawHistogram(t *testing.T) {
	filename := "histogram.png"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove(filename)
	}()

	buckets, counts := MakeTestHistogram()
	if err = drawHistogram(f, buckets, counts); err != nil {
		t.Errorf("draw histogram: %v", err)
	}
}

func MakeTestTimeSeries(n int) report.TimeSeries {
	ts := time.Now().Unix()
	series := make(report.TimeSeries, n)
	for i := 0; i < n; i++ {
		series[i].Timestamp = ts + int64(i)
		series[i].MinLatency = time.Duration(tests.Random(1, 50)) * time.Millisecond
		series[i].MaxLatency = time.Duration(tests.Random(100, 200)) * time.Millisecond
		series[i].AvgLatency = time.Duration(tests.Random(50, 100)) * time.Millisecond
		series[i].Throughput = int64(tests.Random(1000, 2000))
	}
	return series
}

func TestDrawLatencies(t *testing.T) {
	filename := "latencies.png"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove(filename)
	}()

	series := MakeTestTimeSeries(100)
	if err = drawLatencies(f, series); err != nil {
		t.Errorf("draw latencies: %v", err)
	}
}

func TestDrawThroughputs(t *testing.T) {
	filename := "throughputs.png"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove(filename)
	}()

	series := MakeTestTimeSeries(100)
	if err = drawThroughputs(f, series); err != nil {
		t.Errorf("draw throughputs: %v", err)
	}
}
