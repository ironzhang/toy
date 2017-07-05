package benchmark

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/report/builders/text-report"
	"github.com/ironzhang/toy/framework/robot"
)

const maxRecordNumPerResult = 5000

// Benchmark 性能测试
type Benchmark struct {
	Ask        bool
	Verbose    int
	Encoder    report.Encoder
	Robots     []robot.Robot
	Schedulers []Scheduler
}

func (w *Benchmark) Run() {
	for _, s := range w.Schedulers {
		if s.N != 0 {
			w.schedule(&s)
		}
	}
}

func (w *Benchmark) schedule(s *Scheduler) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	if w.Ask && !ask(s.Name) {
		return
	}

	w.benchmark(ctx, s)
}

func (w *Benchmark) benchmark(ctx context.Context, s *Scheduler) {
	w.writeHeader(s)

	done := 0
	prev := time.Now()
	records := make([]report.Record, 0, maxRecordNumPerResult)

	start := time.Now()
	recordc := s.Run(ctx, w.Robots)
	for rec := range recordc {
		done++
		if w.Verbose >= 2 && time.Since(prev) >= 500*time.Millisecond {
			prev = time.Now()
			fmt.Fprintf(os.Stdout, "%s: %d requests done.\n", s.Name, done)
		}

		records = append(records, rec)
		if len(records) >= maxRecordNumPerResult {
			elapse := time.Since(start)
			w.writeBlock(elapse, records)
			if w.Verbose >= 1 {
				w.printResult(s, elapse, records)
			}

			start = time.Now()
			records = records[:0]
		}
	}
	if len(records) > 0 {
		elapse := time.Since(start)
		w.writeBlock(elapse, records)
		if w.Verbose >= 1 {
			w.printResult(s, elapse, records)
		}
	}
}

func (w *Benchmark) writeHeader(s *Scheduler) {
	if w.Encoder != nil {
		header := &report.Header{
			Name:       s.Name,
			QPS:        s.QPS,
			Request:    s.N * len(w.Robots),
			Concurrent: s.C,
		}
		if err := w.Encoder.EncodeHeader(header); err != nil {
			log.Printf("encode header: %v", err)
		}
	}
}

func (w *Benchmark) writeBlock(total time.Duration, records []report.Record) {
	if w.Encoder != nil {
		block := &report.Block{
			Total:   total,
			Records: records,
		}
		if err := w.Encoder.EncodeBlock(block); err != nil {
			log.Printf("encode block: %v", err)
		}
	}
}

func (w *Benchmark) printResult(s *Scheduler, total time.Duration, records []report.Record) {
	result := report.Result{
		Name:       s.Name,
		QPS:        s.QPS,
		Request:    s.N * len(w.Robots),
		Concurrent: s.C,
		Total:      total,
		Records:    records,
	}
	b := text_report.Builder{W: os.Stdout}
	b.Build(result)
}

func ask(name string) bool {
	var err error
	var answer string
	for {
		fmt.Printf("execute %s scheduler[yes/no]?", name)
		if _, err = fmt.Scan(&answer); err != nil {
			break
		}

		answer = strings.ToLower(answer)
		if answer == "yes" || answer == "y" {
			return true
		} else if answer == "no" || answer == "n" {
			return false
		} else {
			fmt.Printf("unknown answer: %s\n", answer)
		}
	}
	return false
}
