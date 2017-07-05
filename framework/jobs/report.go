package jobs

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/report/builders/html-report"
	"github.com/ironzhang/toy/framework/report/builders/text-report"
)

type builder interface {
	Build(rs ...report.Result) error
}

type ReportJob struct {
	ResultFiles []string
	Format      string
	Template    string
	OutputDir   string
	SampleSize  int
}

func (p *ReportJob) builder() builder {
	switch p.Format {
	case "html":
		return &html_report.Builder{Template: p.Template, OutputDir: p.OutputDir, SampleSize: p.SampleSize}
	default:
		return &text_report.Builder{W: os.Stdout}
	}
}

func (p *ReportJob) Execute() error {
	rs, err := load(p.ResultFiles)
	if err != nil {
		return err
	}
	return p.builder().Build(rs...)
}

func load(filenames []string) (rs report.Results, err error) {
	for _, filename := range filenames {
		if rs, err = loadResults(filename, rs); err != nil {
			return nil, fmt.Errorf("load results from %q: %v", filename, err)
		}
	}
	return rs, nil
}

func loadResults(filename string, rs report.Results) (report.Results, error) {
	f, err := os.Open(filename)
	if err != nil {
		return rs, err
	}
	defer f.Close()

	dec := report.NewDecoder(f)
	for {
		r, err := loadResult(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return rs, err
		}
		rs.AddResult(r)
	}
	return rs, nil
}

func loadResult(dec report.Decoder) (report.Result, error) {
	var (
		err     error
		header  report.Header
		block   report.Block
		last    time.Time
		total   time.Duration
		records []report.Record
	)

	if err = dec.DecodeHeader(&header); err != nil {
		return report.Result{}, err
	}
	for {
		if err = dec.DecodeBlock(&block); err != nil {
			fmt.Println(err)
			return report.Result{}, err
		}
		if block.Time.IsZero() {
			break
		}
		last = block.Time
		records = append(records, block.Records...)
	}
	if last.After(header.Time) {
		total = last.Sub(header.Time)
	}
	return report.Result{
		Name:       header.Name,
		QPS:        header.QPS,
		Request:    header.Request,
		Concurrent: header.Concurrent,
		Total:      total,
		Records:    records,
	}, nil
}
