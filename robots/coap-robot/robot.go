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
	coap.EnableCache = false
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
		ok:   true,
		addr: addr,
		done: make(chan struct{}, 1),
	}
}

type Robot struct {
	ok   bool
	addr string
	done chan struct{}
	conn *coap.Conn
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
	case "Echo":
		return p.Echo()
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
	case "/echoFinish":
		p.done <- struct{}{}
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
		log.Printf("new request: %v", err)
		return err
	}
	_, err = p.conn.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}
	return nil
}

func (p *Robot) Echo() error {
	urlstr := fmt.Sprintf("coap://%s/echo", p.addr)
	req, err := coap.NewRequest(true, coap.POST, urlstr, nil)
	if err != nil {
		log.Printf("new request: %v", err)
		return err
	}
	_, err = p.conn.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}

	t := time.NewTimer(5 * time.Second)
	defer t.Stop()
	select {
	case <-p.done:
		return nil
	case <-t.C:
		return errors.New("wait echo finish timeout")
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
		log.Printf("new request: %v", err)
		return err
	}
	_, err = client.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}
	return nil
}
