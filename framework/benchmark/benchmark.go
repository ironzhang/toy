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

const MaxRecordNumPerBlock = 5000

// Benchmark 性能测试
type Benchmark struct {
	Ask        bool
	Verbose    int
	Encoder    report.Encoder
	Robots     []robot.Robot
	Schedulers []Scheduler
}

func (b *Benchmark) Run() {
	for _, s := range b.Schedulers {
		if s.N != 0 {
			b.schedule(&s)
		}
	}
}

func (b *Benchmark) schedule(s *Scheduler) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	if b.Ask && !ask(s.Name) {
		return
	}

	b.benchmark(ctx, s)
}

func (b *Benchmark) benchmark(ctx context.Context, s *Scheduler) {
	b.writeHeader(s)

	done := 0
	prev := time.Now()
	records := make([]report.Record, 0, MaxRecordNumPerBlock)

	start := time.Now()
	recordc := s.Run(ctx, b.Robots)
	for rec := range recordc {
		done++
		if b.Verbose >= 2 && time.Since(prev) >= 500*time.Millisecond {
			prev = time.Now()
			fmt.Fprintf(os.Stdout, "%s: %d requests done.\n", s.Name, done)
		}

		records = append(records, rec)
		if len(records) >= MaxRecordNumPerBlock {
			elapse := time.Since(start)
			b.writeBlock(elapse, records)
			if b.Verbose >= 1 {
				b.printResult(s, elapse, records)
			}

			start = time.Now()
			records = records[:0]
		}
	}
	if len(records) > 0 {
		elapse := time.Since(start)
		b.writeBlock(elapse, records)
		if b.Verbose >= 1 {
			b.printResult(s, elapse, records)
		}
	}
	b.writeBlock(-1, nil) // end of result
}

func (b *Benchmark) writeHeader(s *Scheduler) {
	n := s.N
	if n > 0 {
		n *= len(b.Robots)
	}

	if b.Encoder != nil {
		header := &report.Header{
			Name:       s.Name,
			QPS:        s.QPS,
			Request:    n,
			Concurrent: s.C,
		}
		if err := b.Encoder.EncodeHeader(header); err != nil {
			log.Printf("encode header: %v", err)
		}
	}
}

func (b *Benchmark) writeBlock(total time.Duration, records []report.Record) {
	if b.Encoder != nil {
		block := &report.Block{
			Total:   total,
			Records: records,
		}
		if err := b.Encoder.EncodeBlock(block); err != nil {
			log.Printf("encode block: %v", err)
		}
	}
}

func (b *Benchmark) printResult(s *Scheduler, total time.Duration, records []report.Record) {
	n := s.N
	if n > 0 {
		n *= len(b.Robots)
	}

	result := report.Result{
		Name:       s.Name,
		QPS:        s.QPS,
		Request:    n,
		Concurrent: s.C,
		Total:      total,
		Records:    records,
	}
	(&text_report.Builder{W: os.Stdout}).Build(result)
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
