package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework/robot"
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
	return nil
}

func (r *Robot) PingPong() (err error) {
	r.c.SetWriteDeadline(time.Now().Add(ReadWriteTimeout))
	if _, err = fmt.Fprintf(r.c, "%s\n", Payload); err != nil {
		r.ok = false
		return err
	}

	r.c.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
	rd := bufio.NewReader(r.c)
	line, _, err := rd.ReadLine()
	if err != nil {
		r.ok = false
		return err
	}
	if Verbose {
		fmt.Println(time.Now(), string(line))
	}

	return nil
}

func (r *Robot) Disconnect() error {
	return r.c.Close()
}
