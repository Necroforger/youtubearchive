package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Necroforger/youtubearchive/models"
	"github.com/Necroforger/youtubearchive/util"

	"github.com/jinzhu/gorm"
)

const (
	paramQuery = "q"
	paramLimit = "limit"
	paramPage  = "p"
	paramID    = "id"
)

const (
	tplSearch   = "search"
	tplView     = "view"
	tplHome     = "home"
	tplChannels = "channels"
)

func getSearchParams(r *http.Request) (query string, limit int, page int) {
	query = r.FormValue(paramQuery)

	limit = formValueInt(r, paramLimit, 30)
	page = formValueInt(r, paramPage, 0)
	return
}

func formValueInt(r *http.Request, key string, defaultValue int) int {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}

// queryChannel gives a list of channels
func queryChannels(db *gorm.DB, query string, limit, page int) ([]models.Video, error) {
	var videos []models.Video

	err := db.Group("uploader").Find(&videos).Error

	return videos, err
}

func queryVideos(db *gorm.DB, query string, limit int, page int) ([]models.Video, error) {
	var videos []models.Video

	exquery, tags, keys := util.ParseTags(query)

	raw := "SELECT * FROM videos "
	values := []interface{}{}

	var exset bool
	if exquery != "" {
		raw += "WHERE title LIKE ? "
		values = append(values, "%"+exquery+"%")
		exset = true
	}

	for i, key := range keys {
		if i >= 3 {
			break
		}
		var prefix string
		if i == 0 {
			if !exset {
				prefix = "WHERE"
			} else {
				prefix = "AND"
			}
		} else {
			prefix = "OR"
		}

		v := strings.ToLower(key)
		switch v {
		case "uploader":
			raw += prefix + " uploader LIKE ?"
			values = append(values, "%"+tags[key]+"%")
		case "description":
			raw += prefix + " description LIKE ?"
			values = append(values, "%"+tags[key]+"%")
		}
	}

	raw += " ORDER BY upload_date DESC LIMIT ? OFFSET ? "
	values = append(values, limit, page*limit)

	log.Println(values)

	log.Println(raw)
	err := db.Raw(raw, values...).Scan(&videos).Error

	return videos, err
}

func getVideosByVideoID(db *gorm.DB, ID string, limit, page int) ([]models.Video, error) {
	var videos []models.Video

	err := db.Where("video_id = ?", ID).
		Limit(limit).
		Offset(page * limit).
		Find(&videos).
		Error

	return videos, err
}

// HandleSearch handles searches
func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var (
		query, limit, page = getSearchParams(r)
	)

	videos, err := queryVideos(s.DB, query, limit, page)

	s.ExecuteTemplate(w, r, tplSearch, map[string]interface{}{
		"pages":  []int{page - 1, page, page + 1},
		"query":  query,
		"title":  query,
		"limit":  limit,
		"videos": videos,
		"err":    err,
	})
}

// HandleView views a video
func (s *Server) HandleView(w http.ResponseWriter, r *http.Request) {
	var (
		id   = r.FormValue(paramID)
		page = formValueInt(r, paramPage, 0)
	)

	videos, err := getVideosByVideoID(s.DB, id, 100, page)

	s.ExecuteTemplate(w, r, tplView, map[string]interface{}{
		"videos": videos,
		"err":    err,
	})
}

// HandleHome ...
func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	s.ExecuteTemplate(w, r, tplHome, map[string]interface{}{
		"title": "home",
	})
}

// HandleChannels ...
func (s *Server) HandleChannels(w http.ResponseWriter, r *http.Request) {
	var (
		query, limit, page = getSearchParams(r)
	)

	videos, err := queryChannels(s.DB, query, limit, page)

	s.ExecuteTemplate(w, r, tplChannels, map[string]interface{}{
		"channels": videos,
		"query":    query,
		"pages":    []int{page - 1, page, page + 1},
		"err":      err,
	})
}
