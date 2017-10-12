package main

import (
	"fmt"
	"log"

	"github.com/ironzhang/coap"
	"github.com/ironzhang/toy/framework/robot"
)

var client coap.Client

func NewRobots(n int, file string) ([]robot.Robot, error) {
	robots := make([]robot.Robot, 0, n)
	for i := 1; i <= n; i++ {
		robots = append(robots, &Robot{addr: "localhost:5683"})
	}
	return robots, nil
}

type Robot struct {
	addr string
}

func (r *Robot) OK() bool {
	return true
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Ping":
		return r.Ping()
	default:
		return fmt.Errorf("unknown task(%s)", name)
	}
}

func (r *Robot) Ping() error {
	urlstr := fmt.Sprintf("coap://%s/ping", r.addr)
	req, err := coap.NewRequest(true, coap.POST, urlstr, nil)
	if err != nil {
		return err
	}
	_, err = client.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}
	return nil
}
