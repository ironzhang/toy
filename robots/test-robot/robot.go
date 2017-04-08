package main

import "github.com/ironzhang/toy/framework"

func NewRobots(n int, file string) ([]framework.Robot, error) {
	robots := make([]framework.Robot, 0, n)
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
	return nil
}
