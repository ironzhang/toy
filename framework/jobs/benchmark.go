package jobs

import (
	"fmt"
	"os"
	"plugin"

	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework/benchmark"
	"github.com/ironzhang/toy/framework/report"
	"github.com/ironzhang/toy/framework/robot"
)

const (
	robot_so        = "robot.so"
	robot_json      = "robot.json"
	schedulers_json = "schedulers.json"
)

type BenchmarkJob struct {
	Verbose    int
	Ask        bool
	Record     bool
	RecordFile string
	RobotNum   int
	RobotPath  string
}

func (p *BenchmarkJob) Execute() error {
	schedulers, err := p.loadSchedulers()
	if err != nil {
		return err
	}
	robots, err := p.newRobots()
	if err != nil {
		return err
	}

	var encoder report.Encoder
	if p.Record {
		f, err := os.Create(p.RecordFile)
		if err != nil {
			return err
		}
		defer f.Close()
		encoder = report.NewEncoder(f)
	}

	(&benchmark.Benchmark{
		Verbose:    p.Verbose,
		Ask:        p.Ask,
		Encoder:    encoder,
		Robots:     robots,
		Schedulers: schedulers,
	}).Run()

	return nil
}

func (p *BenchmarkJob) loadSchedulers() ([]benchmark.Scheduler, error) {
	var schedulers []benchmark.Scheduler
	if err := jsoncfg.LoadFromFile(fmt.Sprintf("%s/%s", p.RobotPath, schedulers_json), &schedulers); err != nil {
		return nil, fmt.Errorf("load schedulers json from file: %v", err)
	}
	return schedulers, nil
}

func (p *BenchmarkJob) newRobots() ([]robot.Robot, error) {
	pg, err := plugin.Open(fmt.Sprintf("%s/%s", p.RobotPath, robot_so))
	if err != nil {
		return nil, fmt.Errorf("plugin open: %v", err)
	}
	s, err := pg.Lookup("NewRobots")
	if err != nil {
		return nil, fmt.Errorf("plugin lookup: %v", err)
	}
	NewRobots, ok := s.(func(int, string) ([]robot.Robot, error))
	if !ok {
		return nil, fmt.Errorf("NewRobots is unexpect: %T", s)
	}
	robots, err := NewRobots(p.RobotNum, fmt.Sprintf("%s/%s", p.RobotPath, robot_json))
	if err != nil {
		return nil, fmt.Errorf("new robots: %v", err)
	}
	return robots, nil
}
