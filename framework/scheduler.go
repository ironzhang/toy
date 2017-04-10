package framework

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type result struct {
	err      error
	duration time.Duration
}

type Robot interface {
	OK() bool
	Do(name string) error
}

type Scheduler struct {
	N    int
	C    int
	QPS  int
	Name string

	Sample      int
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

func (s *Scheduler) Run(ctx context.Context, robots []Robot) {
	start := time.Now()
	resultc := s.runWorkers(s.produceTasks(ctx, robots))

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
	results := make([]result, 0, nres)
	for res := range resultc {
		done++
		if s.Display && time.Since(prev) >= 500*time.Millisecond {
			prev = time.Now()
			fmt.Fprintf(s.writer(), "%s: %d requests done.\n", s.Name, done)
		}

		results = append(results, res)
		if len(results) >= nres {
			if s.PrintReport {
				makeReport(s.Name, request, s.QPS, time.Since(start), results).print(s.writer())
			}
			start = time.Now()
			results = results[:0]
		}
	}

	if s.PrintReport && len(results) > 0 {
		makeReport(s.Name, request, s.QPS, time.Since(start), results).print(s.writer())
	}
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

func (s *Scheduler) dispatchTasks(ctx context.Context, robots []Robot, robotc chan<- Robot) {
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
				robotc <- r
			}
		}
	}
}

func (s *Scheduler) produceTasks(ctx context.Context, robots []Robot) <-chan Robot {
	robotc := make(chan Robot, s.C)
	go func() {
		s.dispatchTasks(ctx, robots, robotc)
		close(robotc)
	}()
	return robotc
}

func (s *Scheduler) runWorkers(robotc <-chan Robot) <-chan result {
	resultc := make(chan result, s.C)

	var wg sync.WaitGroup
	wg.Add(s.C)

	for i := 0; i < s.C; i++ {
		go func() {
			runWorker(s.Name, robotc, resultc)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultc)
	}()

	return resultc
}

func runWorker(name string, robotc <-chan Robot, resultc chan<- result) {
	for r := range robotc {
		resultc <- call(name, r)
	}
}

func call(name string, r Robot) result {
	s := time.Now()
	e := r.Do(name)
	d := time.Since(s)
	return result{err: e, duration: d}
}
