package report

//type latency struct {
//	Percent  int
//	Duration time.Duration
//}
//
//type html struct {
//	Name        string
//	Total       time.Duration
//	Slowest     time.Duration
//	Fastest     time.Duration
//	Average     time.Duration
//	Concurrent  int
//	Request     int
//	RealRequest int
//	QPS         int
//	RealQPS     float64
//	LatencyImg  string
//	Latencys    []latency
//	Errs        map[string]int
//}

/*
func OutputHTML(outdir string, reports []Report) (err error) {
	htmls := make([]html, 0, len(reports))
	for _, r := range reports {
		if err = renderRecords(fmt.Sprintf("%s/%s.png", outdir, r.Name), r.Records); err != nil {
			return err
		}
		htmls = append(htmls, html{
			Name:  r.Name,
			Total: r.Total,
		})
	}
	return nil
}
*/
