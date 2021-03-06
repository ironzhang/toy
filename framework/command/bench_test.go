package command

import "testing"

func TestBenchCmdParse(t *testing.T) {
	args := []string{
		"-verbose", "2",
		"-ask",
		"-robot-num", "100",
		"-robot-path", "test-robot",
	}
	var got BenchCmd
	var want = BenchCmd{
		verbose:   2,
		ask:       true,
		robotNum:  100,
		robotPath: "test-robot",
	}
	if err := got.parse(args); err != nil {
		t.Fatalf("cmd parse: %v", err)
	}
	if got != want {
		t.Errorf("%v != %v", got, want)
	}
}

func TestBenchCmdExecute(t *testing.T) {
	cmd := BenchCmd{
		verbose: 0,
		ask:     false,
		//record:     true,
		output:    "bench.tbr",
		robotNum:  1,
		robotPath: "../../robots/test-robot",
	}
	if err := cmd.execute(); err != nil {
		t.Errorf("cmd execute: %v", err)
	}
}
