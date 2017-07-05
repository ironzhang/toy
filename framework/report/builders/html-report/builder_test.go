package html_report

import (
	"testing"

	"github.com/ironzhang/toy/framework/report/builders/tests"
)

func TestBuilder(t *testing.T) {
	r1 := tests.MakeTestResult("Connect")
	r2 := tests.MakeTestResult("Disconnect")
	b := Builder{
		Template:   "./templates/report.template",
		OutputDir:  "./output",
		SampleSize: 500,
	}
	if err := b.Build(r1, r2); err != nil {
		t.Error(err)
	}
}
