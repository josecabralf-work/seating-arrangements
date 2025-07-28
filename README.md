# Induction Seatings Arrangements

This repository provides a Python script (`arrange.py`) to generate optimal seating arrangements for group events, such as induction ceremonies, where the goal is to maximize unique connections between participants across multiple rounds.

## Features
- Distributes people across a specified number of tables and seats per table for several rounds.
- Maximizes the number of unique connections (i.e., ensures participants sit with as many different people as possible).
- Outputs the best seating arrangement found after several attempts.
- Easy configuration via environment variables or by editing the script.

## How It Works
1. **Input**: Reads a list of people from a text file (default: `people.txt`, one name per line).
2. **Parameters**: Number of tables, seats per table, and rounds can be set via environment variables or default values.
3. **Algorithm**: Randomly shuffles and assigns people to tables for each round, tracking connections to maximize unique pairings.
4. **Output**: Writes the best seating arrangement to `seatings.txt` and prints the number of unique connections achieved.

## Usage

### 1. Prepare the People List
Create a `people.txt` file in the same directory, listing each participant on a separate line:

```
Alice
Bob
Charlie
Diana
...etc.
```

### 2. Run the Script
You can run the script directly with Python 3:

```bash
python3 arrange.py
```

#### Optional: Set Parameters
You can override the defaults using environment variables:

- `PEOPLE`: Path to the people list file (default: `people.txt`)
- `TABLES`: Number of tables (default: 15)
- `SEATS`: Number of seats per table (default: 6)
- `ROUNDS`: Number of rounds (default: 2)

Example:

```bash
TABLES=10 SEATS=5 ROUNDS=3 python3 arrange.py
```

### 3. View the Output
The script will generate a `seatings.txt` file with the seating arrangements for each round and table. It will also print the number of unique connections achieved.

## Example Output (`seatings.txt`)
```
Round 1:
Table 1:
  - Alice
  - Bob
  ...
Table 2:
  - Charlie
  - Diana
  ...

Round 2:
Table 1:
  - ...
```

## Requirements
- Python 3.x

## License
See [LICENSE](LICENSE) for details.
