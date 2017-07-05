package jobs

import "testing"

func TestTextReportJob(t *testing.T) {
	job := ReportJob{
		ResultFiles: []string{"./testdata/test.tbr"},
		Format:      "text",
	}
	if err := job.Execute(); err != nil {
		t.Errorf("job execute: %v", err)
	}
}

func TestHTMLReportJob(t *testing.T) {
	job := ReportJob{
		ResultFiles: []string{"./testdata/test.tbr"},
		Format:      "html",
		Template:    "../report/builders/html-report/templates/report.template",
		OutputDir:   "./output",
	}
	if err := job.Execute(); err != nil {
		t.Errorf("job execute: %v", err)
	}
}
