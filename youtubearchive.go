package youtubearchive

import (
	"github.com/Necroforger/youtubearchive/models"
	"github.com/jinzhu/gorm"
)

// InitDB initializes the gorm DB with models
func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Video{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Tag{}).Error
	if err != nil {
		return err
	}

	// Create a table for storing channel metadata as JSON
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channel_metadata (
			ID               INTEGER PRIMARY KEY,
			created          TEXT,
			uploader_url     TEXT,
			json             TEXT
		);
	`).Error
	if err != nil {
		return err
	}

	return nil
}
