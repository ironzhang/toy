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

	Display     bool
	PrintReport bool
	W           io.Writer
}

func (s *Scheduler) writer() io.Writer {
	if s.W == nil {
		return os.Stdout
	}
	return s.W
}

func (s *Scheduler) Run(ctx context.Context, robots []Robot) {
	resultc := make(chan result, s.N)
	robotc := make(chan Robot, s.C)

	if s.Display {
		go display(ctx, s.writer(), s.Name, resultc)
	}

	var wg sync.WaitGroup
	wg.Add(s.C)

	for i := 0; i < s.C; i++ {
		go func() {
			runWorker(s.Name, robotc, resultc)
			wg.Done()
		}()
	}

	var throttle <-chan time.Time
	if s.QPS > 0 {
		t := time.NewTicker(time.Second / time.Duration(s.QPS))
		defer t.Stop()
		throttle = t.C
	}

	n := len(robots)
	start := time.Now()
L:
	for i := 0; i < s.N; i++ {
		select {
		case <-ctx.Done():
			break L
		default:
			r := robots[i%n]
			if r.OK() {
				if s.QPS > 0 {
					<-throttle
				}
				robotc <- r
			}
		}
	}

	close(robotc)
	wg.Wait()
	close(resultc)

	if s.PrintReport {
		makeReport(s.Name, s.N, s.QPS, time.Since(start), resultc).print(s.writer())
	}
}

func runWorker(name string, robotc chan Robot, resultc chan<- result) {
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

func display(ctx context.Context, w io.Writer, name string, resultc chan result) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	var prev int
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n := len(resultc)
			if prev < n {
				prev = n
				fmt.Fprintf(w, "%s: %d requests done.\n", name, n)
			}
		}
	}
}
