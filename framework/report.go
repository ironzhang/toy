package framework

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type report struct {
	name    string
	request int
	qps     int
	total   time.Duration
	average time.Duration
	lats    []time.Duration
	errs    map[string]int
}

func makeReport(name string, request, qps int, total time.Duration, resultc chan result) *report {
	var sum time.Duration
	errs := make(map[string]int)
	lats := make([]time.Duration, 0, len(resultc))
	for res := range resultc {
		if res.err != nil {
			errs[res.err.Error()]++
		} else {
			sum += res.duration
			lats = append(lats, res.duration)
		}
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
	return &report{
		name:    name,
		request: request,
		qps:     qps,
		total:   total,
		average: sum / time.Duration(len(lats)),
		lats:    lats,
		errs:    errs,
	}
}

func (r *report) print(w io.Writer) {
	if len(r.lats) > 0 {
		fmt.Fprintf(w, "\nSummary: %s\n", r.name)
		fmt.Fprintf(w, "  Total:\t%s\n", r.total)
		fmt.Fprintf(w, "  Slowest:\t%s\n", r.lats[len(r.lats)-1])
		fmt.Fprintf(w, "  Fastest:\t%s\n", r.lats[0])
		fmt.Fprintf(w, "  Average:\t%s\n", r.average)
		fmt.Fprintf(w, "  Requests:\t%d/%d\n", len(r.lats), r.request)
		fmt.Fprintf(w, "  Requests/sec:\t%4.4f/%d\n", float64(len(r.lats))/r.total.Seconds(), r.qps)
		r.printHistogram(w)
		r.printLatencies(w)
	}

	if len(r.errs) > 0 {
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
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]time.Duration, len(pctls))

	n := len(r.lats)
	for i, p := range pctls {
		k := (p*n - 1) / 100
		data[i] = r.lats[k]
	}
	for i := range data {
		if i+1 < len(data) && data[i] == data[i+1] {
			data[i] = 0
		}
	}

	fmt.Fprintf(w, "\nLatency distribution:\n")
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			fmt.Fprintf(w, "  %d%% in %s\n", pctls[i], data[i])
		}
	}
}

func (r *report) printErrors(w io.Writer) {
	fmt.Fprintf(w, "\nError distribution:\n")
	for err, num := range r.errs {
		fmt.Fprintf(w, "  [%d]\t%s\n", num, err)
	}
}
