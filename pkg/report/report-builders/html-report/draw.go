package html_report

import (
	"fmt"
	"io"
	"time"

	"github.com/ironzhang/toy/pkg/report"
	"github.com/wcharczuk/go-chart"
)

func drawHistogram(w io.Writer, buckets []time.Duration, counts []int) error {
	values := make([]chart.Value, len(buckets))
	for i := 0; i < len(values); i++ {
		values[i].Value = float64(counts[i])
		values[i].Label = fmt.Sprintf("%s [%d]", buckets[i], counts[i])
	}

	graph := chart.BarChart{
		Width:  1280,
		Height: 720,
		XAxis:  chart.Style{Show: true},
		Bars:   values,
	}
	return graph.Render(chart.PNG, w)
}

func drawLatencies(w io.Writer, series report.TimeSeries) error {
	n := len(series)
	xvalues := make([]time.Time, n)
	maxLatencies := make([]float64, n)
	minLatencies := make([]float64, n)
	avgLatencies := make([]float64, n)
	for i, p := range series {
		xvalues[i] = time.Unix(p.Timestamp, 0)
		maxLatencies[i] = p.MaxLatency.Seconds() * 1000
		minLatencies[i] = p.MinLatency.Seconds() * 1000
		avgLatencies[i] = p.AvgLatency.Seconds() * 1000
	}

	graph := chart.Chart{
		Width:      1280,
		Height:     720,
		Background: chart.Style{Padding: chart.Box{Top: 50}},
		XAxis: chart.XAxis{
			Style:          chart.Style{Show: true},
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d ms", int(v.(float64)))
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "max",
				Style:   chart.Style{Show: true},
				XValues: xvalues,
				YValues: maxLatencies,
			},
			chart.TimeSeries{
				Name:    "min",
				Style:   chart.Style{Show: true},
				XValues: xvalues,
				YValues: minLatencies,
			},
			chart.TimeSeries{
				Name:    "avg",
				Style:   chart.Style{Show: true},
				XValues: xvalues,
				YValues: avgLatencies,
			},
		},
	}
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
	return graph.Render(chart.PNG, w)
}

func drawThroughputs(w io.Writer, series report.TimeSeries) error {
	n := len(series)
	xvalues := make([]time.Time, n)
	throughputs := make([]float64, n)
	for i, p := range series {
		xvalues[i] = time.Unix(p.Timestamp, 0)
		throughputs[i] = float64(p.Throughput)
	}

	graph := chart.Chart{
		Width:      1280,
		Height:     720,
		Background: chart.Style{Padding: chart.Box{Top: 50}},
		XAxis: chart.XAxis{
			Style:          chart.Style{Show: true},
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "throughputs",
				Style:   chart.Style{Show: true},
				XValues: xvalues,
				YValues: throughputs,
			},
		},
	}
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
	return graph.Render(chart.PNG, w)
}
