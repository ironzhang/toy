package main

import (
	"testing"

	"github.com/ironzhang/golang/jsoncfg"
	"github.com/ironzhang/toy/framework"
)

func TestWriteSchedulersJSON(t *testing.T) {
	n := 1000
	schedulers := []framework.Scheduler{
		{
			N:           n,
			C:           10,
			QPS:         1000,
			Name:        "Connect",
			Display:     false,
			PrintReport: false,
		},
		{
			N:           n,
			C:           10,
			QPS:         1000,
			Name:        "Prepare",
			Display:     false,
			PrintReport: false,
		},
		{
			N:           n * 100,
			C:           100,
			QPS:         5000,
			Name:        "Publish",
			Display:     true,
			PrintReport: true,
		},
		{
			N:           n,
			C:           10,
			QPS:         1000,
			Name:        "Disconnect",
			Display:     false,
			PrintReport: false,
		},
	}
	jsoncfg.WriteToFile("./schedulers.json", schedulers)
}
