#!/usr/bin/env python3
"""Simple CLI wrapper to extract EMBER v2 feature vector from stdin bytes.

Reads raw file bytes from stdin.buffer, runs the local `PEFeatureExtractor`
and writes a JSON array of feature floats to stdout.
"""
import sys
import json

try:
    # features.py lives in the same directory as this script
    from features import PEFeatureExtractor
except Exception:
    # fallback: try package import
    try:
        from ember.features import PEFeatureExtractor
    except Exception:
        raise


def main():
    data = sys.stdin.buffer.read()
    extractor = PEFeatureExtractor(feature_version=2)
    vec = extractor.feature_vector(data)
    # Convert numpy array to plain Python list
    try:
        out = vec.tolist()
    except Exception:
        out = [float(x) for x in vec]
    json.dump(out, sys.stdout)


if __name__ == '__main__':
    main()
