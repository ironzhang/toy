package main

import (
	"ac"
	"testing"
)

func TestRobot(t *testing.T) {
	r := Robot{
		c:     ac.NewZServiceClient("localhost:5029", "zc-mqtt-broker", "v1"),
		major: "3",
		sub:   "4",
		id:    "1",
	}

	if err := r.Do("Control"); err != nil {
		t.Errorf("do control: %v", err)
	}
}
