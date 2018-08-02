package youtubedl

// FlatVideo is returned by --flat-playlist
type FlatVideo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	IeKey string `json:"ie_key"`
	Type  string `json:"_type"`
	ID    string `json:"id"`
}

// Video is a Youtube video
type Video struct {
	UploaderURL  string `json:"uploader_url"`
	Extractor    string `json:"extractor"`
	Abr          int    `json:"abr"`
	FormatID     string `json:"format_id"`
	DislikeCount int    `json:"dislike_count"`
	DisplayID    string `json:"display_id"`
	// Categories         []string    `json:"categories"`
	Description        string      `json:"description"`
	IsLive             bool        `json:"is_live"`
	Filename           string      `json:"_filename"`
	WebpageURL         string      `json:"webpage_url"`
	ViewCount          int         `json:"view_count"`
	AverageRating      float64     `json:"average_rating"`
	Height             int         `json:"height"`
	Fulltitle          string      `json:"fulltitle"`
	UploadDate         string      `json:"upload_date"`
	Playlist           interface{} `json:"playlist"`
	UploaderID         string      `json:"uploader_id"`
	Vcodec             string      `json:"vcodec"`
	Width              int         `json:"width"`
	Thumbnail          string      `json:"thumbnail"`
	Ext                string      `json:"ext"`
	ID                 string      `json:"id"`
	Tags               []string    `json:"tags"`
	Acodec             string      `json:"acodec"`
	Duration           int         `json:"duration"`
	WebpageURLBasename string      `json:"webpage_url_basename"`
	ExtractorKey       string      `json:"extractor_key"`
	Artist             string      `json:"artist"`
	Format             string      `json:"format"`
	Uploader           string      `json:"uploader"`
	PlaylistIndex      int         `json:"playlist_index"`
	AgeLimit           int         `json:"age_limit"`
	LikeCount          int         `json:"like_count"`
	Fps                int         `json:"fps"`
	Title              string      `json:"title"`
}
