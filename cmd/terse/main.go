package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ProgName = "terse"
)

var (
	version = "undefined"

	showVersion = flag.Bool("version", false, "show program version and exit")
	limit  = flag.Int("n", 25, "number of lines to sample")
)

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "%s [OPTION]...\n", ProgName)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	flag.PrintDefaults()
}

func run() int {
	flag.CommandLine.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return 0
	}

	return 0
}

func main() {
	os.Exit(run())
}
