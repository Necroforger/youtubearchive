package youtubearchive

import (
	"time"

	"github.com/Necroforger/youtubearchive/youtubedl"
	"github.com/jinzhu/gorm"
)

// Video ...
type Video struct {
	gorm.Model

	LastScanned time.Time

	Title       string
	Views       int
	Likes       int
	Thumbnail   string
	Duration    int
	Description string
	Uploader    string
	UploaderURL string
	UploaderID  string
	UploadDate  string
	VideoID     string
	Tags        []Tag
	WebpageURL  string
}

// VideoFromYoutubedl returns a Video model from a youtube-dl video
func VideoFromYoutubedl(v youtubedl.Video) Video {
	return Video{
		Title:       v.Title,
		Views:       v.ViewCount,
		Likes:       v.LikeCount,
		Thumbnail:   v.Thumbnail,
		Duration:    v.Duration,
		Description: v.Description,
		Uploader:    v.Uploader,
		UploaderURL: v.UploaderURL,
		UploaderID:  v.UploaderID,
		UploadDate:  v.UploadDate,
		VideoID:     v.ID,
		Tags:        MakeTags(v.Tags),
		WebpageURL:  v.WebpageURL,
	}
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
