package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type block struct {
	Height    int   `json:"height"`
	Timestamp int64 `json:"timestamp"`
}

var client = http.Client{Timeout: 20 * time.Second}

func getJSON(url string, out any) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: HTTP %d", url, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func main() {
	api := flag.String("api", "https://mempool.space/api", "mempool API")
	n := flag.Int("blocks", 500, "number of recent blocks")
	flag.Parse()

	var tip int
	if err := getJSON(*api+"/blocks/tip/height", &tip); err != nil {
		log.Fatal(err)
	}
	seen := map[int]block{}
	for h := tip; len(seen) < *n+1 && h >= 0; {
		var page []block
		if err := getJSON(fmt.Sprintf("%s/v1/blocks/%d", *api, h), &page); err != nil {
			log.Fatal(err)
		}
		if len(page) == 0 {
			break
		}
		for _, b := range page {
			seen[b.Height] = b
		}
		h = page[len(page)-1].Height - 1
		time.Sleep(100 * time.Millisecond)
	}
	blocks := make([]block, 0, len(seen))
	for _, b := range seen {
		blocks = append(blocks, b)
	}
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].Height < blocks[j].Height })
	if len(blocks) > *n+1 {
		blocks = blocks[len(blocks)-*n-1:]
	}
	iv := make([]float64, 0, len(blocks)-1)
	for i := 1; i < len(blocks); i++ {
		iv = append(iv, float64(blocks[i].Timestamp-blocks[i-1].Timestamp)/60)
	}

	s := Summarize(iv)
	fmt.Printf("blocks %d-%d (%d intervals)\n", blocks[0].Height, blocks[len(blocks)-1].Height, s.N)
	fmt.Printf("mean %.2f min, median %.2f min (exponential predicts %.2f), p90 %.1f min, max %.1f min\n",
		s.Mean, s.Median, s.Mean*0.6931, s.P90, s.Max)
	fmt.Printf("negative intervals (timestamp before parent): %d\n", s.Negative)
	for _, t := range []float64{20, 30, 60} {
		obs := 0
		for _, v := range iv {
			if v > t {
				obs++
			}
		}
		fmt.Printf("gaps > %2.0f min: observed %3d, expected %5.1f\n", t, obs, ExpectedAbove(t, s.Mean)*float64(s.N))
	}

	width, max := 2.0, 40.0
	h := Histogram(iv, width, max)
	fmt.Println("\nminutes     observed  expected")
	for i, c := range h {
		lo := float64(i) * width
		var exp float64
		label := fmt.Sprintf("%2.0f-%-2.0f", lo, lo+width)
		if i == len(h)-1 {
			exp = ExpectedAbove(lo, s.Mean) * float64(s.N)
			label = fmt.Sprintf("%2.0f+   ", lo)
		} else {
			exp = (ExpectedAbove(lo, s.Mean) - ExpectedAbove(lo+width, s.Mean)) * float64(s.N)
		}
		fmt.Printf("%s  %5d  %7.1f  %s\n", label, c, exp, strings.Repeat("#", c*60/s.N))
	}
}
