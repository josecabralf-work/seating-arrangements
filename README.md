# Induction Seating Arrangements

Generate multi-day seating plans that maximise how many different people meet
each other.

Given a list of people, the tool produces a seating chart for each day that
maximises the number of **unique pairwise connections** (people who share a table
at least once across the whole event), then breaks ties by favouring **diverse**
tables.

## How it works

The objective is lexicographic:

1. **Unique connections**: the count of distinct person-pairs that share a table
   on any day. This is the primary score.
2. **Weighted diversity**: among arrangements with equal unique connections,
   prefer pairings across different attributes. Each pair contributes a base value
   plus bonuses for differing department, role, and country (see weights below).

Diversity weights ([internal/arrange/dataset.go](internal/arrange/dataset.go)):

| Attribute  | Weight | Notes                                    |
| ---------- | -----: | ---------------------------------------- |
| Connection |  1000  | base value for any new unique pairing    |
| Department |   100  | bonus if the two people differ           |
| Role       |    10  | bonus if the two people differ           |
| Country    |     1  | bonus if the two people differ           |

The base connection value (1000) dominates all diversity bonuses combined, so the
solver never trades away a new connection to gain diversity. Diversity only breaks
ties.

### Algorithm

1. **Greedy construction** ([seating.go](internal/arrange/seating.go)): for each
   day, people are shuffled and seated one at a time at the table that gives the
   best immediate-plus-expected score. Tables are kept balanced: every table must
   reach the floor size before any table grows toward the ceiling, so people are
   spread evenly rather than piling onto a few tables.
2. **Hill climbing** ([swap.go](internal/arrange/swap.go)): pairs of people are
   swapped between tables on the same day whenever the swap improves the
   `(unique, weighted)` score. Swap deltas are computed incrementally in O(seats),
   not by rescoring the whole plan.
3. **Parallel attempts** ([cmd/arrange/main.go](cmd/arrange/main.go)): many
   independent attempts run concurrently across worker goroutines, each with its
   own random seed. The best result wins.

## Requirements

- Go 1.22 or newer

## Usage

```bash
make build      # builds bin/arrange
make run        # builds, loads .env, and runs
```

Or run directly:

```bash
go build -o bin/arrange ./cmd/arrange
PEOPLE=people.txt TABLES=30 SEATS=8 DAYS=3 ./bin/arrange
```

`make run` sources a local `.env` file for configuration before running.

### Input

A CSV formatted file (default `people.txt`) with a header row. Recognised columns:
`name`, `department`, `role`, `country`. Only `name` is required; missing columns
are treated as empty and simply contribute no diversity bonus.

```csv
name,department,role,country
Alice Johnson,Engineering,Senior,USA
Brian Smith,Marketing,Manager,UK
Catherine Lee,Product,Mid,Canada
```

### Output

A text file (default `seatings.txt`) listing the tables for each day, with the
people seated at each, sorted by name:

```
Day 1:
Table 1:
  - Alice Johnson
  - Brian Smith
  ...
```

The program also prints per-attempt progress and a final summary to stdout:

```
Best arrangement: 412 unique connections, weighted value 458300. (1.23s over 50 attempts, 8 workers)
```

## Configuration

All options are set via environment variables (read by
[cmd/arrange/main.go](cmd/arrange/main.go)):

| Variable   | Default            | Description                                   |
| ---------- | ------------------ | --------------------------------------------- |
| `PEOPLE`   | `people.txt`       | Path to the input CSV                         |
| `TABLES`   | `30`               | Maximum number of tables per day              |
| `SEATS`    | `8`                | Seats per table                               |
| `DAYS`     | `3`                | Number of days/rounds to schedule             |
| `ATTEMPTS` | `50`               | Independent attempts; best result is kept     |
| `WORKERS`  | number of CPUs     | Concurrent worker goroutines                  |
| `OUTPUT`   | `seatings.txt`     | Path to the output file                       |

If `TABLES` exceeds what's needed to seat everyone given `SEATS`, the count is
reduced automatically. The program requires `TABLES × SEATS ≥ number of people`.

## License

[MIT](LICENSE)
