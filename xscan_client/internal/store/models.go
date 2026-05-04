package store

import (
	"time"
)

// ConfigSingletonID is the fixed primary key for the single settings row.
const ConfigSingletonID = uint(1)

// Config persists lightweight UI / scanner preferences (one row).
type Config struct {
	ID uint `gorm:"primaryKey"`

	ModelStrategy string `gorm:"size:32;default:auto"` // e.g. auto-routed Ember2018/2024
	SoundEnabled  bool   `gorm:"default:false"`
	Language      string `gorm:"size:8;default:zh"`

	Ember2018Path string  `gorm:"size:2048"`
	Ember2024Path string  `gorm:"size:2048"`
	Threshold2018 float64 `gorm:"default:0.65"`
	Threshold2024 float64 `gorm:"default:0.85"`
}

// ScanHistory is one finished scan record for the history view.
type ScanHistory struct {
	ID        uint `gorm:"primaryKey"`
	ScannedAt time.Time `gorm:"not null;index;column:scanned_at"`
	FilePath  string    `gorm:"size:4096;not null;column:file_path"`
	FileHash  string    `gorm:"size:128;column:file_hash"`
	Verdict   string    `gorm:"size:32;not null"` // malicious | safe
	FileSize  int64     `gorm:"default:0;column:file_size"`
	Engine    string    `gorm:"size:16;column:engine"` // e.g. 2018 | 2024
	Score     float32   `gorm:"column:score"`
}

func (ScanHistory) TableName() string {
	return "scan_histories"
}
