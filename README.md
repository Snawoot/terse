# terse
Output randomly sampled lines from input stream or file. Uses simple reservoir sampling algorithm to process input with linear time complexity. Suitable for processing streams, seeing each line only once. Retains relative order of lines.

## Usage example

```
> seq 1000000 | terse -n 5
349893
539678
576919
738393
758023
```

## Performance

```
> time seq 100000000 | bin/terse -n 5
41432706
56746242
61118996
70135895
93968158

real	0m5,106s
user	0m5,676s
sys	0m0,430s
```

It processes about tens of millions of lines per second on modern computer. Most likely I/O will become bottleneck in such sampling rather than application performance will be an issue.

## Installation

#### Binaries

Pre-built binaries are available [here](https://github.com/Snawoot/terse/releases/latest).

#### Build from source

Alternatively, you may install terse from source. Run the following within the source directory:

```
make install
```

#### Docker

A docker image is available as well. Here is an example of running terse as a background service:

```sh
docker run -it --rm yarmak/terse
```

## Synopsis

```
$ terse -h
Usage:

terse [OPTION]...

Options:
  -buffered
    	buffer control (default true)
  -i string
    	use input file instead of stdin
  -n int
    	number of lines to sample (default 25)
  -o string
    	use output file instead of stdout
  -seed value
    	use fixed random seed (default is a value from CSPRNG)
  -version
    	show program version and exit
  -z	line delimiter is NUL, not newline
```
