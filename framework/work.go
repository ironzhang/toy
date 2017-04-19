package framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/ironzhang/toy/framework/robot"
)

type Work struct {
	Ask        bool
	Robots     []robot.Robot
	Schedulers []Scheduler
}

func (w *Work) Run(ctx context.Context) {
	for _, s := range w.Schedulers {
		if s.N != 0 {
			if w.Ask && !ask(s.Name) {
				continue
			}
			s.Run(ctx, w.Robots)
		}
	}
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
