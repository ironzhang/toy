package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type Record struct {
	Err    string
	Start  time.Time
	Elapse time.Duration
}

type Report struct {
	Name       string
	Total      time.Duration
	Concurrent int
	Request    int
	QPS        int
	Records    []Record
}

func (r *Report) Print(w io.Writer) {
	makeReport(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records).print(w)
}

type latency struct {
	Percent  int
	Duration time.Duration
}

type report struct {
	Name        string
	Total       time.Duration
	Slowest     time.Duration
	Fastest     time.Duration
	Average     time.Duration
	Concurrent  int
	Request     int
	RealRequest int
	QPS         int
	RealQPS     float64
	Latencies   []latency
	Errs        map[string]int

	lats []time.Duration
}

func makeReport(name string, request, concurrent, qps int, total time.Duration, records []Record) *report {
	var sum time.Duration
	errs := make(map[string]int)
	lats := make([]time.Duration, 0, len(records))
	for _, r := range records {
		if r.Err != "" {
			errs[r.Err]++
		} else {
			sum += r.Elapse
			lats = append(lats, r.Elapse)
		}
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	var average, slowest, fastest time.Duration
	if len(lats) > 0 {
		average = sum / time.Duration(len(lats))
		slowest = lats[len(lats)-1]
		fastest = lats[0]
	}

	return &report{
		Name:        name,
		Total:       total,
		Slowest:     slowest,
		Fastest:     fastest,
		Average:     average,
		Concurrent:  concurrent,
		Request:     request,
		RealRequest: len(lats),
		QPS:         qps,
		RealQPS:     float64(len(lats)) / total.Seconds(),
		Latencies:   latencyDistribution(lats),
		Errs:        errs,
		lats:        lats,
	}
}

func (r *report) print(w io.Writer) {
	fmt.Fprintf(w, "\nSummary: %s\n", r.Name)
	if len(r.lats) > 0 {
		fmt.Fprintf(w, "  Total:\t%s\n", r.Total)
		fmt.Fprintf(w, "  Slowest:\t%s\n", r.Slowest)
		fmt.Fprintf(w, "  Fastest:\t%s\n", r.Fastest)
		fmt.Fprintf(w, "  Average:\t%s\n", r.Average)
		fmt.Fprintf(w, "  Concurrent:\t%d\n", r.Concurrent)
		fmt.Fprintf(w, "  Requests:\t%d/%d\n", r.RealRequest, r.Request)
		fmt.Fprintf(w, "  Requests/sec:\t%4.4f/%d\n", r.RealQPS, r.QPS)
		r.printHistogram(w)
		r.printLatencies(w)
	}

	if len(r.Errs) > 0 {
		r.printErrors(w)
	}

	fmt.Fprintln(w)
}

func (r *report) printHistogram(w io.Writer) {
	type bucket struct {
		d time.Duration
		c int
	}

	bc := 10
	buckets := make([]bucket, bc+1)
	fastest := r.lats[0]
	slowest := r.lats[len(r.lats)-1]
	bs := (slowest - fastest) / time.Duration(bc)
	for i := 0; i < bc; i++ {
		buckets[i].d = fastest + bs*time.Duration(i)
	}
	buckets[bc].d = slowest

	bi := 0
	max := 0
	for i := 0; i < len(r.lats); {
		if r.lats[i] <= buckets[bi].d {
			buckets[bi].c++
			if max < buckets[bi].c {
				max = buckets[bi].c
			}
			i++
		} else if bi < len(buckets)-1 {
			bi++
		}
	}

	fmt.Fprintf(w, "\nResponse time histogram:\n")
	for i := 0; i < len(buckets); i++ {
		var barLen = 0
		const barChar = "∎"
		if max > 0 {
			barLen = (buckets[i].c*40 + max/2) / max
		}
		fmt.Fprintf(w, "  %s [%d]\t|%s\n", buckets[i].d, buckets[i].c, strings.Repeat(barChar, barLen))
	}
}

func (r *report) printLatencies(w io.Writer) {
	fmt.Fprintf(w, "\nLatency distribution:\n")
	for _, lat := range r.Latencies {
		fmt.Fprintf(w, "  %d%% in %s\n", lat.Percent, lat.Duration)
	}
}

func (r *report) printErrors(w io.Writer) {
	fmt.Fprintf(w, "\nError distribution:\n")
	for err, num := range r.Errs {
		fmt.Fprintf(w, "  [%d]\t%s\n", num, err)
	}
}

func latencyDistribution(lats []time.Duration) []latency {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]time.Duration, len(pctls))

	n := len(lats)
	for i, p := range pctls {
		k := (p*n - 1) / 100
		data[i] = lats[k]
	}
	for i := range data {
		if i+1 < len(data) && data[i] == data[i+1] {
			data[i] = 0
		}
	}

	latencies := make([]latency, 0, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			latencies = append(latencies, latency{Percent: pctls[i], Duration: data[i]})
		}
	}
	return latencies
}
