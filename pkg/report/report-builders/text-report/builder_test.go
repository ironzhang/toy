package text_report

import (
	"os"
	"testing"

	"github.com/ironzhang/toy/pkg/report/report-builders/tests"
)

func TestBuilder(t *testing.T) {
	r := tests.MakeTestResult("test")
	b := Builder{W: os.Stdout}
	if err := b.Build(r); err != nil {
		t.Error(err)
	}
}
