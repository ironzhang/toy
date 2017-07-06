package command

import (
	"flag"
	"fmt"
	"os"
	"plugin"
	"time"

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

type BenchCmd struct {
	verbose    int
	ask        bool
	record     bool
	recordFile string
	robotNum   int
	robotPath  string
}

func (c *BenchCmd) Run(args []string) error {
	if err := c.parse(args); err != nil {
		return nil
	}
	return c.execute()
}

func (c *BenchCmd) parse(args []string) error {
	var fs flag.FlagSet
	fs.Usage = func() {
		fmt.Print("Usage: toy bench [OPTIONS]\n\n")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.verbose, "verbose", 1, "verbose level")
	fs.BoolVar(&c.ask, "ask", false, "ask execute scheduler")
	fs.BoolVar(&c.record, "record", false, "record the result")
	fs.StringVar(&c.recordFile, "record-file", "", "record file")
	fs.IntVar(&c.robotNum, "robot-num", 1, "robot num")
	fs.StringVar(&c.robotPath, "robot-path", "./robots/test-robot", "robot path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return nil
}

func (c *BenchCmd) execute() error {
	schedulers, err := c.loadSchedulers()
	if err != nil {
		return err
	}
	robots, err := c.newRobots()
	if err != nil {
		return err
	}

	var encoder report.Encoder
	if c.record {
		if c.recordFile == "" {
			c.recordFile = fmt.Sprintf("record.%s.tbr", time.Now().Format(time.RFC3339))
		}
		f, err := os.Create(c.recordFile)
		if err != nil {
			return err
		}
		defer f.Close()
		encoder = report.NewEncoder(f)
	}

	(&benchmark.Benchmark{
		Verbose:    c.verbose,
		Ask:        c.ask,
		Encoder:    encoder,
		Robots:     robots,
		Schedulers: schedulers,
	}).Run()

	return nil
}

func (c *BenchCmd) loadSchedulers() ([]benchmark.Scheduler, error) {
	var schedulers []benchmark.Scheduler
	if err := jsoncfg.LoadFromFile(fmt.Sprintf("%s/%s", c.robotPath, schedulers_json), &schedulers); err != nil {
		return nil, fmt.Errorf("load schedulers json from file: %v", err)
	}
	return schedulers, nil
}

func (c *BenchCmd) newRobots() ([]robot.Robot, error) {
	p, err := plugin.Open(fmt.Sprintf("%s/%s", c.robotPath, robot_so))
	if err != nil {
		return nil, fmt.Errorf("plugin open: %v", err)
	}
	s, err := p.Lookup("NewRobots")
	if err != nil {
		return nil, fmt.Errorf("plugin lookup: %v", err)
	}
	NewRobots, ok := s.(func(int, string) ([]robot.Robot, error))
	if !ok {
		return nil, fmt.Errorf("NewRobots is unexpect: %T", s)
	}
	robots, err := NewRobots(c.robotNum, fmt.Sprintf("%s/%s", c.robotPath, robot_json))
	if err != nil {
		return nil, fmt.Errorf("new robots: %v", err)
	}
	return robots, nil
}
