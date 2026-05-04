#!/usr/bin/env python3
"""Dump go_server/ember_cert feature vector to JSON (aligned with Go ExtractFeatures layout)."""

from __future__ import annotations

import argparse
import importlib.util
import json
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent


def _load_pe_feature_extractor():
    """Load PEFeatureExtractor without importing ember_cert package __init__ (no lightgbm)."""
    feat_py = REPO_ROOT / "go_server" / "ember_cert" / "features.py"
    spec = importlib.util.spec_from_file_location("ember_cert_features", feat_py)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load {feat_py}")
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod.PEFeatureExtractor


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--path", type=Path, required=True, help="PE file path")
    ap.add_argument(
        "-o",
        "--out",
        type=Path,
        default=Path("python_features.json"),
        help="output JSON array (default: ./python_features.json)",
    )
    args = ap.parse_args()

    PEFeatureExtractor = _load_pe_feature_extractor()
    data = args.path.read_bytes()
    vec = PEFeatureExtractor().feature_vector(data)
    arr = vec.tolist() if hasattr(vec, "tolist") else [float(x) for x in vec]
    args.out.write_text(json.dumps(arr), encoding="utf8")
    print(f"wrote {len(arr)} floats to {args.out}", file=sys.stderr)


if __name__ == "__main__":
    main()
