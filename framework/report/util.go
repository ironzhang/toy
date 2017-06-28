package report

import (
	"html/template"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

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

func renderRecords(filename string, records []Record) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
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
	return graph.Render(chart.PNG, f)
}

func fastestRecord(records []Record) Record {
	var x Record
	for i, r := range records {
		if i == 0 {
			x = r
		} else if x.Elapse > r.Elapse {
			x = r
		}
	}
	return x
}
