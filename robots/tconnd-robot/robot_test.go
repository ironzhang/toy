package main

import "testing"

func TestRobot(t *testing.T) {
	var r = Robot{ok: true, addr: "localhost:2000"}

	var err error
	var actions = []string{
		"Connect",
		"PingPong",
		"PingPong",
		"PingPong",
		"Disconnect",
	}
	for i, name := range actions {
		if err = r.Do(name); err != nil {
			t.Fatalf("actions[%d]: %s: %v", i, name, err)
		}
	}
}
