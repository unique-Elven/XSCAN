package ember

import (
	"encoding/json"
	"math"
	"regexp"
	"runtime"
	"sync"

	_ "embed"
)

//go:embed string_regexes.json
var stringRegexesJSON []byte

var stringRegexTable []*regexp.Regexp

func init() {
	var dump struct {
		N     int `json:"n"`
		Items []struct {
			K string `json:"k"`
			P string `json:"p"`
			I bool   `json:"i"`
		} `json:"items"`
	}
	if err := json.Unmarshal(stringRegexesJSON, &dump); err != nil {
		panic("ember: string_regexes.json: " + err.Error())
	}
	for _, it := range dump.Items {
		pat := it.P
		if it.I {
			pat = "(?i)" + pat
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			panic("ember: compile regex " + it.K + ": " + err.Error())
		}
		stringRegexTable = append(stringRegexTable, re)
	}
	if len(stringRegexTable) != dump.N {
		panic("ember: regex table length mismatch")
	}
}

func extractPrintableRuns(data []byte) [][]byte {
	var runs [][]byte
	i := 0
	for i < len(data) {
		start := i
		for i < len(data) && data[i] >= 0x20 && data[i] <= 0x7f {
			i++
		}
		if i-start >= 5 {
			runs = append(runs, data[start:i])
		}
		if i < len(data) {
			i++
		}
	}
	return runs
}

func entropyPrintable(counts *[96]float64, total int) float64 {
	if total <= 0 {
		return 0
	}
	var h float64
	inv := 1.0 / float64(total)
	for _, c := range counts {
		if c <= 0 {
			continue
		}
		p := c * inv
		h -= p * math.Log2(p)
	}
	return h
}

func countRegexMatchesParallel(runs [][]byte, regs []*regexp.Regexp) []float32 {
	out := make([]float32, len(regs))
	if len(runs) == 0 {
		return out
	}
	n := runtime.GOMAXPROCS(0)
	if n < 1 {
		n = 1
	}
	if n > len(runs) {
		n = len(runs)
	}
	step := (len(runs) + n - 1) / n
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < len(runs); i += step {
		j := i + step
		if j > len(runs) {
			j = len(runs)
		}
		chunk := runs[i:j]
		wg.Add(1)
		go func(seg [][]byte) {
			defer wg.Done()
			local := make([]float32, len(regs))
			for _, run := range seg {
				for ri, re := range regs {
					if re.Match(run) {
						local[ri]++
					}
				}
			}
			mu.Lock()
			for k := range out {
				out[k] += local[k]
			}
			mu.Unlock()
		}(chunk)
	}
	wg.Wait()
	return out
}

// ProcessStringExtractor writes StringExtractor.dim floats (ember_cert).
func ProcessStringExtractor(data []byte, dst []float32) {
	if len(dst) < DimStrings {
		return
	}
	clear(dst[:DimStrings])

	runs := extractPrintableRuns(data)
	nRuns := len(runs)
	dst[0] = float32(nRuns)

	var hist [96]float64
	totalPrint := 0
	sumLen := 0
	for _, run := range runs {
		sumLen += len(run)
		for _, b := range run {
			idx := int(b - 0x20)
			if idx >= 0 && idx < 96 {
				hist[idx]++
				totalPrint++
			}
		}
	}

	var avLen float32
	if nRuns > 0 {
		avLen = float32(float64(sumLen) / float64(nRuns))
	}
	dst[1] = avLen
	dst[2] = float32(totalPrint)

	div := float64(totalPrint)
	if div <= 0 {
		div = 1
	}
	for i := 0; i < 96; i++ {
		dst[3+i] = float32(hist[i] / div)
	}
	dst[99] = float32(entropyPrintable(&hist, totalPrint))

	regCounts := countRegexMatchesParallel(runs, stringRegexTable)
	copy(dst[100:], regCounts)
}
