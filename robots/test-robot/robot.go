package main

import (
	"time"

	"github.com/ironzhang/toy/framework/robot"
)

func NewRobots(n int, file string) ([]robot.Robot, error) {
	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &Robot{})
	}
	return robots, nil
}

type Robot struct {
}

func (r *Robot) OK() bool {
	return true
}

func (r *Robot) Do(name string) error {
	//fmt.Println(name)
	switch name {
	case "Connect":
		time.Sleep(10 * time.Millisecond)
	case "Prepare":
		time.Sleep(20 * time.Millisecond)
	case "Publish":
		time.Sleep(100 * time.Microsecond)
	case "Disconnect":
		time.Sleep(10 * time.Microsecond)
	}
	return nil
}
