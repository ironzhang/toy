package main

import (
	"log"
	"testing"
)

func TestRobot(t *testing.T) {
	n := 1000
	for i := 0; i < n; i++ {
		r := Robot{addr: "localhost:5683"}
		if err := r.Do("ShortPing"); err != nil {
			t.Errorf("do short ping: %v", err)
		}
		log.Printf("robot(%d) done", i)
	}
}
