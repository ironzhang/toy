package benchmark

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/robot"
)

type WorkRobot struct {
	ConnectCount    int64
	PrepareCount    int64
	PublishCount    int64
	DisconnectCount int64
}

func (r *WorkRobot) OK() bool {
	return true
}

func (r *WorkRobot) Do(name string) error {
	switch name {
	case "Connect":
		atomic.AddInt64(&r.ConnectCount, 1)
	case "Prepare":
		atomic.AddInt64(&r.PrepareCount, 1)
	case "Publish":
		atomic.AddInt64(&r.PublishCount, 1)
	case "Disconnect":
		atomic.AddInt64(&r.DisconnectCount, 1)
	}
	return nil
}

func TestBenchmark(t *testing.T) {
	n := 100
	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &WorkRobot{})
	}

	filename := "benchmark.tbr"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b := Benchmark{
		Ask:     false,
		Verbose: 0,
		Encoder: report.NewGobEncoder(f),
		Robots:  robots,
		Schedulers: []Scheduler{
			{
				N:    1,
				C:    10,
				QPS:  1000,
				Name: "Connect",
			},
			{
				N:    1,
				C:    10,
				QPS:  1000,
				Name: "Prepare",
			},
			{
				N:    100,
				C:    100,
				QPS:  5000,
				Name: "Publish",
			},
			{
				N:    1,
				C:    10,
				QPS:  1000,
				Name: "Disconnect",
			},
		},
	}
	b.Run()

	for _, r := range robots {
		wr := r.(*WorkRobot)
		if wr.ConnectCount != 1 {
			t.Errorf("ConnectCount: %d != 1", wr.ConnectCount)
		}
		if wr.PrepareCount != 1 {
			t.Errorf("PrepareCount: %d != 1", wr.PrepareCount)
		}
		if wr.PublishCount != 100 {
			t.Errorf("PublishCount: %d != 100", wr.PublishCount)
		}
		if wr.DisconnectCount != 1 {
			t.Errorf("DisconnectCount: %d != 1", wr.DisconnectCount)
		}
	}

	os.Remove(filename)
}

func ExampleAsk() {
	if ask("TestAsk") {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
