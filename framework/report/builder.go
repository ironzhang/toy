package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Builder struct {
	Template   string
	OutputDir  string
	SampleSize int
}

func (b *Builder) MakeHTML(reports []Report) error {
	os.MkdirAll(b.OutputDir, os.ModePerm)

	reports = mergeReports(reports)
	data := make([]*report, 0, len(reports))
	for _, r := range reports {
		img, err := renderLatencyImage(fmt.Sprintf("%s/%s.png", b.OutputDir, r.Name), r.Records, b.SampleSize)
		if err != nil {
			return err
		}
		r := makeReport(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records)
		r.LatencyImg = img
		data = append(data, r)
	}
	return renderTemplate(b.Template, fmt.Sprintf("%s/report.html", b.OutputDir), data)
}

func mergeReports(reports []Report) []Report {
	results := make([]Report, 0)
	for _, r := range reports {
		lookup := false
		for i := range results {
			if results[i].Name == r.Name {
				results[i] = mergeReport(results[i], r)
				lookup = true
				break
			}
		}
		if !lookup {
			results = append(results, r)
		}
	}
	return results
}

func mergeReport(a Report, b Report) Report {
	var c Report
	c.Name = a.Name
	if a.Total > b.Total {
		c.Total = a.Total
	} else {
		c.Total = b.Total
	}
	c.Concurrent = a.Concurrent + b.Concurrent
	c.Request = a.Request + b.Request
	c.QPS = a.QPS + b.QPS
	c.Records = append(a.Records, b.Records...)
	return c
}

func renderLatencyImage(filename string, records []Record, sampleSize int) (string, error) {
	if len(records) <= 1 {
		return "", nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sort.Slice(records, func(i, j int) bool { return records[i].Start.Before(records[j].Start) })
	records = sampling(records, sampleSize)

	n := len(records)
	xvalues := make([]time.Time, n)
	yvalues := make([]float64, n)
	for i := 0; i < n; i++ {
		xvalues[i] = records[i].Start
		yvalues[i] = float64(records[i].Elapse) / float64(time.Millisecond)
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "latency(ms)",
			NameStyle: chart.Style{Show: true},
			Style:     chart.Style{Show: true},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xvalues,
				YValues: yvalues,
			},
		},
	}
	return filepath.Base(filename), graph.Render(chart.PNG, f)
}

func renderTemplate(srcFile, dstFile string, data interface{}) error {
	t, err := template.ParseFiles(srcFile)
	if err != nil {
		return err
	}

	f, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}
