package arrange

import (
	"encoding/csv"
	"os"
	"strings"
)

// Diversity weights: department > role > country. BaseConnectionValue ensures
// every new unique connection outweighs any diversity bonus, so the greedy
// still prioritises raw new pairings.
const (
	BaseConnectionValue = 1000
	WeightDept          = 100
	WeightRole          = 10
	WeightCountry       = 1
)

type Person struct {
	Name       string
	Department string
	Role       string
	Country    string
	Dept       int
	RoleID     int
	CountryID  int
}

type Dataset struct {
	People    []Person
	N         int
	PairScore [][]int // symmetric, PairScore[i][j] == PairScore[j][i]
}

func LoadPeople(path string) (*Dataset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &Dataset{}, nil
	}

	header := rows[0]
	colIdx := func(name string) int {
		for i, h := range header {
			if strings.TrimSpace(h) == name {
				return i
			}
		}
		return -1
	}
	nameCol := colIdx("name")
	deptCol := colIdx("department")
	roleCol := colIdx("role")
	countryCol := colIdx("country")

	deptIdx, roleIdx, countryIdx := map[string]int{}, map[string]int{}, map[string]int{}
	intern := func(m map[string]int, s string) int {
		if v, ok := m[s]; ok {
			return v
		}
		id := len(m)
		m[s] = id
		return id
	}
	get := func(row []string, c int) string {
		if c < 0 || c >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[c])
	}

	people := make([]Person, 0, len(rows)-1)
	for _, row := range rows[1:] {
		name := get(row, nameCol)
		if name == "" {
			continue
		}
		dept := get(row, deptCol)
		role := get(row, roleCol)
		country := get(row, countryCol)
		people = append(people, Person{
			Name:       name,
			Department: dept,
			Role:       role,
			Country:    country,
			Dept:       intern(deptIdx, dept),
			RoleID:     intern(roleIdx, role),
			CountryID:  intern(countryIdx, country),
		})
	}

	n := len(people)
	pair := make([][]int, n)
	for i := 0; i < n; i++ {
		pair[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			s := BaseConnectionValue
			if people[i].Dept != people[j].Dept {
				s += WeightDept
			}
			if people[i].RoleID != people[j].RoleID {
				s += WeightRole
			}
			if people[i].CountryID != people[j].CountryID {
				s += WeightCountry
			}
			pair[i][j] = s
		}
	}
	return &Dataset{People: people, N: n, PairScore: pair}, nil
}

func (ds *Dataset) AveragePairValue() float64 {
	if ds.N < 2 {
		return float64(BaseConnectionValue)
	}
	total := 0
	count := 0
	for i := 0; i < ds.N; i++ {
		for j := i + 1; j < ds.N; j++ {
			total += ds.PairScore[i][j]
			count++
		}
	}
	return float64(total) / float64(count)
}
