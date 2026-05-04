package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"xscan_client/internal/ember"
)

func main() {
	path := flag.String("path", "", "path to PE file")
	includeCert := flag.Bool("cert", false, "parse security directory / authenticode block like Go ExtractFeatures(..., true)")
	out := flag.String("o", "go_features.json", "output JSON array path")
	flag.Parse()
	if *path == "" {
		log.Fatal("missing -path")
	}
	feat, err := ember.ExtractFeatures(*path, *includeCert)
	if err != nil {
		log.Fatal(err)
	}
	if err := ember.WriteFeaturesJSON(*out, feat); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d floats to %s\n", len(feat), *out)
}
