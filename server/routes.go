package server

import (
	"fmt"
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

	limit = formValueInt(r, paramLimit, 0)
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

	// err := db.Group("uploader").Find(&videos).Error
	err := db.
		Group("uploader").
		Where("uploader LIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(page * limit).
		Find(&videos).Error

	return videos, err
}

func countChannels(db *gorm.DB, query string) (int, error) {
	var count int
	err := db.
		Model(&models.Video{}).
		Group("uploader").
		Where("uploader LIKE ?", "%"+query+"%").
		Count(&count).Error

	return count, err
}

func queryVideos(db *gorm.DB, query string, limit int, page int) ([]models.Video, error) {
	var videos []models.Video

	raw, values := buildVideosQuery(query, limit, page)
	err := db.Raw(raw, values...).Scan(&videos).Error

	return videos, err
}

func countVideos(db *gorm.DB, query string, limit, page int) (int, error) {
	var count int
	raw, values := buildVideosQuery(query, limit, page)
	err := db.Raw("SELECT count(*) FROM ( "+raw+" )", values...).Row().Scan(&count)

	return count, err
}

func buildVideosQuery(query string, limit, page int) (string, []interface{}) {
	exquery, tags, keys := util.ParseTags(query)

	raw := "SELECT * FROM videos"
	values := []interface{}{}

	var exset bool
	if exquery != "" {
		raw += " WHERE title LIKE ?"
		values = append(values, "%"+exquery+"%")
		exset = true
	}

	for i, key := range keys {
		var prefix string
		if i == 0 {
			if !exset {
				prefix = " WHERE"
			} else {
				prefix = " AND"
			}
		} else {
			prefix = " AND"
		}

		v := strings.ToLower(key)
		switch v {
		case "uploader":
			raw += prefix + " uploader LIKE ?"
			values = append(values, "%"+tags[key]+"%")
		case "description":
			raw += prefix + " description LIKE ?"
			values = append(values, "%"+tags[key]+"%")
		case "title":
			parts := strings.Split(tags[key], ",")
			for i := 0; i < len(parts); i++ {
				if i == 0 {
					raw += prefix + " title LIKE ?"
				} else {
					raw += " AND title LIKE ?"
				}
				values = append(values, "%"+parts[i]+"%")
			}
		}
	}

	raw += " ORDER BY upload_date DESC"
	if limit >= 0 {
		raw += " LIMIT ?"
		values = append(values, limit)
		if page >= 0 {
			raw += " OFFSET ?"
			values = append(values, page*limit)
		}
	}

	fmt.Println(raw)
	return raw, values
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
	if limit == 0 {
		limit = 100
	}

	var reterr error

	videos, err := queryVideos(s.DB, query, limit, page)
	if err != nil {
		s.Log("error querying videos: ", err)
		reterr = err
	}

	total, err := countVideos(s.DB, query, -1, -1)
	if err != nil {
		s.Log("error counting videos: ", err)
		reterr = err
	}
	total = int(float64(total)/float64(limit) + 0.9)

	s.ExecuteTemplate(w, r, tplSearch, map[string]interface{}{
		"pages":     []int{page - 1, page, page + 1},
		"query":     query,
		"title":     query,
		"limit":     limit,
		"videos":    videos,
		"err":       reterr,
		"paginator": NewPaginator(page, total, 31, query, limit, "/search"),
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
	if limit == 0 {
		limit = 100
	}

	var reterror error

	videos, err := queryChannels(s.DB, query, limit, page)
	if err != nil {
		reterror = err
	}
	total, err := countChannels(s.DB, query)
	if err != nil {
		reterror = err
	}
	total = int(float64(total)/float64(limit) + 0.9)

	s.ExecuteTemplate(w, r, tplChannels, map[string]interface{}{
		"channels":  videos,
		"query":     query,
		"pages":     []int{page - 1, page, page + 1},
		"err":       reterror,
		"paginator": NewPaginator(page, total, 31, query, limit, "/channels"),
	})
}
