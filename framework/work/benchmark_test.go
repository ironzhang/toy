package work

import (
	"fmt"
	"sync/atomic"
	"testing"

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

func TestWork(t *testing.T) {
	n := 100
	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &WorkRobot{})
	}

	w := Work{
		Ask:    false,
		Robots: robots,
		Schedulers: []Scheduler{
			{
				N:           1,
				C:           10,
				QPS:         1000,
				Name:        "Connect",
				Display:     false,
				PrintReport: false,
			},
			{
				N:           1,
				C:           10,
				QPS:         1000,
				Name:        "Prepare",
				Display:     false,
				PrintReport: false,
			},
			{
				N:      100,
				C:      100,
				QPS:    5000,
				Name:   "Publish",
				Sample: 100,
				//Display:     true,
				//PrintReport: true,
			},
			{
				N:           1,
				C:           10,
				QPS:         1000,
				Name:        "Disconnect",
				Display:     false,
				PrintReport: false,
			},
		},
	}
	w.Run()

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
}

func ExampleAsk() {
	if ask("TestAsk") {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}