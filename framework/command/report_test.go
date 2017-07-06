package command

import (
	"reflect"
	"testing"
)

func TestReportCmdParse(t *testing.T) {
	args := []string{
		"-format", "text",
		"-output-dir", "output_dir",
		"-sample-size", "1000",
		"a.tbr", "b.tbr", "c.tbr",
	}
	var got ReportCmd
	var want = ReportCmd{
		format:      "text",
		outputDir:   "output_dir",
		sampleSize:  1000,
		resultFiles: []string{"a.tbr", "b.tbr", "c.tbr"},
	}
	if err := got.parse(args); err != nil {
		t.Fatalf("cmd parse: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestTextReportCmdExecute(t *testing.T) {
	cmd := ReportCmd{
		format:      "text",
		resultFiles: []string{"./testdata/test.tbr"},
	}
	if err := cmd.execute(); err != nil {
		t.Errorf("cmd execute: %v", err)
	}
}

func TestHTMLReportCmdExecute(t *testing.T) {
	cmd := ReportCmd{
		Template:    "../report/builders/html-report/templates/report.template",
		format:      "html",
		outputDir:   "output",
		sampleSize:  500,
		resultFiles: []string{"./testdata/test.tbr"},
	}
	if err := cmd.execute(); err != nil {
		t.Errorf("cmd execute: %v", err)
	}
}
