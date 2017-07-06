package main

import (
	"flag"
	"fmt"
)

func usage() {
	fmt.Println("Usage: toy COMMAND [arg...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    bench")
	fmt.Println("    report")
	fmt.Println()
	fmt.Println("run 'toy COMMAND --help' for more information on a command")
}

func main() {
	flag.Usage = usage
	flag.Parse()
}
