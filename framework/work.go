package framework

import (
	"context"

	"github.com/ironzhang/toy/framework/robot"
)

type Work struct {
	Robots     []robot.Robot
	Schedulers []Scheduler
}

func (w *Work) Run(ctx context.Context) {
	for _, s := range w.Schedulers {
		if s.N != 0 {
			s.Run(ctx, w.Robots)
		}
	}
}
