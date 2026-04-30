package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"arrange/internal/arrange"
)

func envInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func envStr(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}

func main() {
	peopleFile := envStr("PEOPLE", "people.txt")
	tables := envInt("TABLES", 30)
	seats := envInt("SEATS", 8)
	days := envInt("DAYS", 3)
	attempts := envInt("ATTEMPTS", 50)
	workers := envInt("WORKERS", runtime.NumCPU())
	output := envStr("OUTPUT", "seatings.txt")

	ds, err := arrange.LoadPeople(peopleFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load error:", err)
		os.Exit(1)
	}
	avgPair := ds.AveragePairValue()

	type result struct {
		seatings arrange.Seatings
		unique   int
		weighted int
	}

	results := make(chan result, attempts)
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	var done int64
	start := time.Now()

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			seed1 := uint64(time.Now().UnixNano()) ^ uint64(idx)*0x9e3779b97f4a7c15
			seed2 := uint64(idx)*2862933555777941757 + 3037000493
			rng := rand.New(rand.NewPCG(seed1, seed2))

			seatings := arrange.CreateSittings(ds, tables, seats, days, avgPair, rng)
			seatings, u, w := arrange.ImproveBySwaps(seatings, ds)
			n := atomic.AddInt64(&done, 1)
			fmt.Printf("Attempt %d/%d done: %d unique, %d weighted\n", n, attempts, u, w)
			results <- result{seatings, u, w}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	bestU, bestW := -1, math.MinInt
	var bestSeatings arrange.Seatings
	for r := range results {
		if r.unique > bestU || (r.unique == bestU && r.weighted > bestW) {
			bestU = r.unique
			bestW = r.weighted
			bestSeatings = r.seatings
		}
	}

	if err := arrange.WriteSeatings(bestSeatings, ds, output); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}
	fmt.Printf("Best arrangement: %d unique connections, weighted value %d. (%.2fs over %d attempts, %d workers)\n",
		bestU, bestW, time.Since(start).Seconds(), attempts, workers)
}
