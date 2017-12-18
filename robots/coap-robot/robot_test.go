package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestShortPing(t *testing.T) {
	n := 1000
	for i := 0; i < n; i++ {
		r := newRobot("localhost:5683")
		if err := r.Do("ShortPing"); err != nil {
			t.Errorf("do short ping: %v", err)
		}
		//log.Printf("robot(%d) done", i)
	}
}

func RunTasks(r *Robot, tasks []string) error {
	for _, task := range tasks {
		if err := r.Do(task); err != nil {
			return fmt.Errorf("do %q: %v", task, err)
		}
		//log.Printf("robot do %q success", task)
	}
	return nil
}

func TestPing(t *testing.T) {
	n := 100
	tasks := []string{"Connect", "Ping", "Sleep", "Disconnect"}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		r := newRobot("localhost:5683")
		wg.Add(1)
		go func(r *Robot) {
			defer wg.Done()
			if err := RunTasks(r, tasks); err != nil {
				t.Errorf("run tasks: %v", err)
			}
		}(r)
	}
	wg.Wait()
}

func TestEcho(t *testing.T) {
	n := 100
	tasks := []string{"Connect", "Echo", "Sleep", "Disconnect"}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		r := newRobot("localhost:5683")
		wg.Add(1)
		go func(r *Robot) {
			defer wg.Done()
			if err := RunTasks(r, tasks); err != nil {
				t.Errorf("run tasks: %v", err)
			}
		}(r)
	}
	wg.Wait()
}
