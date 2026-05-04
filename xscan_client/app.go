package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"xscan_client/internal/ember"
	"xscan_client/internal/lgbm"
	"xscan_client/internal/pemeta"
	"xscan_client/internal/store"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx context.Context

	lgbm     lgbm.Engine
	scanPool *lgbm.ScanModelPool
	db       *gorm.DB

	scanMu     sync.RWMutex
	scannerCfg scannerModelPaths
}

type scannerModelPaths struct {
	Path2018 string
	Path2024 string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		scanPool: &lgbm.ScanModelPool{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	exe, err := os.Executable()
	if err != nil {
		println("xscan executable path:", err.Error())
		return
	}
	db, err := store.OpenDatabase(filepath.Dir(exe))
	if err != nil {
		println("xscan_data.db:", err.Error())
		return
	}
	a.db = db
}

// --- Local persistence (SQLite + GORM, pure Go driver) ---

// AppConfigDTO is persisted settings for the Vue client.
type AppConfigDTO struct {
	ModelStrategy string  `json:"modelStrategy"`
	SoundEnabled  bool    `json:"soundEnabled"`
	Language      string  `json:"language"`
	Path2018      string  `json:"path2018"`
	Path2024      string  `json:"path2024"`
	Threshold2018 float64 `json:"threshold2018"`
	Threshold2024 float64 `json:"threshold2024"`
}

// HistorySaveEntry is one row appended by the frontend after a scan.
type HistorySaveEntry struct {
	ScannedAt int64   `json:"scannedAt"`
	FilePath  string  `json:"filePath"`
	FileHash  string  `json:"fileHash"`
	Verdict   string  `json:"verdict"`
	FileSize  int64   `json:"fileSize"`
	Engine    string  `json:"engine"`
	Score     float64 `json:"score"`
}

// ScanHistoryRow is returned to the history page.
type ScanHistoryRow struct {
	ID        uint    `json:"id"`
	ScannedAt int64   `json:"scannedAt"`
	FilePath  string  `json:"filePath"`
	FileHash  string  `json:"fileHash"`
	Verdict   string  `json:"verdict"`
	FileSize  int64   `json:"fileSize"`
	Engine    string  `json:"engine"`
	Score     float64 `json:"score"`
}

// GetAppConfig loads the singleton Config row (defaults live next to the executable).
func (a *App) GetAppConfig() (AppConfigDTO, error) {
	var dto AppConfigDTO
	if a.db == nil {
		return dto, fmt.Errorf("database unavailable")
	}
	var c store.Config
	if err := a.db.First(&c, store.ConfigSingletonID).Error; err != nil {
		return dto, err
	}
	dto.ModelStrategy = c.ModelStrategy
	dto.SoundEnabled = c.SoundEnabled
	dto.Language = c.Language
	dto.Path2018 = c.Ember2018Path
	dto.Path2024 = c.Ember2024Path
	dto.Threshold2018 = c.Threshold2018
	dto.Threshold2024 = c.Threshold2024
	return dto, nil
}

// SaveAppConfig updates persisted UI / scanner settings.
func (a *App) SaveAppConfig(cfg AppConfigDTO) error {
	if a.db == nil {
		return fmt.Errorf("database unavailable")
	}
	var c store.Config
	if err := a.db.First(&c, store.ConfigSingletonID).Error; err != nil {
		return err
	}
	ms := strings.TrimSpace(cfg.ModelStrategy)
	if ms == "" {
		ms = "auto"
	}
	lang := strings.TrimSpace(cfg.Language)
	if lang == "" {
		lang = "zh"
	}
	c.ModelStrategy = ms
	c.SoundEnabled = cfg.SoundEnabled
	c.Language = lang
	c.Ember2018Path = strings.TrimSpace(cfg.Path2018)
	c.Ember2024Path = strings.TrimSpace(cfg.Path2024)
	c.Threshold2018 = cfg.Threshold2018
	c.Threshold2024 = cfg.Threshold2024
	return a.db.Save(&c).Error
}

func normalizeHistoryVerdict(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	if s == "malicious" || s == "malware" {
		return "malicious"
	}
	return "safe"
}

// SaveHistory inserts scan rows; empty fileHash is filled from SHA-256 when the path is readable.
func (a *App) SaveHistory(entries []HistorySaveEntry) error {
	if a.db == nil {
		return fmt.Errorf("database unavailable")
	}
	if len(entries) == 0 {
		return nil
	}
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, e := range entries {
			fp := strings.TrimSpace(e.FilePath)
			if fp == "" {
				continue
			}
			hash := strings.TrimSpace(e.FileHash)
			if hash == "" {
				if h, err := store.FileSHA256Hex(fp); err == nil {
					hash = h
				}
			}
			ts := time.UnixMilli(e.ScannedAt)
			if e.ScannedAt <= 0 {
				ts = time.Now()
			}
			eng := strings.TrimSpace(e.Engine)
			if eng == "" {
				eng = "2018"
			}
			row := store.ScanHistory{
				ScannedAt: ts,
				FilePath:  fp,
				FileHash:  hash,
				Verdict:   normalizeHistoryVerdict(e.Verdict),
				FileSize:  e.FileSize,
				Engine:    eng,
				Score:     float32(e.Score),
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetHistory returns recent rows ordered by scan time desc (newest first).
func (a *App) GetHistory(limit int) ([]ScanHistoryRow, error) {
	if a.db == nil {
		return nil, fmt.Errorf("database unavailable")
	}
	if limit <= 0 {
		limit = 500
	}
	var rows []store.ScanHistory
	if err := a.db.Order("scanned_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]ScanHistoryRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, ScanHistoryRow{
			ID:        r.ID,
			ScannedAt: r.ScannedAt.UnixMilli(),
			FilePath:  r.FilePath,
			FileHash:  r.FileHash,
			Verdict:   r.Verdict,
			FileSize:  r.FileSize,
			Engine:    r.Engine,
			Score:     float64(r.Score),
		})
	}
	return out, nil
}

// ClearScanHistory removes all rows from the scan history table.
func (a *App) ClearScanHistory() error {
	if a.db == nil {
		return fmt.Errorf("database unavailable")
	}
	return a.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&store.ScanHistory{}).Error
}

// DeleteHistory removes scan history rows by primary key.
func (a *App) DeleteHistory(ids []uint) error {
	if a.db == nil {
		return fmt.Errorf("database unavailable")
	}
	if len(ids) == 0 {
		return nil
	}
	return a.db.Where("id IN ?", ids).Delete(&store.ScanHistory{}).Error
}

func engineTagFromModelUsed(modelUsed string) string {
	if strings.Contains(modelUsed, "2024") {
		return "2024"
	}
	return "2018"
}

func (a *App) scanThresholds() (th18, th24 float64) {
	th18, th24 = 0.65, 0.85
	if a.db == nil {
		return
	}
	var c store.Config
	if err := a.db.First(&c, store.ConfigSingletonID).Error; err != nil {
		return
	}
	return c.Threshold2018, c.Threshold2024
}

func maliciousPathsFromResults(out []FileScoreResult, th18, th24 float64) []string {
	var paths []string
	for _, r := range out {
		if r.ErrMsg != "" {
			continue
		}
		eng := engineTagFromModelUsed(r.ModelUsed)
		th := th18
		if eng == "2024" {
			th = th24
		}
		if r.Score >= th {
			paths = append(paths, r.Path)
		}
	}
	return paths
}

// persistScanHistoryFromResults writes successful scan rows using thresholds from persisted Config.
func (a *App) persistScanHistoryFromResults(results []FileScoreResult) error {
	if a.db == nil || len(results) == 0 {
		return nil
	}
	var cfg store.Config
	if err := a.db.First(&cfg, store.ConfigSingletonID).Error; err != nil {
		return fmt.Errorf("load config for history: %w", err)
	}
	for _, r := range results {
		if r.ErrMsg != "" {
			continue
		}
		eng := engineTagFromModelUsed(r.ModelUsed)
		th := cfg.Threshold2018
		if eng == "2024" {
			th = cfg.Threshold2024
		}
		verdict := "safe"
		if r.Score >= th {
			verdict = "malicious"
		}
		hash := ""
		if h, err := store.FileSHA256Hex(r.Path); err == nil {
			hash = h
		}
		row := store.ScanHistory{
			ScannedAt: time.Now(),
			FilePath:  r.Path,
			FileHash:  hash,
			Verdict:   verdict,
			FileSize:  r.Size,
			Engine:    eng,
			Score:     float32(r.Score),
		}
		if err := a.db.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

// LightGBMModelStatus describes the currently loaded classifier model.
type LightGBMModelStatus struct {
	Loaded    bool   `json:"loaded"`
	ModelPath string `json:"modelPath"`
	NFeatures int    `json:"nFeatures"`
}

// LoadLightGBMModel loads a LightGBM model from disk (LightGBM **text** dump: .txt or text .model).
func (a *App) LoadLightGBMModel(modelPath string) error {
	return a.lgbm.EnsureLoaded(modelPath)
}

func isScanTargetExt(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".exe", ".dll", ".exx":
		return true
	default:
		return false
	}
}

// LightGBMModelStatus returns load state for the active model.
func (a *App) LightGBMModelStatus() LightGBMModelStatus {
	p := a.lgbm.LoadedPath()
	return LightGBMModelStatus{
		Loaded:    p != "",
		ModelPath: p,
		NFeatures: a.lgbm.NFeatures(),
	}
}

// UnloadLightGBMModel clears the in-memory model.
func (a *App) UnloadLightGBMModel() {
	a.lgbm.Unload()
}

// PredictLightGBM runs inference on a feature vector (must match model input size).
func (a *App) PredictLightGBM(features []float32) (float64, error) {
	return a.lgbm.Predict(features)
}

// PredictFileLightGBM extracts Ember partial features (general+histogram) then predicts.
// Returns an error if feature dimension does not match the loaded model.
func (a *App) PredictFileLightGBM(filePath string, includeCert bool) (float64, error) {
	feat, err := ember.ExtractFeatures(filePath, includeCert)
	if err != nil {
		return 0, err
	}
	return a.lgbm.Predict(feat)
}

// ScannerModelPaths is pushed from the frontend (settings) so batch scans need only file paths.
type ScannerModelPaths struct {
	Path2018 string `json:"path2018"`
	Path2024 string `json:"path2024"`
}

// ConfigureScanner validates paths, preloads both LightGBM ensembles into memory, then stores paths for ScanFiles.
func (a *App) ConfigureScanner(paths ScannerModelPaths) error {
	p18 := strings.TrimSpace(paths.Path2018)
	p24 := strings.TrimSpace(paths.Path2024)
	if p18 == "" || p24 == "" {
		return fmt.Errorf("both path2018 and path2024 are required")
	}
	if err := a.scanPool.Reload(p18, p24); err != nil {
		return err
	}
	a.scanMu.Lock()
	a.scannerCfg.Path2018 = p18
	a.scannerCfg.Path2024 = p24
	a.scanMu.Unlock()
	return nil
}

// FileScoreResult is one file outcome from a batch scan.
type FileScoreResult struct {
	Path      string  `json:"path"`
	Score     float64 `json:"score"`
	Size      int64   `json:"size"`
	ModelUsed string  `json:"modelUsed,omitempty"`
	ErrMsg    string  `json:"errMsg,omitempty"`
}

// ScanFiles scores each path using saferwall/pe–based signing detection, the matching Ember feature
// extractor, and the preloaded pool model (2018 vs 2024). Raw feature length is aligned to each
// ensemble’s NFeatures() via padding or truncation before PredictSingle.
func (a *App) ScanFiles(filePaths []string) ([]FileScoreResult, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no files to scan")
	}
	a.scanMu.RLock()
	cfg := a.scannerCfg
	a.scanMu.RUnlock()
	if cfg.Path2018 == "" || cfg.Path2024 == "" {
		return nil, fmt.Errorf("scanner not configured: call ConfigureScanner with both models")
	}
	if err := a.scanPool.Ensure(cfg.Path2018, cfg.Path2024); err != nil {
		return nil, err
	}
	defer debug.FreeOSMemory()

	type job struct {
		idx       int
		path      string
		useSigned bool
	}
	unsigned := make([]job, 0)
	signed := make([]job, 0)
	for i, p := range filePaths {
		useSigned := pemeta.ShouldUseSignedPipelineFromPath(p)
		j := job{idx: i, path: p, useSigned: useSigned}
		if useSigned {
			signed = append(signed, j)
		} else {
			unsigned = append(unsigned, j)
		}
	}

	out := make([]FileScoreResult, len(filePaths))

	runBatch := func(jobs []job, includeCert bool, modelLabel string, predict func([]float32) (float64, error)) error {
		if len(jobs) == 0 {
			return nil
		}
		for _, j := range jobs {
			p := j.path
			var sz int64
			if fi, err := os.Stat(p); err == nil {
				sz = fi.Size()
			}
			feat, err := ember.ExtractFeatures(p, includeCert)
			if err != nil {
				out[j.idx] = FileScoreResult{Path: p, Size: sz, ModelUsed: modelLabel, ErrMsg: err.Error()}
				continue
			}
			score, err := predict(feat)
			if err != nil {
				out[j.idx] = FileScoreResult{Path: p, Size: sz, ModelUsed: modelLabel, ErrMsg: err.Error()}
				continue
			}
			out[j.idx] = FileScoreResult{Path: p, Score: score, Size: sz, ModelUsed: modelLabel}
		}
		return nil
	}

	if err := runBatch(unsigned, false, "Ember2018", a.scanPool.PredictUnsigned); err != nil {
		return nil, err
	}
	if err := runBatch(signed, true, "Ember2024", a.scanPool.PredictSigned); err != nil {
		return nil, err
	}
	if err := a.persistScanHistoryFromResults(out); err != nil {
		return nil, fmt.Errorf("persist scan history: %w", err)
	}
	th18, th24 := a.scanThresholds()
	malicious := maliciousPathsFromResults(out, th18, th24)
	runtime.EventsEmit(a.ctx, "scan_completed", map[string]interface{}{
		"maliciousPaths": malicious,
	})
	return out, nil
}

// ScanFilesLightGBM loads modelPath (if needed), then scores each file path on disk.
// includeCert is passed to the PE feature extractor (security directory when true).
func (a *App) ScanFilesLightGBM(modelPath string, filePaths []string, includeCert bool) ([]FileScoreResult, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no files to scan")
	}
	if err := a.lgbm.EnsureLoaded(modelPath); err != nil {
		return nil, err
	}
	out := make([]FileScoreResult, 0, len(filePaths))
	for _, p := range filePaths {
		var sz int64
		if fi, err := os.Stat(p); err == nil {
			sz = fi.Size()
		}
		feat, err := ember.ExtractFeatures(p, includeCert)
		if err != nil {
			out = append(out, FileScoreResult{Path: p, Size: sz, ErrMsg: err.Error()})
			continue
		}
		score, err := a.lgbm.Predict(feat)
		if err != nil {
			out = append(out, FileScoreResult{Path: p, Size: sz, ErrMsg: err.Error()})
			continue
		}
		out = append(out, FileScoreResult{Path: p, Score: score, Size: sz})
	}
	return out, nil
}

// QuarantineFiles renames each path by appending ".exx" (e.g. a.exe → a.exe.exx).
// Returns how many renames succeeded.
func (a *App) QuarantineFiles(filePaths []string) (int, error) {
	if len(filePaths) == 0 {
		return 0, nil
	}
	n := 0
	var firstErr error
	for _, p := range filePaths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		dest := p + ".exx"
		if err := os.Rename(p, dest); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		n++
	}
	if n == 0 && firstErr != nil {
		return 0, firstErr
	}
	return n, nil
}

// PickLightGBMModelFile opens a native file dialog to choose a LightGBM text model.
func (a *App) PickLightGBMModelFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 LightGBM 模型",
		Filters: []runtime.FileFilter{
			{DisplayName: "LightGBM model (*.txt;*.model)", Pattern: "*.txt;*.model"},
			{DisplayName: "All files", Pattern: "*.*"},
		},
	})
}

// PickScanFiles opens a native multi-select file dialog for scanning.
func (a *App) PickScanFiles() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择要扫描的文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "PE / executables (*.exe;*.dll;*.exx)", Pattern: "*.exe;*.dll;*.exx"},
			{DisplayName: "All files", Pattern: "*.*"},
		},
	})
}

// PickScanDirectory opens a native folder dialog and returns all regular files recursively (simple depth scan).
func (a *App) PickScanDirectory() ([]string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择要扫描的目录",
	})
	if err != nil || dir == "" {
		return nil, err
	}
	var paths []string
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		if !isScanTargetExt(path) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
