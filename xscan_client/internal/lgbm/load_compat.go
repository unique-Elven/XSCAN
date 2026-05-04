package lgbm

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/dmitryikh/leaves"
)

// patchLightGBMTextVersion rewrites unsupported text-model version lines (e.g. v4/v5)
// to v3 so github.com/dmitryikh/leaves can parse; v2/v3 are left unchanged.
// Binary models are returned unchanged.
func patchLightGBMTextVersion(data []byte) []byte {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	trimLeft := bytes.TrimLeft(data, " \t\r\n")
	if !bytes.HasPrefix(trimLeft, []byte("tree")) {
		return data
	}
	br := bufio.NewReader(bytes.NewReader(data))
	var out bytes.Buffer
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return data
		}
		if len(bytes.TrimSpace(line)) == 0 {
			out.Write(line)
			rest, _ := io.ReadAll(br)
			out.Write(rest)
			break
		}
		line = rewriteVersionLine(line)
		out.Write(line)
		if err == io.EOF {
			break
		}
	}
	return out.Bytes()
}

func rewriteVersionLine(line []byte) []byte {
	s := strings.TrimSpace(string(line))
	if !strings.HasPrefix(s, "version=") {
		return line
	}
	ver := strings.TrimPrefix(s, "version=")
	if ver == "v2" || ver == "v3" {
		return line
	}
	var nl string
	switch {
	case strings.HasSuffix(string(line), "\r\n"):
		nl = "\r\n"
	case strings.HasSuffix(string(line), "\n"):
		nl = "\n"
	default:
		nl = "\n"
	}
	return append([]byte("version=v3"), nl...)
}

// ensembleFromPath loads a LightGBM **text** dump (.txt / text-based .model).
// Unsupported header versions are normalized before parsing.
func ensembleFromPath(path string, loadTransformation bool) (*leaves.Ensemble, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	patched := patchLightGBMTextVersion(raw)
	r := bufio.NewReader(bytes.NewReader(patched))
	return leaves.LGEnsembleFromReader(r, loadTransformation)
}
