package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/wcharczuk/go-chart"
)

func OutputHTML(templateFile string, outdir string, reports []Report) error {
	data := make([]*report, 0, len(reports))
	for _, r := range reports {
		img, err := renderLatencyImage(fmt.Sprintf("%s/%s.png", outdir, r.Name), r.Records)
		if err != nil {
			return err
		}
		r := makeReport(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records)
		r.LatencyImg = img
		data = append(data, r)
	}
	return renderTemplate(templateFile, fmt.Sprintf("%s/report.html", outdir), data)
}

func renderLatencyImage(filename string, records []Record) (string, error) {
	if len(records) <= 1 {
		return "", nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

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
