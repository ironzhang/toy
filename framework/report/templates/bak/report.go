package report

import (
	"io"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Record2 struct {
	Err    string
	Start  time.Time
	Elapse time.Duration
}

type Report2 struct {
	Name       string
	Total      time.Duration
	Concurrent int
	Request    int
	QPS        int
	Records    []Record
}

func (r *Report) PrintText(w io.Writer) {
	makeText(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records).print(w)
}

func (r *Report) render(w io.Writer) error {
	n := len(r.Records)
	xvalues := make([]time.Time, n)
	yvalues := make([]float64, n)
	for i := 0; i < n; i++ {
		xvalues[i] = r.Records[i].Start
		yvalues[i] = float64(r.Records[i].Elapse) / float64(time.Millisecond)
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
