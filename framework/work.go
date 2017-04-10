package framework

import "context"

type Work struct {
	Robots     []Robot
	Schedulers []Scheduler
}

func (w *Work) Run(ctx context.Context) {
	for _, s := range w.Schedulers {
		if s.N != 0 {
			s.Run(ctx, w.Robots)
		}
	}
}
