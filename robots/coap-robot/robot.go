package main

import (
	"fmt"
	"log"

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
		robots = append(robots, &Robot{ok: true, addr: "localhost:5683"})
	}
	return robots, nil
}

type Robot struct {
	ok   bool
	addr string
	conn *coap.Conn
}

func (r *Robot) OK() bool {
	return r.ok
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Connect":
		return r.Connect()
	case "Disconnect":
		return r.Disconnect()
	case "Ping":
		return r.Ping()
	case "ShortPing":
		return r.ShortPing()
	default:
		return fmt.Errorf("unknown task(%s)", name)
	}
}

func (r *Robot) ShortPing() error {
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

func (r *Robot) Connect() (err error) {
	urlstr := fmt.Sprintf("coap://%s", r.addr)
	r.conn, err = client.Dial(urlstr, nil, nil)
	if err != nil {
		r.ok = false
		return err
	}
	return nil
}

func (r *Robot) Disconnect() error {
	return r.conn.Close()
}

func (r *Robot) Ping() error {
	urlstr := fmt.Sprintf("coap://%s/ping", r.addr)
	req, err := coap.NewRequest(true, coap.POST, urlstr, nil)
	if err != nil {
		return err
	}
	_, err = r.conn.SendRequest(req)
	if err != nil {
		log.Printf("send request: %v", err)
		return err
	}
	return nil
}
