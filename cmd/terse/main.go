package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

const (
	ProgName = "terse"
)

var (
	version = "undefined"

	showVersion  = flag.Bool("version", false, "show program version and exit")
	limit        = flag.Int("n", 25, "number of lines to sample")
	nulDelimiter = flag.Bool("z", false, "line delimiter is NUL, not newline")
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

	scanner := bufio.NewScanner(os.Stdin)
	if *nulDelimiter {
		scanner.Split(scanZeroTerminatedLines)
	}
	for scanner.Scan() {
		fmt.Printf("read: %q\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v", err)
	}

	return 0
}

func main() {
	os.Exit(run())
}

// scanZeroTerminatedLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing zero-byte end-of-line marker.
// The returned line may be empty.
// The last non-empty line of input will be returned even if it has no newline.
func scanZeroTerminatedLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full zero-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
