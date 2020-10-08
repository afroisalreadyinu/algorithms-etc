Solutions to the eight puzzle with various search algorithms in different
languages. Here is a comparison of the performances using
[hyperfine](https://github.com/sharkdp/hyperfine):

Python:

```
$ hyperfine 'python eight.py start.txt'
Benchmark #1: python eight.py start.txt
  Time (mean ± σ):     246.2 ms ±   8.7 ms    [User: 236.3 ms, System: 8.3 ms]
  Range (min … max):   239.6 ms … 270.4 ms    11 runs
```
Go:

```
$ go build eight.go && hyperfine './eight start.txt'
Benchmark #1: ./eight start.txt
  Time (mean ± σ):      15.7 ms ±   1.4 ms    [User: 21.1 ms, System: 2.3 ms]
  Range (min … max):    14.1 ms …  23.3 ms    124 runs
```

C:

```
$ make eight && hyperfine './eight start.txt'
cc     eight.c   -o eight
Benchmark #1: ./eight start.txt
  Time (mean ± σ):      10.3 ms ±   1.3 ms    [User: 8.3 ms, System: 2.1 ms]
  Range (min … max):     8.6 ms …  15.0 ms    204 runs

```
