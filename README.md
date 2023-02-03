# terse
Output randomly sampled lines from input stream or file. Uses simple [reservoir sampling](http://www.cs.umd.edu/~samir/498/vitter.pdf) algorithm to process input with linear time complexity. Suitable for processing streams, seeing each line only once. Retains relative order of lines.

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

Comparison against `shuf -n`  on real data: 5.1GB nginx log with 17451712  lines in it.

```
root@logger:~# ls -lh /var/log/remote/nginx/2023_02_02_18.log
-rw-r----- 1 root logs 5.1G Feb  2 18:59 /var/log/remote/nginx/2023_02_02_18.log
root@logger:~# wc -l /var/log/remote/nginx/2023_02_02_18.log
17451712 /var/log/remote/nginx/2023_02_02_18.log
root@logger:~# time terse -i /var/log/remote/nginx/2023_02_02_18.log -n 25 > /dev/null

real    0m2.656s
user    0m1.315s
sys     0m1.372s
root@logger:~# time shuf -n 25 /var/log/remote/nginx/2023_02_02_18.log > /dev/null

real    0m22.784s
user    0m21.059s
sys     0m1.703s
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

A docker image is available as well. Here is an example of running terse in a pipeline with docker:

```sh
seq 5 | docker run -i --rm yarmak/terse
```

## Synopsis

```
> terse -h
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
