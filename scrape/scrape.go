/*
Package scrape provides some methods for scraping information from youtube

- TODO
[ ] Channel subscriptions
[ ] Channel playlists
*/
package scrape

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
TODO: Classes and other scraping information below

.channel-header-profile-image
.about-description
.about-stats
.yt-lockup-playlist

You can also find playlists by using a CSS selector looking for links that look like
/playlist?list=PLNCRTSKrIMvss_8PSICTJJxUWKUgSu2nU
*/

// Subscription ...
type Subscription struct {
}

// Link stores link information
type Link struct {
	Name string
	URL  string
}

// ChannelInfo contains various channel information and statistics
type ChannelInfo struct {
	Name         string
	ProfileImage string
	Description  string
	HeaderLinks  []Link
	Related      []Link
	AllStats     []string
	Views        int
	Subscribers  int
	Joined       string
}

// Subscriptions ...
func Subscriptions(URL string) Subscription {
	return Subscription{}
}

func filterStr(str string, f func(s rune) bool) (retval string) {
	for _, v := range str {
		if f(v) {
			retval += string(v)
		}
	}
	return
}

// GetChannelInfo scrapes channel information from a channel URL
func GetChannelInfo(URL string) (info ChannelInfo, err error) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return
	}

	info.Name = doc.Find(".qualified-channel-title-text").Text()
	info.Description = doc.Find(".about-description").First().Text()
	info.ProfileImage = doc.Find(".channel-header-profile-image").First().AttrOr("src", "")

	// stats := []string{}
	doc.Find(".about-stat").Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		info.AllStats = append(info.AllStats, t)
		c := func(str string) bool {
			return strings.Contains(t, str)
		}
		switch {
		case c("subscribers"):
			n, _ := parseNumberFromStat(t)
			info.Subscribers = n
		case c("views"):
			n, _ := parseNumberFromStat(t)
			info.Views = n
		case c("Joined"):
			info.Joined = strings.Replace(t, "Joined ", "", 1)
		}
	})

	doc.Find("[class*='related-channel'] a[title]").Each(func(_ int, s *goquery.Selection) {
		r := Link{
			Name: s.Text(),
			URL:  s.AttrOr("href", ""),
		}
		info.Related = append(info.Related, r)
	})

	doc.Find("#header-links a").Each(func(_ int, s *goquery.Selection) {
		r := Link{
			Name: s.Text(),
			URL:  s.AttrOr("href", ""),
		}
		info.HeaderLinks = append(info.HeaderLinks, r)
	})

	return
}

// parseNumberFromStat parses a number statistic from a stat
func parseNumberFromStat(stat string) (number int, err error) {
	numbers := filterStr(stat, func(s rune) bool {
		return !(s <= '0' || s >= '9')
	})
	n, err := strconv.ParseInt(numbers, 10, 64)
	return int(n), err
}
