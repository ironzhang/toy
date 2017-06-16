package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/ironzhang/gomqtt/pkg/packet"
	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework/robot"
)

var (
	errRead  = errors.New("read error")
	errWrite = errors.New("write error")
)

var (
	Verbose          = true
	Payload          = "hello, world"
	ConnectTimeout   = 20 * time.Second
	ReadWriteTimeout = 10 * time.Second
)

type Options struct {
	Addrs            []string
	Verbose          bool
	PayloadSize      int
	ConnectTimeout   jsoncfg.Duration
	ReadWriteTimeout jsoncfg.Duration
}

func NewRobots(n int, file string) ([]robot.Robot, error) {
	var opts Options
	err := jsoncfg.LoadFromFile(file, &opts)
	if err != nil {
		return nil, err
	}
	if len(opts.Addrs) <= 0 {
		return nil, errors.New("addrs is empty")
	}

	Verbose = opts.Verbose
	Payload = strings.Repeat("0", opts.PayloadSize)
	ConnectTimeout = time.Duration(opts.ConnectTimeout)
	ReadWriteTimeout = time.Duration(opts.ReadWriteTimeout)

	l := len(opts.Addrs)
	robots := make([]robot.Robot, 0, n)
	for i := 0; i < n; i++ {
		robots = append(robots, &Robot{ok: true, addr: opts.Addrs[i%l]})
	}
	return robots, nil
}

type Robot struct {
	ok   bool
	addr string

	c net.Conn
}

func (r *Robot) OK() bool {
	return r.ok
}

func (r *Robot) Do(name string) error {
	switch name {
	case "Connect":
		return r.Connect()
	case "PingPong":
		return r.PingPong()
	case "Disconnect":
		return r.Disconnect()
	}
	return fmt.Errorf("unknown task name: %q", name)
}

func (r *Robot) Connect() (err error) {
	if r.c, err = net.DialTimeout("tcp", r.addr, ConnectTimeout); err != nil {
		r.ok = false
		return err
	}

	connect := packet.NewConnectPacket()
	connect.Keepalive = 60
	r.c.SetWriteDeadline(time.Now().Add(ReadWriteTimeout))
	if err = connect.Write(r.c); err != nil {
		r.ok = false
		return err
	}

	r.c.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
	cp, err := packets.ReadPacket(r.c)
	if err != nil {
		r.ok = false
		return err
	}
	connack, ok := cp.(*packets.ConnackPacket)
	if !ok {
		r.ok = false
		return fmt.Errorf("read packet not a connack packet")
	}
	if connack.ReturnCode != packets.Accepted {
		r.ok = false
		if e, ok := packets.ConnErrors[connack.ReturnCode]; ok {
			return e
		}
		return fmt.Errorf("connack.ReturnCode=%d", connack.ReturnCode)
	}

	return nil
}

func (r *Robot) PingPong() (err error) {
	ping := packet.NewPingreqPacket()
	r.c.SetWriteDeadline(time.Now().Add(ReadWriteTimeout))
	if err = ping.Write(r.c); err != nil {
		log.Printf("write pingreq: %v", err)
		return err
	}

	r.c.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
	cp, err := packets.ReadPacket(r.c)
	if err != nil {
		log.Printf("read pingresp: %v", err)
		return err
	}
	_, ok := cp.(*packets.PingrespPacket)
	if !ok {
		log.Printf("read packet not a pingresp packet")
		return fmt.Errorf("read packet not a pingresp packet")
	}

	return nil
}

func (r *Robot) Disconnect() error {
	return r.c.Close()
}
