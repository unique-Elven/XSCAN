"""Shim: CLI and legacy imports expect ``features.PEFeatureExtractor``."""

from .features_nocert import PEFeatureExtractor

__all__ = ["PEFeatureExtractor"]
