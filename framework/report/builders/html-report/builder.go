package html_report

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/ironzhang/toy/framework/report"
)

type latpct struct {
	Percent float64
	Latency time.Duration
}

type targ struct {
	Name           string
	QPS            int
	RealQPS        int
	Request        int
	RealRequest    int
	Concurrent     int
	Total          time.Duration
	Fastest        time.Duration
	Slowest        time.Duration
	Average        time.Duration
	Errs           map[string]int
	Latpcts        []latpct
	HistogramImg   string
	LatenciesImg   string
	ThroughputsImg string
}

type Builder struct {
	Template   string
	OutputDir  string
	SampleSize int
}

func (b *Builder) Build(rs ...report.Result) (err error) {
	os.MkdirAll(b.OutputDir, os.ModePerm)
	args := make([]targ, len(rs))
	for i, r := range rs {
		if args[i], err = b.build(r.Stats(b.SampleSize)); err != nil {
			return err
		}
	}
	return b.renderHTML(args)
}

func (b *Builder) renderHTML(args []targ) error {
	t, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return err
	}

	filename := "report.html"
	f, err := os.Create(b.OutputDir + "/" + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, args)
}

func (b *Builder) build(s *report.Stats) (a targ, err error) {
	a.Name = s.Name
	a.QPS = s.QPS
	a.RealQPS = s.RealQPS
	a.Request = s.Request
	a.RealRequest = s.RealRequest
	a.Concurrent = s.Concurrent
	a.Total = s.Total
	a.Fastest = s.Fastest
	a.Slowest = s.Slowest
	a.Average = s.Average
	a.Errs = s.Errs
	pcs, data := s.Latencies()
	for i := 0; i < len(pcs); i++ {
		a.Latpcts = append(a.Latpcts, latpct{Percent: pcs[i], Latency: data[i]})
	}
	if a.HistogramImg, err = b.buildHistogramImage(s); err != nil {
		return a, err
	}
	if a.LatenciesImg, err = b.buildLatenciesImage(s); err != nil {
		return a, err
	}
	if a.ThroughputsImg, err = b.buildThroughputsImage(s); err != nil {
		return a, err
	}
	return a, nil
}

func (b *Builder) buildHistogramImage(s *report.Stats) (string, error) {
	buckets, counts, _ := s.Histogram()
	if len(buckets) <= 0 {
		return "", nil
	}

	filename := fmt.Sprintf("%s-histogram.png", s.Name)
	f, err := os.Create(b.OutputDir + "/" + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return filename, drawHistogram(f, buckets, counts)
}

func (b *Builder) buildLatenciesImage(s *report.Stats) (string, error) {
	if len(s.Series) <= 1 {
		return "", nil
	}

	filename := fmt.Sprintf("%s-latencies.png", s.Name)
	f, err := os.Create(b.OutputDir + "/" + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return filename, drawLatencies(f, s.Series)
}

func (b *Builder) buildThroughputsImage(s *report.Stats) (string, error) {
	if len(s.Series) <= 1 {
		return "", nil
	}

	filename := fmt.Sprintf("%s-throughputs.png", s.Name)
	f, err := os.Create(b.OutputDir + "/" + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return filename, drawThroughputs(f, s.Series)
}
