package framework

import (
	"context"
	"fmt"
	"io"
	"os"
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
	d := time.Since(s)
	return report.Record{Err: errorstr(e), Start: s.UTC(), Elapse: d}
}

type Scheduler struct {
	Name   string
	N      int
	C      int
	QPS    int
	Sample int

	Display     bool
	PrintReport bool

	W io.Writer `json:"-"`
}

func (s *Scheduler) writer() io.Writer {
	if s.W == nil {
		return os.Stdout
	}
	return s.W
}

func (s *Scheduler) infinite() bool {
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
	for i := 0; i < N || s.infinite(); i++ {
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

func (s *Scheduler) Run(ctx context.Context, robots []robot.Robot) {
	start := time.Now()
	recordc := s.runWorkers(s.produceTasks(ctx, robots))

	var nres int
	var request int
	if s.infinite() {
		nres = s.Sample * len(robots)
		request = -1
	} else {
		nres = s.N * len(robots)
		request = nres
	}

	done := 0
	prev := start
	records := make([]report.Record, 0, nres)
	for res := range recordc {
		done++
		if s.Display && time.Since(prev) >= 500*time.Millisecond {
			prev = time.Now()
			fmt.Fprintf(s.writer(), "%s: %d requests done.\n", s.Name, done)
		}

		records = append(records, res)
		if len(records) >= nres {
			if s.PrintReport {
				(&report.Report{
					Name:       s.Name,
					Total:      time.Since(start),
					Concurrent: s.C,
					Request:    request,
					QPS:        s.QPS,
					Records:    records,
				}).Print(s.writer())
			}
			start = time.Now()
			records = records[:0]
		}
	}

	if s.PrintReport && len(records) > 0 {
		(&report.Report{
			Name:       s.Name,
			Total:      time.Since(start),
			Concurrent: s.C,
			Request:    request,
			QPS:        s.QPS,
			Records:    records,
		}).Print(s.writer())
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

func runWorker(taskc <-chan task, recordc chan<- report.Record) {
	for t := range taskc {
		recordc <- t.execute()
	}
}
