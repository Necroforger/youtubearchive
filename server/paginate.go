package server

import (
	"net/url"
	"strconv"
)

// Paginator stores pagination information for templates
type Paginator struct {
	Query   string
	Limit   string
	Base    string
	Current string

	Beginning string
	Middle    []string
	Ending    string
}

// GetURL gets the url for your given page
func (p Paginator) GetURL(page string) string {
	return p.Base + "?q=" + url.QueryEscape(p.Query) + "&p=" + page + "&limit=" + p.Limit
}

// NewPaginator creates a new paginator
func NewPaginator(position, total, length int, query string, limit int, base string) Paginator {
	start, middle, end := Paginate(position, total, length)
	return Paginator{
		Query:   query,
		Limit:   strconv.Itoa(limit),
		Base:    base,
		Current: strconv.Itoa(position),

		Beginning: start,
		Middle:    middle,
		Ending:    end,
	}
}

// Paginate generates paginated links
func Paginate(position, total, length int) (beginning string, middle []string, ending string) {
	if length <= 0 {
		length = 1
	}

	position -= length / 2

	if position > 0 {
		beginning = "0"
	}
	if position+length < total {
		ending = strconv.Itoa(total - 1)
	}

	if position+length >= total {
		position -= (position + length - total)
	}
	if position < 0 {
		position = 0
	}

	for i := position; i < position+length && i < total; i++ {
		middle = append(middle, strconv.Itoa(i))
	}
	return
}
