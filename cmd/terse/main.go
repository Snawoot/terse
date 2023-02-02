package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Snawoot/terse/reservoir"
	"github.com/Snawoot/terse/rng"
)

const (
	ProgName = "terse"
)

var (
	version = "undefined"

	showVersion    = flag.Bool("version", false, "show program version and exit")
	limit          = flag.Int("n", 25, "number of lines to sample")
	nulDelimiter   = flag.Bool("z", false, "line delimiter is NUL, not newline")
	inputFilename  = flag.String("i", "", "use input file instead of stdin")
	outputFilename = flag.String("o", "", "use output file instead of stdout")
	buffered       = flag.Bool("buffered", true, "buffer control")
	seed           *int64
)

func init() {
	flag.Func("seed", "use fixed random seed (default is a value from CSPRNG)", func(val string) error {
		seedVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse seed value: %w", err)
		}
		seed = &seedVal
		return nil
	})
}

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

	if *limit < 0 {
		fmt.Fprintln(os.Stderr, "error: negative limit value")
		usage()
		return 2
	}

	var input io.Reader = os.Stdin
	if *inputFilename != "" {
		f, err := os.Open(*inputFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open input file: %v\n", err)
			return 3
		}
		defer f.Close()
		input = f
	}
	var output io.Writer = os.Stdout
	if *outputFilename != "" {
		f, err := os.Create(*outputFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open output file: %v\n", err)
			return 4
		}
		defer func() {
			f.Sync()
			f.Close()
		}()
		output = f
	}

	if *buffered {
		buffin := bufio.NewReader(input)
		input = buffin
		buffout := bufio.NewWriter(output)
		defer buffout.Flush()
		output = buffout
	}

	r := reservoir.NewReservoir[string](*limit, rng.NewRNG(seed))

	scanner := bufio.NewScanner(input)
	if *nulDelimiter {
		scanner.Split(scanZeroTerminatedLines)
	}

	for scanner.Scan() {
		r.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
	}

	delimiter := []byte{'\n'}
	if *nulDelimiter {
		delimiter = []byte{0}
	}
	for _, line := range r.Items() {

		if _, err := output.Write([]byte(line)); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}

		if _, err := output.Write(delimiter); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}
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
