---
name: use-workshop
description: Run, test, and build the seating arranger inside its Workshop environment. Use when the `workshop` CLI is involved.
user-invocable: false
---

# Use Workshop

Workshop builds a LXD environment with a Go 1.27 toolchain (snap) and `make`. The arranger is a CLI with no services or databases, so nothing is tunneled to the host.

```bash
workshop launch --verbose
workshop shell arrange
```

Run project actions with `workshop run arrange <action> [args]`; `workshop.yaml` defines `go`, `make` and `arrange` (e.g. `workshop run arrange make test`). Anything else goes through `workshop shell`.

The `setup-project` hook runs `make deps` and `go install ./...` on every launch/refresh, so `arrange` is on the PATH inside the environment.

Sheets mode needs `INPUT_SHEET_URL`, `OUTPUT_SHEET_URL` and `GOOGLE_ACCESS_TOKEN`, which only the workflow supplies. Inside the environment, run against a local CSV instead.

## Reset

`workshop remove arrange` then `workshop launch`. Nothing outside the project directory is worth preserving.
