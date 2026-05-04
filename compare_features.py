#!/usr/bin/env python3
"""Compare python_features.json vs go_features.json element-wise (tolerance 1e-4)."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


def load_array(path: Path) -> list[float]:
    with path.open(encoding="utf8") as f:
        data = json.load(f)
    if not isinstance(data, list):
        raise SystemExit(f"{path}: expected JSON array, got {type(data).__name__}")
    out: list[float] = []
    for i, x in enumerate(data):
        try:
            out.append(float(x))
        except (TypeError, ValueError) as e:
            raise SystemExit(f"{path}: index {i} not numeric: {x!r}") from e
    return out


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument(
        "--py",
        type=Path,
        default=Path("python_features.json"),
        help="Python JSON array (default: ./python_features.json)",
    )
    ap.add_argument(
        "--go",
        type=Path,
        default=Path("go_features.json"),
        help="Go JSON array (default: ./go_features.json)",
    )
    ap.add_argument(
        "--tol",
        type=float,
        default=1e-4,
        help="absolute difference threshold (default: 0.0001)",
    )
    args = ap.parse_args()

    py_v = load_array(args.py)
    go_v = load_array(args.go)
    if len(py_v) != len(go_v):
        print(
            f"Length mismatch: Python {len(py_v)} vs Go {len(go_v)}",
            file=sys.stderr,
        )
        sys.exit(2)

    n_bad = 0
    for i, (a, b) in enumerate(zip(py_v, go_v)):
        if abs(a - b) > args.tol:
            print(f"Index: {i} - Python: {a} - Go: {b}")
            n_bad += 1

    if n_bad:
        sys.exit(1)
    print(f"OK: {len(py_v)} elements within tol={args.tol}")


if __name__ == "__main__":
    main()
