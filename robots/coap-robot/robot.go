package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ironzhang/coap"
	"github.com/ironzhang/toy/framework/robot"
)

var client coap.Client

func init() {
	coap.Verbose = 0
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	robots := make([]robot.Robot, 0, n)
	for i := 1; i <= n; i++ {
		robots = append(robots, newRobot("localhost:5683"))
	}
	return robots, nil
}

func newRobot(addr string) *Robot {
	return &Robot{
		ok:       true,
		addr:     addr,
		pingc:    make(chan struct{}, 1),
		observec: make(chan struct{}, 1),
	}
}

type Robot struct {
	ok       bool
	addr     string
	conn     *coap.Conn
	pingc    chan struct{}
	observec chan struct{}
}

func (p *Robot) OK() bool {
	return p.ok
}

func (p *Robot) Do(name string) error {
	switch name {
	case "Connect":
		return p.Connect()
	case "Disconnect":
		return p.Disconnect()
	case "Ping":
		return p.Ping()
	case "Observe":
		return p.Observe()
	case "Sleep":
		return p.Sleep()
	case "ShortPing":
		return p.ShortPing()
	default:
		return fmt.Errorf("unknown task(%s)", name)
	}
}

func (p *Robot) ServeCOAP(w coap.ResponseWriter, r *coap.Request) {
	//log.Printf("ServeCOAP: %q", r.URL.Path)
	switch r.URL.Path {
	case "/ping":
		p.pingc <- struct{}{}
	case "/observe":
		p.observec <- struct{}{}
	}
}

func (p *Robot) Connect() (err error) {
	urlstr := fmt.Sprintf("coap://%s", p.addr)
	p.conn, err = client.Dial(urlstr, p, nil)
	if err != nil {
		p.ok = false
		return err
	}
	return nil
}

func (p *Robot) Disconnect() error {
	return p.conn.Close()
}

func (p *Robot) Ping() error {
	urlstr := fmt.Sprintf("coap://%s/ping", p.addr)
	req, err := coap.NewRequest(true, coap.POST, urlstr, nil)
	if err != nil {
		return err
	}
	_, err = p.conn.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}

	t := time.NewTimer(10 * time.Second)
	defer t.Stop()
	select {
	case <-p.pingc:
		return nil
	case <-t.C:
		return errors.New("wait ping timeout")
	}

	return nil
}

func (p *Robot) Observe() error {
	urlstr := fmt.Sprintf("coap://%s/observe", p.addr)
	req, err := coap.NewRequest(true, coap.POST, urlstr, nil)
	if err != nil {
		return err
	}
	_, err = p.conn.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}

	t := time.NewTimer(10 * time.Second)
	defer t.Stop()
	select {
	case <-p.observec:
		return nil
	case <-t.C:
		return errors.New("wait observe timeout")
	}

	return nil
}

func (p *Robot) Sleep() error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (p *Robot) ShortPing() error {
	urlstr := fmt.Sprintf("coap://%s/short/ping", p.addr)
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
