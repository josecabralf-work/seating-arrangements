# Induction Seatings Arrangements

This repository provides a Python script (`arrange.py`) to generate optimal seating arrangements for group events, such as induction ceremonies, where the goal is to maximize **unique** connections between participants across multiple rounds while also **boosting the value** of each connection by diversity of department, role, and country.

## Features
- Distributes people across a specified number of tables and seats per table for several rounds.
- Maximizes the number of unique connections (i.e., ensures participants sit with as many different people as possible).
- Weights each connection by how different the two participants are, in the order **department > role > country**.
- Outputs the best seating arrangement found after several attempts.
- Easy configuration via environment variables.

## How It Works
1. **Input**: Reads a CSV list of people from a text file (default: `people.txt`) with columns `name, department, role, country`.
2. **Parameters**: Number of tables, seats per table, rounds, and attempts can be set via environment variables.
3. **Scoring**: Every unique pairing is worth a large base value (so raw unique connections remain the primary goal). Added on top are diversity bonuses:
   - Different **department** → +100
   - Different **role** → +10
   - Different **country** → +1
   The greedy placement uses this value on each seat assignment, and runs are ranked lexicographically by `(unique_connections, weighted_value)`.
4. **Algorithm**: Two-stage per attempt.
   - **Greedy placement**: for each round, shuffles participants and, for each person, picks the table that maximizes the sum of immediate diversity-weighted value with existing seat-mates (only counting *new* connections) plus the expected value of the remaining empty seats.
   - **Swap local search**: after greedy, iteratively swaps pairs of people between tables in the same round whenever doing so strictly improves `(unique_connections, weighted_value)`. Continues until no improving swap exists.
5. **Output**: Writes the best seating arrangement to `seatings.txt` (annotated with each person's attributes) and prints the number of unique connections along with the total weighted value.

## Usage

### 1. Prepare the People List
Create a `people.txt` file in the same directory. The first line is a CSV header; one participant per subsequent line:

```
name,department,role,country
Alice Johnson,Engineering,Senior,USA
Brian Smith,Marketing,Manager,UK
Catherine Lee,Product,Mid,Canada
...
```

All four columns are required. Use any consistent labels you like — the script only cares whether two values are equal or different.

### 2. Run the Script
Run directly with Python 3:

```bash
python3 arrange.py
```

#### Optional: Set Parameters
Override the defaults using environment variables:

- `PEOPLE`: Path to the people list file (default: `people.txt`)
- `TABLES`: Number of tables (default: 15)
- `SEATS`: Number of seats per table (default: 6)
- `ROUNDS`: Number of rounds (default: 2)
- `ATTEMPTS`: Number of random attempts to try (default: 50)

Example:

```bash
TABLES=10 SEATS=5 ROUNDS=3 ATTEMPTS=200 python3 arrange.py
```

### 3. View the Output
The script generates a `seatings.txt` file with the seating arrangements for each round and table, annotated with department, role, and country. It also prints the number of unique connections and the total weighted value achieved.

## Example Output (`seatings.txt`)
```
Round 1:
Table 1:
  - Alice Johnson (Engineering, Senior, USA)
  - Brian Smith (Marketing, Manager, UK)
  ...
```

## Tuning the Weights
If department/role/country should be weighted differently for your event, edit the constants at the top of `arrange.py`:

```python
BASE_CONNECTION_VALUE = 1000
WEIGHT_DEPT = 100
WEIGHT_ROLE = 10
WEIGHT_COUNTRY = 1
```

Keep `BASE_CONNECTION_VALUE` strictly larger than the sum of all diversity weights if you want unique-pair count to stay the dominant objective.

## Requirements
- Python 3.9+

## License
See [LICENSE](LICENSE) for details.
