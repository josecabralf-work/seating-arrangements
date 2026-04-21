import csv
import math
import random
from itertools import combinations


# Diversity weights: department matters most, then role, then country.
# BASE ensures every new unique connection is worth far more than any diversity
# bonus, so the greedy still prioritises raw new pairings.
BASE_CONNECTION_VALUE = 1000
WEIGHT_DEPT = 100
WEIGHT_ROLE = 10
WEIGHT_COUNTRY = 1


def connection_value(a_attrs: dict, b_attrs: dict) -> int:
    score = BASE_CONNECTION_VALUE
    if a_attrs["department"] != b_attrs["department"]:
        score += WEIGHT_DEPT
    if a_attrs["role"] != b_attrs["role"]:
        score += WEIGHT_ROLE
    if a_attrs["country"] != b_attrs["country"]:
        score += WEIGHT_COUNTRY
    return score


def load_people(path: str) -> tuple[list[str], dict[str, dict]]:
    """Read a CSV with columns: name, department, role, country."""
    people: list[str] = []
    attrs: dict[str, dict] = {}
    with open(path, "r", newline="") as f:
        reader = csv.DictReader(f)
        for row in reader:
            name = (row.get("name") or "").strip()
            if not name:
                continue
            attrs[name] = {
                "department": (row.get("department") or "").strip(),
                "role": (row.get("role") or "").strip(),
                "country": (row.get("country") or "").strip(),
            }
            people.append(name)
    return people, attrs


def _average_pair_value(people: list[str], attrs: dict[str, dict]) -> float:
    if len(people) < 2:
        return float(BASE_CONNECTION_VALUE)
    total = 0
    count = 0
    for a, b in combinations(people, 2):
        total += connection_value(attrs[a], attrs[b])
        count += 1
    return total / count


def create_sittings(
    people: list[str],
    attrs: dict[str, dict],
    tables: int,
    seats: int,
    rounds: int,
) -> tuple[list[list[set]], int, int]:
    if not people or tables <= 0 or seats <= 0 or rounds <= 0:
        return [], 0, 0
    if len(people) > tables * seats:
        raise ValueError(f"Not enough seats for all people: {len(people)} > {tables * seats}")

    tables = min(tables, len(people) // seats + (1 if len(people) % seats > 0 else 0))
    avg_pair_value = _average_pair_value(people, attrs)

    connections: dict[str, set] = {person: set() for person in people}
    seatings: list[list[set]] = [[set() for _ in range(tables)] for _ in range(rounds)]

    for round_idx in range(rounds):
        random.shuffle(people)
        for person in people:
            best_score = -math.inf
            best_table = -1

            for table in range(tables):
                seated = seatings[round_idx][table]
                remaining_seats = seats - len(seated)
                if remaining_seats <= 0:
                    continue

                # Immediate value: diversity-weighted sum over NEW connections only.
                immediate = 0
                for other in seated:
                    if other not in connections[person]:
                        immediate += connection_value(attrs[person], attrs[other])

                # Future potential: empty seats this person will still share, priced at
                # the expected value of a random future connection.
                future = (remaining_seats - 1) * avg_pair_value

                score = immediate + future
                if score > best_score:
                    best_score = score
                    best_table = table

            if best_table == -1:
                raise ValueError("Not enough seats available")

            chosen_table = seatings[round_idx][best_table]
            connections[person].update(chosen_table)
            for p in chosen_table:
                connections[p].add(person)
            chosen_table.add(person)

    unique_connections, weighted_value = score_seatings(seatings, attrs)
    return seatings, unique_connections, weighted_value


def score_seatings(
    seatings: list[list[set]], attrs: dict[str, dict]
) -> tuple[int, int]:
    """Count unique pairs and their total diversity-weighted value."""
    seen: set[tuple[str, str]] = set()
    unique = 0
    weighted = 0
    for round_tables in seatings:
        for table in round_tables:
            for a, b in combinations(sorted(table), 2):
                if (a, b) in seen:
                    continue
                seen.add((a, b))
                unique += 1
                weighted += connection_value(attrs[a], attrs[b])
    return unique, weighted


def improve_by_swaps(
    seatings: list[list[set]], attrs: dict[str, dict]
) -> tuple[list[list[set]], int, int]:
    """Hill-climb: swap pairs of people between tables in the same round whenever
    doing so improves (unique_connections, weighted_value) lexicographically."""
    base_u, base_w = score_seatings(seatings, attrs)
    while True:
        improving = _find_improving_swap(seatings, attrs, base_u, base_w)
        if improving is None:
            return seatings, base_u, base_w
        base_u, base_w = improving


def _find_improving_swap(seatings, attrs, base_u, base_w):
    for round_tables in seatings:
        n_tables = len(round_tables)
        for i in range(n_tables):
            for j in range(i + 1, n_tables):
                for a in list(round_tables[i]):
                    for b in list(round_tables[j]):
                        round_tables[i].discard(a); round_tables[i].add(b)
                        round_tables[j].discard(b); round_tables[j].add(a)
                        u, w = score_seatings(seatings, attrs)
                        if (u, w) > (base_u, base_w):
                            return u, w
                        round_tables[i].discard(b); round_tables[i].add(a)
                        round_tables[j].discard(a); round_tables[j].add(b)
    return None


def write_seatings_to_txt(
    seatings: list[list[set]],
    attrs: dict[str, dict] | None = None,
    filename: str = "seatings.txt",
) -> None:
    with open(filename, "w") as f:
        for round_index, round_tables in enumerate(seatings, start=1):
            f.write(f"Round {round_index}:\n")
            for table_index, table in enumerate(round_tables, start=1):
                f.write(f"Table {table_index}:\n")
                for person in sorted(table):
                    if attrs and person in attrs:
                        a = attrs[person]
                        f.write(
                            f"  - {person} ({a['department']}, {a['role']}, {a['country']})\n"
                        )
                    else:
                        f.write(f"  - {person}\n")
            f.write("\n")


if __name__ == "__main__":
    import os

    people_file = os.getenv("PEOPLE", "people.txt")
    tables = int(os.getenv("TABLES", 15))
    seats = int(os.getenv("SEATS", 6))
    rounds = int(os.getenv("ROUNDS", 2))
    attempts = int(os.getenv("ATTEMPTS", 50))

    people, attrs = load_people(people_file)

    best_unique = -1
    best_weighted = -math.inf
    best_seatings: list[list[set]] = []
    for _ in range(attempts):
        seatings, _, _ = create_sittings(list(people), attrs, tables, seats, rounds)
        seatings, unique_conns, weighted = improve_by_swaps(seatings, attrs)
        # Primary: maximise unique connections. Tiebreaker: diversity-weighted value.
        if (unique_conns, weighted) > (best_unique, best_weighted):
            best_unique = unique_conns
            best_weighted = weighted
            best_seatings = seatings

    write_seatings_to_txt(best_seatings, attrs)
    print(
        f"Best arrangement: {best_unique} unique connections, "
        f"weighted value {best_weighted}."
    )
