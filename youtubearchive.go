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
			updated          TEXT,
			uploader_url     TEXT,
			json             TEXT
		);

		-- A view of all channels most recent metadata archive
		CREATE VIEW IF NOT EXISTS recent_channel_metadata AS
			select * from channel_metadata where ID in (
			select 
				(select ID from channel_metadata where a.uploader_url = uploader_url order by created desc) as ID
			from
				(select uploader_url from channel_metadata group by uploader_url) a)
	`).Error
	if err != nil {
		return err
	}

	return nil
}
