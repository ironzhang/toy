package main

import (
	"flag"
	"fmt"

	"github.com/ironzhang/toy/framework/command"
)

func usage() {
	fmt.Println("Usage: toy COMMAND [arg...]")
	fmt.Println()
	fmt.Println("A extensible benchmark tool")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    bench\tdo benchmark test")
	fmt.Println("    report\tmake benchmark report")
	fmt.Println()
	fmt.Println("run 'toy COMMAND --help' for more information on a command")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		usage()
		return
	}

	switch args[0] {
	case "bench":
		RunBenchCmd(args[1:])
	case "report":
		RunReportCmd(args[1:])
	default:
		usage()
	}
}

func RunBenchCmd(args []string) {
	var cmd command.BenchCmd
	if err := cmd.Run(args); err != nil {
		fmt.Println(err)
	}
}

func RunReportCmd(args []string) {
	var cmd = command.ReportCmd{Template: "./framework/report/builders/html-report/templates/report.template"}
	if err := cmd.Run(args); err != nil {
		fmt.Println(err)
	}
}
