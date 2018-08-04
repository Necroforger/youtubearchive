package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Video ...
type Video struct {
	gorm.Model

	LastScanned time.Time

	Title       string `gorm:"index"`
	Views       int
	Likes       int
	Thumbnail   string
	Duration    int
	Description string
	Uploader    string
	UploaderURL string
	UploaderID  string
	UploadDate  string `gorm:"index"`
	VideoID     string
	Tags        []Tag
	WebpageURL  string
}

// Tag is a tag
type Tag struct {
	gorm.Model
	Value string
}

// MakeTags makes a slice of strings into a slice of tags
func MakeTags(data []string) []Tag {
	retval := make([]Tag, len(data))
	for i := 0; i < len(data); i++ {
		retval[i] = Tag{
			Value: data[i],
		}
	}
	return retval
}
