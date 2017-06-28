package main

import (
	"testing"

	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/toy/framework"
)

func TestWriteSchedulersJSON(t *testing.T) {
	schedulers := []framework.Scheduler{
		{
			N:           1,
			C:           10,
			QPS:         1000,
			Name:        "Connect",
			Display:     true,
			PrintReport: false,
		},
		{
			N:           1,
			C:           10,
			QPS:         1000,
			Name:        "Prepare",
			Display:     true,
			PrintReport: false,
		},
		{
			N:           100,
			C:           100,
			QPS:         5000,
			Name:        "Publish",
			Display:     true,
			PrintReport: true,
		},
		{
			N:           1,
			C:           10,
			QPS:         1000,
			Name:        "Disconnect",
			Display:     true,
			PrintReport: false,
		},
	}
	jsoncfg.WriteToFile("./schedulers.json", schedulers)
}
