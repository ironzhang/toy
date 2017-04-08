package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"plugin"

	"github.com/ironzhang/toy/framework"
)

func main() {
	p, err := plugin.Open("./robots/test-robot/robot.so")
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

	robots, err := NewRobots(100, "./robots/test-robot/robot.json")
	if err != nil {
		fmt.Printf("new robots: %v", err)
		return
	}

	var schedulers []framework.Scheduler
	data, err := ioutil.ReadFile("./robots/test-robot/schedulers.json")
	if err != nil {
		fmt.Printf("read file: %v", err)
		return
	}
	if err = json.Unmarshal(data, &schedulers); err != nil {
		fmt.Printf("json unmarshal: %v", err)
		return
	}

	(&framework.Work{
		Robots:     robots,
		Schedulers: schedulers,
	}).Run(context.Background())
}
