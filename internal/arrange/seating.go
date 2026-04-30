package arrange

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// Seatings: [round][table] -> slice of person IDs.
type Seatings [][][]int

func CreateSittings(ds *Dataset, tables, seats, rounds int, avgPair float64, rng *rand.Rand) Seatings {
	n := ds.N
	if n == 0 || tables <= 0 || seats <= 0 || rounds <= 0 {
		return nil
	}
	if n > tables*seats {
		panic(fmt.Sprintf("Not enough seats for all people: %d > %d", n, tables*seats))
	}
	needed := n / seats
	if n%seats > 0 {
		needed++
	}
	if tables > needed {
		tables = needed
	}

	connections := make([][]bool, n)
	for i := range connections {
		connections[i] = make([]bool, n)
	}

	seatings := make(Seatings, rounds)
	for r := range seatings {
		seatings[r] = make([][]int, tables)
		for t := range seatings[r] {
			seatings[r][t] = make([]int, 0, seats)
		}
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for r := 0; r < rounds; r++ {
		rng.Shuffle(n, func(i, j int) { order[i], order[j] = order[j], order[i] })
		for _, p := range order {
			bestScore := math.Inf(-1)
			bestTable := -1
			for t := 0; t < tables; t++ {
				seated := seatings[r][t]
				rem := seats - len(seated)
				if rem <= 0 {
					continue
				}
				imm := 0
				for _, o := range seated {
					if !connections[p][o] {
						imm += ds.PairScore[p][o]
					}
				}
				future := float64(rem-1) * avgPair
				score := float64(imm) + future
				if score > bestScore {
					bestScore = score
					bestTable = t
				}
			}
			if bestTable == -1 {
				panic("Not enough seats available")
			}
			for _, o := range seatings[r][bestTable] {
				connections[p][o] = true
				connections[o][p] = true
			}
			seatings[r][bestTable] = append(seatings[r][bestTable], p)
		}
	}
	return seatings
}
