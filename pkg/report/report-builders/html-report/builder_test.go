package html_report

import (
	"testing"

	"github.com/ironzhang/toy/pkg/report/report-builders/tests"
)

func TestBuilder(t *testing.T) {
	r1 := tests.MakeTestResult("test1")
	r2 := tests.MakeTestResult("test2")
	b := Builder{
		Template:  "./templates/report.template",
		OutputDir: "./output",
	}
	if err := b.Build(r1, r2); err != nil {
		t.Error(err)
	}
}
