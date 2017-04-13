package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ironzhang/toy/framework/robot"
)

var errTimeout = errors.New("timeout")

var (
	timeout = 5 * time.Second
	payload = strings.Repeat("A", 100)
)

func NewRobots(n int, file string) ([]robot.Robot, error) {
	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &Robot{addr: "tcp://localhost:1883", id: i + 1, ok: true})
	}
	return robots, nil
}

type Robot struct {
	addr string
	id   int

	ok bool
	c  mqtt.Client
}

func (r *Robot) OK() bool {
	return r.ok
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Connect":
		return r.Connect()
	case "Subscribe":
		return r.Subscribe()
	case "Publish":
		return r.Publish()
	case "Disconnect":
		return r.Disconnect()
	}
	return fmt.Errorf("unknown task name: %q", name)
}

func (r *Robot) Connect() error {
	r.c = mqtt.NewClient(r.MqttClientOptions())
	t := r.c.Connect()
	if !t.WaitTimeout(timeout) {
		r.ok = false
		return errTimeout
	}
	if err := t.Error(); err != nil {
		r.ok = false
		return err
	}
	return nil
}

func (r *Robot) Subscribe() error {
	t := r.c.Subscribe(fmt.Sprint(r.id), 0, nil)
	if !t.WaitTimeout(timeout) {
		r.ok = false
		return errTimeout
	}
	if err := t.Error(); err != nil {
		r.ok = false
		return err
	}
	return nil
}

func (r *Robot) Publish() error {
	t := r.c.Publish(fmt.Sprint(r.id), 0, false, payload)
	if !t.WaitTimeout(timeout) {
		return errTimeout
	}
	if err := t.Error(); err != nil {
		return err
	}
	return nil
}

func (r *Robot) Disconnect() error {
	r.c.Disconnect(0)
	return nil
}

func (r *Robot) MqttClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(r.addr)
	opts.AutoReconnect = false
	opts.DefaultPublishHander = r.OnMessage
	opts.OnConnectionLost = r.OnConnectionLost
	return opts
}

func (r *Robot) OnMessage(c mqtt.Client, msg mqtt.Message) {
}

func (r *Robot) OnConnectionLost(c mqtt.Client, err error) {
	r.ok = false
}
