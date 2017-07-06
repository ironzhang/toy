package jobs

import "testing"

func TestBenchmarkJob(t *testing.T) {
	job := BenchmarkJob{
		Verbose:    1,
		Ask:        false,
		Record:     false,
		RecordFile: "benchmark.tbr",
		RobotNum:   1,
		RobotPath:  "../../robots/test-robot",
	}
	if err := job.Execute(); err != nil {
		t.Errorf("job execute: %v", err)
	}
}
