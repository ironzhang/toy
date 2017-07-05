package benchmark

import (
	"context"
	"sync"
	"time"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/robot"
)

type task struct {
	name  string
	robot robot.Robot
}

func errorstr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func (t *task) execute() report.Record {
	s := time.Now()
	e := t.robot.Do(t.name)
	return report.Record{Err: errorstr(e), Start: s.UTC(), Latency: time.Since(s)}
}

// Scheduler 调度器
type Scheduler struct {
	Name string
	N    int
	C    int
	QPS  int
}

func (s *Scheduler) Infinite() bool {
	return s.N < 0
}

func (s *Scheduler) throttleCycle() (d time.Duration, c int) {
	if s.QPS <= 0 {
		return
	}
	c = 1
	d = time.Second / time.Duration(s.QPS)
	for d < time.Millisecond {
		c *= 10
		d *= 10
	}
	return
}

func (s *Scheduler) dispatchTasks(ctx context.Context, robots []robot.Robot, taskc chan<- task) {
	var throttle <-chan time.Time
	d, cycle := s.throttleCycle()
	if cycle > 0 {
		t := time.NewTicker(d)
		defer t.Stop()
		throttle = t.C
	}

	n := len(robots)
	N := n * s.N
	for i := 0; i < N || s.Infinite(); i++ {
		select {
		case <-ctx.Done():
			return
		default:
			r := robots[i%n]
			if r.OK() {
				if cycle > 0 && i%cycle == 0 {
					<-throttle
				}
				taskc <- task{name: s.Name, robot: r}
			}
		}
	}
}

func (s *Scheduler) produceTasks(ctx context.Context, robots []robot.Robot) <-chan task {
	taskc := make(chan task, s.C)
	go func() {
		s.dispatchTasks(ctx, robots, taskc)
		close(taskc)
	}()
	return taskc
}

func runWorker(taskc <-chan task, recordc chan<- report.Record) {
	for t := range taskc {
		recordc <- t.execute()
	}
}

func (s *Scheduler) runWorkers(taskc <-chan task) <-chan report.Record {
	recordc := make(chan report.Record, s.C)

	var wg sync.WaitGroup
	wg.Add(s.C)

	for i := 0; i < s.C; i++ {
		go func() {
			runWorker(taskc, recordc)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(recordc)
	}()

	return recordc
}

func (s *Scheduler) Run(ctx context.Context, robots []robot.Robot) <-chan report.Record {
	return s.runWorkers(s.produceTasks(ctx, robots))
}
