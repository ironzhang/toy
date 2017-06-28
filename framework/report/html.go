package report

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

func OutputHTML(templateFile string, outdir string, reports []Report) (err error) {
	data := make([]*report, 0, len(reports))
	for _, r := range reports {
		if err = renderLatencyImage(fmt.Sprintf("%s/%s.png", outdir, r.Name), r.Records); err != nil {
			return err
		}
		data = append(data, makeReport(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records))
	}
	return renderTemplate(templateFile, fmt.Sprintf("%s/report.html", outdir), data)
}

func renderLatencyImage(filename string, records []Record) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(records) <= 1 {
		return nil
	}

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
	return graph.Render(chart.PNG, f)
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
