package text_report

import (
	"fmt"
	"io"
	"strings"

	"github.com/ironzhang/toy/pkg/report"
)

type Builder struct {
	w io.Writer
}

func (b *Builder) Build(rs ...report.Result) error {
	for _, r := range rs {
		printStats(b.w, r.Stats())
	}
	return nil
}

func printStats(w io.Writer, s *report.Stats) {
	fmt.Fprintf(w, "\nSummary: %s\n", s.Name)
	fmt.Fprintf(w, "  Total:\t%s\n", s.Total)
	fmt.Fprintf(w, "  Slowest:\t%s\n", s.Slowest)
	fmt.Fprintf(w, "  Fastest:\t%s\n", s.Fastest)
	fmt.Fprintf(w, "  Average:\t%s\n", s.Average)
	fmt.Fprintf(w, "  Concurrent:\t%d\n", s.Concurrent)
	fmt.Fprintf(w, "  Requests:\t%d/%d\n", s.RealRequest, s.Request)
	fmt.Fprintf(w, "  Requests/sec:\t%4.4f/%d\n", s.RealQPS, s.QPS)
	if len(s.Lats) > 0 {
		printHistogram(w, s)
	}
	if len(s.Errs) > 0 {
	}
	fmt.Fprintln(w)
}

const barChar = "∎"

func printHistogram(w io.Writer, s *report.Stats) {
	var barLen int
	buckets, counts, max := s.Histogram()
	fmt.Fprintf(w, "\nResponse time histogram:\n")
	for i := 0; i < len(buckets); i++ {
		if max > 0 {
			barLen = (counts[i]*40 + max/2) / max
		}
		fmt.Fprintf(w, "  %s [%d]\t|%s\n", buckets[i], counts[i], strings.Repeat(barChar, barLen))
	}
}

func printLatencies(w io.Writer, s *report.Stats) {
	pcs, data := s.Latencies()
	fmt.Fprintf(w, "\nLatency distribution:\n")
	for i := 0; i < len(pcs); i++ {
		fmt.Fprintf(w, "  %v%% in %s\n", pcs[i], data[i])
	}
}
