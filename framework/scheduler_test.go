package framework

import (
	"context"
	"io/ioutil"
	"sync/atomic"
	"testing"
	"time"
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
	r := &SchedulerRobot{ok: true}
	robots := []Robot{r}
	e := Scheduler{
		N: 1000,
		C: 10,
	}
	e.Run(context.Background(), robots)

	if int(r.count) != e.N {
		t.Errorf("count(%d) != %d", r.count, e.N)
	}
}

func TestSchedulerQPS(t *testing.T) {
	r := &SchedulerRobot{ok: true}
	robots := []Robot{r}

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
	e.Run(context.Background(), robots)
}

func TestSchedulerFunctions(t *testing.T) {
	r1 := &SchedulerRobot{ok: true}
	r2 := &SchedulerRobot{ok: true}
	robots := []Robot{r1, r2}

	e := Scheduler{
		N: 5,
		C: 2,
	}
	e.Run(context.Background(), robots)

	if r1.count != 3 {
		t.Errorf("c1(%d) != 3", r1.count)
	}
	if r2.count != 2 {
		t.Errorf("c2(%d) != 2", r2.count)
	}
}

func TestSchedulerDisplay(t *testing.T) {
	r1 := &SchedulerRobot{ok: true}
	robots := []Robot{r1}

	e := Scheduler{
		N:           100,
		C:           2,
		QPS:         10,
		Name:        "TestSchedulerDisplay",
		Display:     true,
		PrintReport: true,
		W:           ioutil.Discard,
	}
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	e.Run(ctx, robots)
}

func TestSchedulerRobotOK(t *testing.T) {
	r1 := &SchedulerRobot{ok: true}
	r2 := &SchedulerRobot{ok: false}
	robots := []Robot{r1, r2}

	e := Scheduler{
		N: 5,
		C: 2,
		//PrintReport: true,
	}
	e.Run(context.Background(), robots)

	if r1.count != 3 {
		t.Errorf("c1(%d) != 3", r1.count)
	}
	if r2.count != 0 {
		t.Errorf("c2(%d) != 0", r2.count)
	}
}
