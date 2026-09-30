# btc-block-intervals

Fetches recent Bitcoin blocks and compares the time between them with what a Poisson process
predicts. The target is one block every 10 minutes, but because mining is a Poisson process the gaps
follow an exponential distribution, and gaps of 30 minutes or more are normal.

```bash
go run .
go run . -blocks 1000
```

It prints the mean, median and 90th percentile interval, the number of gaps longer than 20, 30 and 60
minutes next to the number the exponential model expects (`P(gap > t) = exp(-t / mean)`), and a
histogram of observed and expected counts in 2-minute buckets.

```
$ go run . -blocks 300
blocks 967987-968287 (300 intervals)
mean 10.20 min, median 7.63 min (exponential predicts 7.07), p90 22.4 min, max 53.0 min
negative intervals (timestamp before parent): 8
gaps > 20 min: observed  40, expected  42.2
gaps > 30 min: observed  13, expected  15.9
gaps > 60 min: observed   0, expected   0.8

minutes     observed  expected
 0-2      42     53.4  ########
 2-4      37     43.9  #######
 4-6      38     36.1  #######
...
```

Two things stand out in the output:

- The median is about 7 minutes (`ln 2 x 10 min`), not 10, because in an exponential distribution
  short gaps are the most common.
- A few intervals are negative. A block's timestamp only has to be later than the median of the
  previous 11 blocks, so a miner's clock can put a block before its parent.

Data comes from `GET /api/blocks/tip/height` and `GET /api/v1/blocks/:height`, 15 blocks per call.

```bash
go test ./...
```
