package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ironzhang/gomqtt/pkg/proto/public"
	"github.com/ironzhang/matrix/restful"
	"github.com/ironzhang/toy/framework/jsoncfg"
	"github.com/ironzhang/toy/framework/robot"
)

var (
	qos     = 0
	payload = bytes.Repeat([]byte{0}, 50)
	client  = &restful.Client{
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          0,
				MaxIdleConnsPerHost:   50,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 20 * time.Second,
		},
	}
)

type Options struct {
	Addrs       []string
	Timeout     jsoncfg.Duration
	Qos         int
	PayloadSize int
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	var opts Options
	if err := jsoncfg.LoadFromFile(file, &opts); err != nil {
		return nil, err
	}
	if len(opts.Addrs) <= 0 {
		return nil, errors.New("addrs is empty")
	}
	qos = opts.Qos
	payload = bytes.Repeat([]byte{0}, opts.PayloadSize)
	client.Client.Timeout = time.Duration(opts.Timeout)

	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &Robot{addr: opts.Addrs[i%len(opts.Addrs)], topic: fmt.Sprint(i)})
	}
	return robots, nil
}

type Robot struct {
	addr  string
	topic string
}

func (r *Robot) OK() bool {
	return true
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Publish":
		return r.Publish()
	}
	return fmt.Errorf("unknown task name: %q", name)
}

func (r *Robot) Publish() error {
	req := public.PublishRequest{Topic: r.topic, Qos: qos, Payload: payload}
	return client.Post(fmt.Sprintf("%s/v1/publish", r.addr), req, nil)
}
