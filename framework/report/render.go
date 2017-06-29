package report

import (
	"html/template"
	"io"
	"time"

	"github.com/wcharczuk/go-chart"
)

func renderLatencies(w io.Writer, records []Record) error {
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
	return graph.Render(chart.PNG, w)
}

func renderTemplate(w io.Writer, filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}
