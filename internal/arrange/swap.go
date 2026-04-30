package arrange

// scoreState tracks how many days each pair appears in. The total
// (unique, weighted) only depends on whether each pair appears at all,
// so swap deltas are O(seats) instead of O(days*tables*seats^2).
type scoreState struct {
	pairCount [][]int // upper-triangular: index with i < j
	unique    int
	weighted  int
}

func newScoreState(seatings Seatings, ds *Dataset) *scoreState {
	n := ds.N
	pc := make([][]int, n)
	for i := range pc {
		pc[i] = make([]int, n)
	}
	s := &scoreState{pairCount: pc}
	for _, rt := range seatings {
		for _, table := range rt {
			for i := 0; i < len(table); i++ {
				for j := i + 1; j < len(table); j++ {
					a, b := table[i], table[j]
					if a > b {
						a, b = b, a
					}
					if pc[a][b] == 0 {
						s.unique++
						s.weighted += ds.PairScore[a][b]
					}
					pc[a][b]++
				}
			}
		}
	}
	return s
}

func (s *scoreState) changePair(x, y, delta int, ds *Dataset) {
	if x == y {
		return
	}
	if x > y {
		x, y = y, x
	}
	before := s.pairCount[x][y]
	after := before + delta
	s.pairCount[x][y] = after
	if before == 0 && after > 0 {
		s.unique++
		s.weighted += ds.PairScore[x][y]
	} else if before > 0 && after == 0 {
		s.unique--
		s.weighted -= ds.PairScore[x][y]
	}
}

// applySwap swaps person at ti[ai] with person at tj[bj] (same day) and
// updates the score state incrementally. Returns the previous (unique, weighted).
func (s *scoreState) applySwap(ti, tj []int, ai, bj int, ds *Dataset) (int, int) {
	a := ti[ai]
	b := tj[bj]
	prevU, prevW := s.unique, s.weighted
	for k, x := range ti {
		if k == ai {
			continue
		}
		s.changePair(a, x, -1, ds)
		s.changePair(b, x, +1, ds)
	}
	for k, y := range tj {
		if k == bj {
			continue
		}
		s.changePair(b, y, -1, ds)
		s.changePair(a, y, +1, ds)
	}
	ti[ai], tj[bj] = b, a
	return prevU, prevW
}

func (s *scoreState) revertSwap(ti, tj []int, ai, bj int, ds *Dataset) {
	// After applySwap, ti[ai] is the formerly-tj person and vice versa.
	b := ti[ai]
	a := tj[bj]
	ti[ai], tj[bj] = a, b
	for k, x := range ti {
		if k == ai {
			continue
		}
		s.changePair(a, x, +1, ds)
		s.changePair(b, x, -1, ds)
	}
	for k, y := range tj {
		if k == bj {
			continue
		}
		s.changePair(b, y, +1, ds)
		s.changePair(a, y, -1, ds)
	}
}

// ImproveBySwaps hill-climbs: swap pairs of people between tables on the same
// day whenever doing so improves (unique, weighted) lexicographically.
func ImproveBySwaps(seatings Seatings, ds *Dataset) (Seatings, int, int) {
	state := newScoreState(seatings, ds)
	for {
		if !findImprovingSwap(seatings, state, ds) {
			return seatings, state.unique, state.weighted
		}
	}
}

func findImprovingSwap(seatings Seatings, state *scoreState, ds *Dataset) bool {
	for _, rt := range seatings {
		nT := len(rt)
		for i := 0; i < nT; i++ {
			for j := i + 1; j < nT; j++ {
				ti := rt[i]
				tj := rt[j]
				for ai := 0; ai < len(ti); ai++ {
					for bj := 0; bj < len(tj); bj++ {
						prevU, prevW := state.applySwap(ti, tj, ai, bj, ds)
						if state.unique > prevU || (state.unique == prevU && state.weighted > prevW) {
							return true
						}
						state.revertSwap(ti, tj, ai, bj, ds)
					}
				}
			}
		}
	}
	return false
}
