import math
import random


def create_sittings(people: list[str], tables: int, seats: int, rounds: int):
    if not people or tables <= 0 or seats <= 0 or rounds <= 0:
        return []
    if len(people) > tables * seats:
        raise ValueError(f"Not enough seats for all people: {len(people)} > {tables * seats}")

    tables = min(tables, len(people) // seats + (1 if len(people) % seats > 0 else 0))
    connections = { person: set() for person in people }
    seatings = [[set() for _ in range(tables)] for _ in range(rounds)]

    for round in range(rounds):
        random.shuffle(people)
        for person in people:
            max_potential = [-math.inf, -1] # (score, table index)
            for table in range(tables):
                potential = seats - len(seatings[round][table]) - 1 # the -1 is for the person themselves
                remaining_seats = seats - len(seatings[round][table])
                if potential > max_potential[0] and remaining_seats > 0:
                    max_potential = [potential, table]

            if max_potential[1] == -1:
                raise ValueError(f"Not enough seats available")
            
            connections[person] = connections[person].union(seatings[round][max_potential[1]])
            for p in seatings[round][max_potential[1]]:
                connections[p].add(person)

            # We add the person after updating connections to avoid self-connection
            seatings[round][max_potential[1]].add(person)

    unique_connections = sum(len(connections[person]) for person in people) / 2

    return seatings, unique_connections


def write_seatings_to_txt(seatings, filename="seatings.txt"):
    with open(filename, "w") as f:
        for round_index, round_tables in enumerate(seatings, start=1):
            f.write(f"Round {round_index}:\n")
            for table_index, table in enumerate(round_tables, start=1):
                f.write(f"Table {table_index}:\n")
                for person in sorted(table):
                    f.write(f"  - {person}\n")
            f.write("\n")


if __name__ == "__main__":
    import os

    people_file = os.getenv("PEOPLE", "people.txt")
    tables = int(os.getenv("TABLES", 15))
    seats = int(os.getenv("SEATS", 6))
    rounds = int(os.getenv("ROUNDS", 2))
    
    with open(people_file, "r") as f:
        people = [line.strip() for line in f if line.strip()]

    max_unique_connections = -math.inf
    best_seatings = []
    for i in range(10):  # Run multiple times to find the best configuration
        seatings, unique_connections = create_sittings(people, tables, seats, rounds)
        if unique_connections > max_unique_connections:
            max_unique_connections = unique_connections
            best_seatings = seatings

    write_seatings_to_txt(best_seatings)
    print(f"Best seating arrangement found with {max_unique_connections} unique connections.")