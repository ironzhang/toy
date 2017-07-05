package scheduler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ironzhang/toy/framework/codec"
	"github.com/ironzhang/toy/framework/robot"
)

type Work struct {
	Ask        bool
	Encoder    codec.Encoder
	Robots     []robot.Robot
	Schedulers []Scheduler
}

func (w *Work) Run() {
	for _, s := range w.Schedulers {
		if s.N != 0 {
			w.schedule(&s)
		}
	}
}

func (w *Work) schedule(s *Scheduler) {
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
	s.Run(ctx, w.Robots, w.Encoder)
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
