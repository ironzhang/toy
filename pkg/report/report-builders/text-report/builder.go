package text_report

import (
	"io"

	"github.com/ironzhang/toy/pkg/report"
)

type Builder struct {
	w io.Writer
}

func (b *Builder) Build(rs ...report.Result) error {
	for _, r := range rs {
		b.printStats(r.Stats())
	}
	return nil
}

func (b *Builder) printStats(s report.Stats) {
}
