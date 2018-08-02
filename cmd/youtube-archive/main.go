package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/Necroforger/youtubearchive"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	fSubscriptions = flag.String("subs", "", "path to subscriptions file")
	fSubsStart     = flag.Int("subs-start", -1, "offset of subs to start from")
	fSubsEnd       = flag.Int("subs-end", -1, "index in subs to end at")
	fOut           = flag.String("o", "archive.db", "path to output file")
	fProcs         = flag.Int("procs", 1, "number of youtube-dl processes to run when downloading")
)

// Subscription ...
type Subscription struct {
	ContentDetails struct {
		ActivityType   string `json:"activityType"`
		NewItemCount   int    `json:"newItemCount"`
		TotalItemCount int    `json:"totalItemCount"`
	} `json:"contentDetails"`
	Etag    string `json:"etag"`
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Snippet struct {
		ChannelID   string    `json:"channelId"`
		Description string    `json:"description"`
		PublishedAt time.Time `json:"publishedAt"`
		ResourceID  struct {
			ChannelID string `json:"channelId"`
			Kind      string `json:"kind"`
		} `json:"resourceId"`
		Thumbnails struct {
			Default struct {
				URL string `json:"url"`
			} `json:"default"`
			High struct {
				URL string `json:"url"`
			} `json:"high"`
			Medium struct {
				URL string `json:"url"`
			} `json:"medium"`
		} `json:"thumbnails"`
		Title string `json:"title"`
	} `json:"snippet"`
}

func readJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func ytchannel(id string) string {
	return "https://www.youtube.com/channel/" + id
}

func main() {
	flag.Parse()

	// Open database
	DB, err := gorm.Open("sqlite3", *fOut)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	app := youtubearchive.NewApp(DB, &youtubearchive.AppOptions{
		Procs: *fProcs,
		Log:   os.Stderr,
	})

	// Download subscriptions
	if *fSubscriptions != "" {
		subs := []Subscription{}
		err := readJSON(*fSubscriptions, &subs)
		if err != nil {
			log.Fatal(err)
		}

		if *fSubsStart > 0 && *fSubsStart < len(subs) {
			subs = subs[*fSubsStart:]
		}
		if *fSubsEnd > 0 && *fSubsEnd <= len(subs) {
			subs = subs[:*fSubsEnd]
		}

		for _, v := range subs {
			log.Println("Downloading channel: ", v.Snippet.Title)
			err := app.DownloadURL(
				ytchannel(v.Snippet.ResourceID.ChannelID),
			)
			if err != nil {
				log.Println("ERRORS: ", err)
			}
		}
	}

	for _, v := range flag.Args() {
		err := app.DownloadURL(v)
		if err != nil {
			log.Println("ERRORS: ", err)
		}
	}
}
