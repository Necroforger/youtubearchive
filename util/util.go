package util

import (
	"strings"

	"github.com/Necroforger/youtubearchive/models"
	"github.com/Necroforger/youtubearchive/youtubedl"
)

// VideoFromYoutubedl returns a Video model from a youtube-dl video
func VideoFromYoutubedl(v youtubedl.Video) models.Video {
	return models.Video{
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
		Tags:        models.MakeTags(v.Tags),
		WebpageURL:  v.WebpageURL,
	}
}

// ParseTags parses a given text into ':' separated tags
// So that "uploader:name blah" will return
// "blah", {"uploader": "name"}, [uploader]
// words are split by semicolons ";"
func ParseTags(data string) (string, map[string]string, []string) {
	words := strings.Split(data, ";")
	unextracted := []string{}
	extracted := map[string]string{}
	keys := []string{}

	for _, word := range words {
		if strings.Contains(word, ":") {
			parts := strings.SplitN(word, ":", 2)
			if len(parts) > 1 {
				extracted[parts[0]] = parts[1]
				keys = append(keys, parts[0])
			}
		} else {
			unextracted = append(unextracted, strings.TrimSpace(word))
		}
	}

	return strings.Join(unextracted, ";"), extracted, keys
}
