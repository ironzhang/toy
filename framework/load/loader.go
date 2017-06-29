package load

import (
	"io"
	"os"

	"github.com/ironzhang/toy/framework/codec"
	"github.com/ironzhang/toy/framework/report"
)

type Loader struct {
	Verbose bool
}

func (l *Loader) loadFile(filename string, reports []report.Report) ([]report.Report, error) {
	f, err := os.Open(filename)
	if err != nil {
		return reports, err
	}
	defer f.Close()

	var r report.Report
	dec := codec.NewDecoder(f)
	for {
		if err = dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			return reports, err
		}
		if l.Verbose {
			r.Print(os.Stdout)
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func (l *Loader) Load(filenames ...string) (reports []report.Report, err error) {
	for _, file := range filenames {
		if reports, err = l.loadFile(file, reports); err != nil {
			return nil, err
		}
	}
	return reports, nil
}
