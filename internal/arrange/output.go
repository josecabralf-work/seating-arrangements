package arrange

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func WriteSeatings(seatings Seatings, ds *Dataset, filename string) error {
	var sb strings.Builder
	for r, rt := range seatings {
		fmt.Fprintf(&sb, "Day %d:\n", r+1)
		for t, table := range rt {
			fmt.Fprintf(&sb, "Table %d:\n", t+1)
			ids := append([]int(nil), table...)
			sort.Slice(ids, func(i, j int) bool {
				return ds.People[ids[i]].Name < ds.People[ids[j]].Name
			})
			for _, p := range ids {
				fmt.Fprintf(&sb, "  - %s\n", ds.People[p].Name)
			}
		}
		sb.WriteByte('\n')
	}
	return os.WriteFile(filename, []byte(sb.String()), 0o644)
}
