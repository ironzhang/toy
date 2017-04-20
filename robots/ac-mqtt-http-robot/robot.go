package main

import (
	"ac"
	"ac-common-go/util/jsoncfg"
	"ac-gateway/server"
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/ironzhang/toy/framework/robot"
)

var httpClient = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100000,
		MaxIdleConnsPerHost: 100000,
	},
	Timeout: 20 * time.Second,
}

var payload = bytes.Repeat([]byte{0}, 50)

type Options struct {
	Addr          string
	PayloadSize   int
	MajorDomainID string
	SubDomainID   string
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	var opts Options
	if err := jsoncfg.LoadFromFile(file, &opts); err != nil {
		return nil, err
	}
	payload = bytes.Repeat([]byte{0}, opts.PayloadSize)

	robots := make([]robot.Robot, 0, n)
	for i := 1; i <= n; i++ {
		r := &Robot{
			c:     ac.NewZServiceClientCustom(opts.Addr, "zc-mqtt-broker", "v1", &httpClient),
			major: opts.MajorDomainID,
			sub:   opts.SubDomainID,
			id:    fmt.Sprint(i),
		}
		robots = append(robots, r)
	}
	return robots, nil
}

type Robot struct {
	c *ac.ZServiceClient

	major string
	sub   string
	id    string
}

func (r *Robot) OK() bool {
	return true
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Control":
		return r.Control()
	}
	return fmt.Errorf("unknown task name: %q", name)
}

func (r *Robot) Control() error {
	ctx := ac.NewZContext(r.major, r.sub, 0, 0, "", "zc-test")
	msg := ac.NewMsg("control")
	msg.SetContext(ctx)
	msg.PutString("physicalDeviceId", r.id)
	msg.PutString("messageCode", fmt.Sprint(server.ZC_CODE_BASE))
	msg.SetStreamPayload(uint32(len(payload)), bytes.NewReader(payload))
	resp, err := r.c.Send(msg)
	if err != nil {
		return err
	}
	if resp.IsErr() {
		return resp.GetACError()
	}
	return nil
}
