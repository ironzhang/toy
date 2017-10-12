package main

import "testing"

func TestRobot(t *testing.T) {
	r := Robot{addr: "localhost:5683"}
	if err := r.Do("Ping"); err != nil {
		t.Errorf("do ping: %v", err)
	}
}
