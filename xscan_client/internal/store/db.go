package store

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// OpenDatabase creates/opens xscan_data.db next to the executable and runs AutoMigrate.
func OpenDatabase(exeDir string) (*gorm.DB, error) {
	dbPath := filepath.Join(exeDir, "xscan_data.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite %s: %w", dbPath, err)
	}
	if err := db.AutoMigrate(&Config{}, &ScanHistory{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	if err := ensureDefaultConfig(db, exeDir); err != nil {
		return nil, err
	}
	return db, nil
}

func ensureDefaultConfig(db *gorm.DB, exeDir string) error {
	var c Config
	err := db.First(&c, ConfigSingletonID).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("load config: %w", err)
	}

	c = Config{
		ID:            ConfigSingletonID,
		ModelStrategy: "auto",
		SoundEnabled:  false,
		Language:      "zh",
		Ember2018Path: filepath.Join(exeDir, "ember_model_2018.txt"),
		Ember2024Path: filepath.Join(exeDir, "EMBER2024_all.model"),
		Threshold2018: 0.65,
		Threshold2024: 0.85,
	}
	return db.Create(&c).Error
}
