package main

import (
	"testing"
	"time"

	"ac-common-go/util/jsoncfg"
)

func TestOptions(t *testing.T) {
	opts := Options{
		Addr:          "tcp://139.219.236.56:1883",
		Timeout:       "5s",
		PayloadSize:   50,
		MajorDomainID: "3",
		SubDomainID:   "6",
	}

	jsoncfg.WriteToFile("example.robot.json", opts)
}

func TestRobot(t *testing.T) {
	r := Robot{
		ok:      true,
		addr:    "tcp://139.219.236.56:1883",
		timeout: 5 * time.Second,
		user:    User("3", "6", "1"),
	}

	if err := r.Connect(); err != nil {
		t.Fatalf("connect: %v", err)
	}

	if err := r.Subscribe(); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := r.Publish(); err != nil {
		t.Fatalf("publish: %v", err)
	}

	time.Sleep(time.Second)

	if err := r.Disconnect(); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
}
