package jobs

import "testing"

func TestTextReportJob(t *testing.T) {
	job := ReportJob{
		ResultFiles: []string{"./testdata/a.tbr"},
		Format:      "text",
	}
	if err := job.Execute(); err != nil {
		t.Errorf("job execute: %v", err)
	}
}
