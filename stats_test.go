package main

import (
	"math"
	"testing"
)

func TestSummarize(t *testing.T) {
	s := Summarize([]float64{10, 2, -1, 30, 9})
	if s.N != 5 || s.Median != 9 || s.Max != 30 || s.Negative != 1 {
		t.Fatalf("%+v", s)
	}
	if math.Abs(s.Mean-10) > 1e-12 {
		t.Fatalf("mean %v", s.Mean)
	}
}

func TestExpectedAbove(t *testing.T) {
	if math.Abs(ExpectedAbove(10, 10)-math.Exp(-1)) > 1e-12 {
		t.Fatal("P(>mean) should be e^-1")
	}
	// median of an exponential is ln2 * mean
	if math.Abs(ExpectedAbove(math.Ln2*10, 10)-0.5) > 1e-12 {
		t.Fatal("median")
	}
}

func TestHistogram(t *testing.T) {
	h := Histogram([]float64{-1, 0.5, 1.9, 2, 5, 100}, 2, 6)
	want := []int{3, 1, 1, 1}
	for i := range want {
		if h[i] != want[i] {
			t.Fatalf("got %v want %v", h, want)
		}
	}
}

func TestQuantileInterpolates(t *testing.T) {
	if q := quantile([]float64{0, 10}, 0.9); math.Abs(q-9) > 1e-12 {
		t.Fatalf("q %v", q)
	}
}
