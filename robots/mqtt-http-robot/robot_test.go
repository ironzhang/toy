package main

import "testing"

func TestRobot(t *testing.T) {
	r := Robot{addr: "http://localhost:8000", topic: "1"}
	if err := r.Do("Publish"); err != nil {
		t.Error(err)
	}
}
