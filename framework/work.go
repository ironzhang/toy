package framework

import "context"

type Work struct {
	Robots     []Robot
	Schedulers []Scheduler
}

func (w *Work) Run(ctx context.Context) {
	for _, s := range w.Schedulers {
		s.Run(ctx, w.Robots)
	}
}
