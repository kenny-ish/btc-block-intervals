package main

import (
	"math"
	"sort"
)

type Stats struct {
	N        int
	Mean     float64
	Median   float64
	P90      float64
	Max      float64
	Negative int
}

func quantile(sorted []float64, q float64) float64 {
	if len(sorted) == 0 {
		return math.NaN()
	}
	pos := q * float64(len(sorted)-1)
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))
	return sorted[lo] + (sorted[hi]-sorted[lo])*(pos-float64(lo))
}

// Summarize takes intervals in minutes.
func Summarize(iv []float64) Stats {
	s := append([]float64(nil), iv...)
	sort.Float64s(s)
	var sum float64
	neg := 0
	for _, v := range s {
		sum += v
		if v < 0 {
			neg++
		}
	}
	return Stats{N: len(s), Mean: sum / float64(len(s)), Median: quantile(s, 0.5), P90: quantile(s, 0.9), Max: s[len(s)-1], Negative: neg}
}

// ExpectedAbove is the exponential model's share of gaps longer than t.
func ExpectedAbove(t, mean float64) float64 {
	return math.Exp(-t / mean)
}

// Histogram counts intervals in buckets of `width` minutes up to `max`; the last bucket is open ended.
func Histogram(iv []float64, width, max float64) []int {
	n := int(max/width) + 1
	h := make([]int, n)
	for _, v := range iv {
		i := int(math.Max(0, v) / width)
		if i >= n {
			i = n - 1
		}
		h[i]++
	}
	return h
}
