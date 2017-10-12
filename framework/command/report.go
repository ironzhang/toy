package command

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/report/builders/text-report"
)

type builder interface {
	Build(rs ...report.Result) error
}

type ReportCmd struct {
	format     string
	outputDir  string
	sampleSize int

	resultFiles []string
}

func (c *ReportCmd) Run(args []string) error {
	if err := c.parse(args); err != nil {
		return nil
	}
	return c.execute()
}

func (c *ReportCmd) parse(args []string) error {
	var fs flag.FlagSet
	fs.Usage = func() {
		fmt.Print("Usage: toy report [OPTIONS] FILE [FILE...]\n\n")
		fmt.Print("make benchmark report with test records\n\n")
		fs.PrintDefaults()
	}
	fs.StringVar(&c.format, "format", "text", "report format, text/html")
	fs.StringVar(&c.outputDir, "output-dir", "output", "output dir")
	fs.IntVar(&c.sampleSize, "sample-size", 500, "sample size")
	if err := fs.Parse(args); err != nil {
		return err
	}
	c.resultFiles = fs.Args()
	if len(c.resultFiles) <= 0 {
		fs.Usage()
		os.Exit(1)
	}
	return nil
}

func (c *ReportCmd) execute() error {
	rs, err := load(c.resultFiles)
	if err != nil {
		return err
	}
	return c.builder().Build(rs...)
}

func (c *ReportCmd) builder() builder {
	return &text_report.Builder{W: os.Stdout}
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
		last    time.Time
		total   time.Duration
		records []report.Record
	)

	if err = dec.DecodeHeader(&header); err != nil {
		return report.Result{}, err
	}
	for {
		var block report.Block // block变量必须在for循环内
		if err = dec.DecodeBlock(&block); err != nil {
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
