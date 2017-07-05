package schedule

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ironzhang/toy/framework/robot"
)

type SchedulerRobot struct {
	ok    bool
	count int64
}

func (r *SchedulerRobot) OK() bool {
	return r.ok
}

func (r *SchedulerRobot) Do(name string) error {
	atomic.AddInt64(&r.count, 1)
	return nil
}

func TestSchedulerN(t *testing.T) {
	r1 := &SchedulerRobot{ok: true}
	r2 := &SchedulerRobot{ok: true}
	robots := []robot.Robot{r1, r2}
	e := Scheduler{
		N: 1000,
		C: 10,
	}
	recordc := e.Run(context.Background(), robots)
	for r := range recordc {
		if r.Err != "" {
			t.Errorf("scheduler run: %v", r.Err)
		}
	}

	if int(r1.count) != e.N {
		t.Errorf("count(%d) != %d", r1.count, e.N)
	}
	if int(r2.count) != e.N {
		t.Errorf("count(%d) != %d", r2.count, e.N)
	}
}

func TestSchedulerQPS1(t *testing.T) {
	r := &SchedulerRobot{ok: true}
	robots := []robot.Robot{r}

	start := time.Now()
	tick := time.Tick(time.Second)
	go func() {
		for {
			tt := <-tick
			sec := int64(tt.Sub(start) / time.Second)
			if r.count > sec {
				t.Errorf("count(%d) > sec(%d)", r.count, sec)
			}
		}
	}()

	e := Scheduler{
		N:   2,
		C:   2,
		QPS: 1,
	}
	recordc := e.Run(context.Background(), robots)
	for range recordc {
	}
}

func TestSchedulerQPS2(t *testing.T) {
	r := &SchedulerRobot{ok: true}
	robots := []robot.Robot{r}

	ch1 := (&Scheduler{
		N:   10000,
		C:   1,
		QPS: 0,
	}).Run(context.Background(), robots)
	for range ch1 {
	}

	ch2 := (&Scheduler{
		N:   10000,
		C:   1,
		QPS: 1000000,
		//PrintReport: true,
	}).Run(context.Background(), robots)
	for range ch2 {
	}
}

func TestSchedulerRobotOK(t *testing.T) {
	r1 := &SchedulerRobot{ok: true}
	r2 := &SchedulerRobot{ok: false}
	robots := []robot.Robot{r1, r2}

	e := Scheduler{
		N: 5,
		C: 2,
	}
	recordc := e.Run(context.Background(), robots)
	for range recordc {
	}

	if r1.count != 5 {
		t.Errorf("c1(%d) != 3", r1.count)
	}
	if r2.count != 0 {
		t.Errorf("c2(%d) != 0", r2.count)
	}
}
