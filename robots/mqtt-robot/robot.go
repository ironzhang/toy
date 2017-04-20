package main

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ironzhang/golang/jsoncfg"
	"github.com/ironzhang/toy/framework/robot"
	"github.com/sirupsen/logrus"
)

var errTimeout = errors.New("timeout")

var (
	addr    = "tcp://localhost:1883"
	timeout = 5 * time.Second
	payload = strings.Repeat("0", 50)
)

type Options struct {
	Addr        string
	Timeout     string
	PayloadSize int
	Start       int
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	var opts Options
	err := jsoncfg.LoadFromFile(file, &opts)
	if err != nil {
		return nil, err
	}

	addr = opts.Addr
	timeout, err = time.ParseDuration(opts.Timeout)
	if err != nil {
		return nil, err
	}
	payload = strings.Repeat("0", opts.PayloadSize)

	robots := make([]robot.Robot, 0, n)
	for i := 1; i <= n; i++ {
		robots = append(robots, &Robot{ok: true, id: opts.Start + i})
	}
	return robots, nil
}

type Robot struct {
	ok bool
	c  mqtt.Client

	id int
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
	opts.AddBroker(addr)
	opts.AutoReconnect = false
	opts.DefaultPublishHander = r.OnMessage
	opts.OnConnectionLost = r.OnConnectionLost
	return opts
}

var msgcnt int64

func (r *Robot) OnMessage(c mqtt.Client, msg mqtt.Message) {
	n := atomic.AddInt64(&msgcnt, 1)
	if n%100000 == 0 {
		logrus.WithField("msgcnt", n).Info("on message")
	}
}

var lost int64

func (r *Robot) OnConnectionLost(c mqtt.Client, err error) {
	r.ok = false

	n := atomic.AddInt64(&lost, 1)
	logrus.WithError(err).WithField("lost", n).Errorf("robot(%s) lost connection", r.id)
}
