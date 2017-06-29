package main

import (
	"flag"
	"fmt"
	"os"
	"plugin"

	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework"
	"github.com/ironzhang/toy/framework/codec"
	"github.com/ironzhang/toy/framework/load"
	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/robot"
)

const (
	ROBOT_SO        = "robot.so"
	ROBOT_JSON      = "robot.json"
	SCHEDULERS_JSON = "schedulers.json"

	OUTPUT_DIR      = "./output"
	REPORT_TEMPLATE = "./framework/report/templates/report.template"
)

var (
	Verbose bool
)

func MakeReport() error {
	l := load.Loader{
		Verbose: Verbose,
	}
	reports, err := l.Load(flag.Args()...)
	if err != nil {
		return fmt.Errorf("load reports: %v", err)
	}

	b := report.Builder{
		Template:   REPORT_TEMPLATE,
		OutputDir:  OUTPUT_DIR,
		SampleSize: 500,
	}
	if err = b.Build(reports); err != nil {
		return fmt.Errorf("build report: %v", err)
	}
	return nil
}

func main() {
	var ask bool
	var report bool
	var robotNum int
	var robotPath string
	var recordFile string
	flag.BoolVar(&Verbose, "verbose", false, "print verbose info")
	flag.BoolVar(&ask, "ask", false, "ask execute task")
	flag.BoolVar(&report, "report", false, "make report")
	flag.StringVar(&recordFile, "record", "", "record file")
	flag.IntVar(&robotNum, "robot-num", 1, "run robot number")
	flag.StringVar(&robotPath, "robot-path", "./robots/test-robot", "robot plugin path")
	flag.Parse()

	if report {
		if err := MakeReport(); err != nil {
			fmt.Printf("make report: %v\n", err)
			return
		}
		return
	}

	p, err := plugin.Open(fmt.Sprintf("%s/%s", robotPath, ROBOT_SO))
	if err != nil {
		fmt.Printf("plugin open: %v\n", err)
		return
	}

	s, err := p.Lookup("NewRobots")
	if err != nil {
		fmt.Printf("plugin lookup: %v\n", err)
		return
	}

	NewRobots, ok := s.(func(int, string) ([]robot.Robot, error))
	if !ok {
		fmt.Printf("%T is unexpect\n", s)
		return
	}

	robots, err := NewRobots(robotNum, fmt.Sprintf("%s/%s", robotPath, ROBOT_JSON))
	if err != nil {
		fmt.Printf("new robots: %v\n", err)
		return
	}

	var schedulers []framework.Scheduler
	err = jsoncfg.LoadFromFile(fmt.Sprintf("%s/%s", robotPath, SCHEDULERS_JSON), &schedulers)
	if err != nil {
		fmt.Printf("load schedulers json from file: %v\n", err)
		return
	}

	var enc codec.Encoder
	if recordFile != "" {
		f, err := os.Create(recordFile)
		if err != nil {
			fmt.Printf("create record file: %v\n", err)
			return
		}
		defer f.Close()
		enc = codec.NewEncoder(f)
	}

	(&framework.Work{Ask: ask, Encoder: enc, Robots: robots, Schedulers: schedulers}).Run()
}
