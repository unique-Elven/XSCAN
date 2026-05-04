package ember

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestImportFQMatchesPythonGoldenABD(t *testing.T) {
	root := findRepoRoot(t)
	exe := filepath.Join(root, "ABD.exe")
	golden := filepath.Join("testdata", "abd_import_fq_golden.txt")
	data, err := os.ReadFile(exe)
	if err != nil {
		t.Skipf("ABD.exe not found: %v", err)
	}
	rawGolden, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("golden file: %v", err)
	}
	want := splitLinesSorted(normalizeLF(strings.TrimSpace(string(rawGolden))))

	p, err := OpenParsedPE(data, false)
	if err != nil || p == nil {
		t.Fatalf("OpenParsedPE: %v", err)
	}
	_, gotFQ := ImportLibrariesAndFQ(p)
	sort.Strings(gotFQ)

	if len(gotFQ) != len(want) {
		t.Fatalf("len got %d want %d", len(gotFQ), len(want))
	}
	for i := range gotFQ {
		if gotFQ[i] != want[i] {
			t.Fatalf("line %d:\n got %q\nwant %q", i, gotFQ[i], want[i])
		}
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "..")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find xscan_client/go.mod from cwd")
		}
		dir = parent
	}
}

func normalizeLF(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}

func splitLinesSorted(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	sort.Strings(lines)
	return lines
}
