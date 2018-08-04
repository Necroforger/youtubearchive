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

	return nil
}
