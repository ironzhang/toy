package main

import (
	"testing"

	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework/benchmark"
)

func TestWriteSchedulersJSON(t *testing.T) {
	schedulers := []benchmark.Scheduler{
		{
			Name: "Connect",
			N:    1,
			C:    10,
			QPS:  1000,
		},
		{
			Name: "Prepare",
			N:    1,
			C:    10,
			QPS:  1000,
		},
		{
			Name: "Publish",
			N:    100,
			C:    10,
			QPS:  5000,
		},
		{
			Name: "Disconnect",
			N:    1,
			C:    10,
			QPS:  1000,
		},
	}
	jsoncfg.WriteToFile("./schedulers.json", schedulers)
}
