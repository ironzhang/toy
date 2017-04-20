package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"ac-common-go/common"
	"ac-common-go/util/jsoncfg"
	"ac-gateway/server"
	"ac-mqtt/topicname"
	"ac-mqtt/username"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ironzhang/toy/framework/robot"
	log "github.com/sirupsen/logrus"
)

var errTimeout = errors.New("timeout")

var ack int64
var control int64
var payload = bytes.Repeat([]byte{0}, 50)

type Options struct {
	Addr          string
	Timeout       string
	PayloadSize   int
	MajorDomainID string
	SubDomainID   string
	Start         int
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	var opts Options
	if err := jsoncfg.LoadFromFile(file, &opts); err != nil {
		return nil, err
	}
	timeout, err := time.ParseDuration(opts.Timeout)
	if err != nil {
		return nil, err
	}
	payload = bytes.Repeat([]byte{0}, opts.PayloadSize)

	robots := make([]robot.Robot, 0, n)
	for i := 1; i <= n; i++ {
		r := &Robot{
			ok:      true,
			addr:    opts.Addr,
			timeout: timeout,
			user:    User(opts.MajorDomainID, opts.SubDomainID, fmt.Sprint(opts.Start+i)),
		}
		robots = append(robots, r)
	}
	return robots, nil
}

func User(majorDomainID, subDomainID, deviceID string) username.User {
	return username.User{
		MajorDomainID: majorDomainID,
		SubDomainID:   subDomainID,
		DeviceID:      deviceID,
		AccessVersion: 4,
		ModVersion:    "v1",
		DevVersion:    "v1",
		DevType:       "L",
		ModType:       "L",
		IsHTTPOTA:     true,
	}
}

type Robot struct {
	ok bool
	c  mqtt.Client

	addr    string
	timeout time.Duration
	user    username.User
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
	case "Nothing":
		return nil
	}
	return fmt.Errorf("unknown task name: %q", name)
}

func (r *Robot) Connect() error {
	r.c = mqtt.NewClient(r.mqttClientOptions())
	t := r.c.Connect()
	if !t.WaitTimeout(r.timeout) {
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
	topic := r.user.DeviceDownstreamTopic()
	t := r.c.Subscribe(topic, 0, nil)
	if !t.WaitTimeout(r.timeout) {
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
	var msg server.Message
	msg.Header.Version = 4
	msg.Header.MsgCode = server.ZC_REPORT_BASE
	msg.Payload = payload
	return r.sendServerMessage(&msg)
}

func (r *Robot) Disconnect() error {
	r.c.Disconnect(0)
	return nil
}

func (r *Robot) OnMessage(c mqtt.Client, m mqtt.Message) {
	var msg server.Message
	if err := msg.Deserialize(m.Payload()); err != nil {
		log.Errorf("message deserialize: %v", err)
		return
	}

	switch {
	case msg.Header.MsgCode >= server.ZC_CODE_BASE && msg.Header.MsgCode < server.ZC_REPORT_BASE:
		if err := r.sendServerMessage(&msg); err != nil {
			log.Errorf("send server message: %v", err)
			return
		}

		n := atomic.AddInt64(&control, 1)
		if n%100000 == 0 {
			log.WithField("control", n).Infoln()
		}
	default:
		n := atomic.AddInt64(&ack, 1)
		if n%100000 == 0 {
			log.WithField("ack", n).Infoln()
		}
	}
}

func (r *Robot) OnConnectionLost(c mqtt.Client, err error) {
	r.ok = false
	log.WithField("robot", r.user.DeviceID).WithError(err).Info("connection lost")
}

func (r *Robot) mqttClientOptions() *mqtt.ClientOptions {
	username, _ := json.Marshal(r.user)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(r.addr)
	opts.Username = string(username)
	opts.KeepAlive = time.Minute
	opts.AutoReconnect = false
	opts.DefaultPublishHander = r.OnMessage
	opts.OnConnectionLost = r.OnConnectionLost
	return opts
}

func (r *Robot) sendServerMessage(msg *server.Message) error {
	msg.Header.Checksum = common.CCITTChecksum(msg.Payload)
	payload, err := msg.Serialize()
	if err != nil {
		return err
	}

	topic := topicname.MakeDeviceUpstreamTopic(r.user.MajorDomainID, r.user.SubDomainID, r.user.DeviceID)
	t := r.c.Publish(topic, 0, false, payload)
	if !t.WaitTimeout(r.timeout) {
		return errTimeout
	}
	if err := t.Error(); err != nil {
		return err
	}
	return nil
}
