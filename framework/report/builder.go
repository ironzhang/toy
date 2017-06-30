package report

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func mergeReports(reports []Report) []Report {
	var i int
	var merges []Report
	for _, r := range reports {
		for i = 0; i < len(merges); i++ {
			if merges[i].Name == r.Name {
				break
			}
		}
		if i < len(merges) {
			merges[i].merge(r)
		} else {
			merges = append(merges, r)
		}
	}
	return merges
}

type Builder struct {
	Template   string
	OutputDir  string
	SampleSize int
}

func (b *Builder) processRecords(orecords []Record) []Record {
	// 过滤错误记录
	nrecords := make([]Record, 0, len(orecords))
	for _, r := range orecords {
		if r.Err == "" {
			nrecords = append(nrecords, r)
		}
	}

	// 排序
	sort.Slice(nrecords, func(i, j int) bool { return nrecords[i].Start.Before(nrecords[j].Start) })

	// 采样
	size := b.SampleSize
	if size <= 0 {
		size = 500
	}
	return sampling(nrecords, size)
}

func (b *Builder) buildImage(r Report) (string, error) {
	records := b.processRecords(r.Records)
	if len(records) <= 1 {
		return "", nil
	}

	filename := fmt.Sprintf("%s/%s.png", b.OutputDir, r.Name)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return filepath.Base(filename), renderLatencies(f, records)
}

func (b *Builder) buildHistogramImage(r *report) (string, error) {
	if len(r.lats) <= 0 {
		return "", nil
	}

	filename := fmt.Sprintf("%s/%s_histogram.png", b.OutputDir, r.Name)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buckets := histogramBucket(r.lats)
	return filepath.Base(filename), renderHistogram(f, buckets)
}

func (b *Builder) buildHTML(data interface{}) error {
	filename := fmt.Sprintf("%s/report.html", b.OutputDir)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderTemplate(f, b.Template, data)
}

func (b *Builder) Build(reports []Report) error {
	os.MkdirAll(b.OutputDir, os.ModePerm)
	var data []map[string]interface{}
	for _, r := range mergeReports(reports) {
		m := make(map[string]interface{})
		report := makeReport(r.Name, r.Request, r.Concurrent, r.QPS, r.Total, r.Records)
		latenciesImg, err := b.buildImage(r)
		if err != nil {
			return err
		}
		histogramImg, err := b.buildHistogramImage(report)
		if err != nil {
			return err
		}
		m["report"] = report
		m["latenciesImg"] = latenciesImg
		m["histogramImg"] = histogramImg
		data = append(data, m)
	}
	return b.buildHTML(data)
}
