package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"plugin"

	"github.com/ironzhang/golang/jsoncfg"
	"github.com/ironzhang/toy/framework"
)

const (
	ROBOT_SO        = "robot.so"
	ROBOT_JSON      = "robot.json"
	SCHEDULERS_JSON = "schedulers.json"
)

func main() {
	var robotNum int
	var robotPath string

	flag.IntVar(&robotNum, "robot-num", 1, "run robot number")
	flag.StringVar(&robotPath, "robot-path", "./robots/test-robot", "robot plugin path")
	flag.Parse()

	p, err := plugin.Open(fmt.Sprintf("%s/%s", robotPath, ROBOT_SO))
	if err != nil {
		fmt.Printf("plugin open: %v", err)
		return
	}

	s, err := p.Lookup("NewRobots")
	if err != nil {
		fmt.Printf("plugin lookup: %v", err)
		return
	}

	NewRobots, ok := s.(func(int, string) ([]framework.Robot, error))
	if !ok {
		fmt.Printf("%T is unexpect", s)
		return
	}

	robots, err := NewRobots(robotNum, fmt.Sprintf("%s/%s", robotPath, ROBOT_JSON))
	if err != nil {
		fmt.Printf("new robots: %v", err)
		return
	}

	var schedulers []framework.Scheduler
	err = jsoncfg.LoadFromFile(fmt.Sprintf("%s/%s", robotPath, SCHEDULERS_JSON), &schedulers)
	if err != nil {
		fmt.Printf("load schedulers json from file: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		defer cancel()
	}()
	(&framework.Work{Robots: robots, Schedulers: schedulers}).Run(ctx)
}
