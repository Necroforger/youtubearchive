package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

var (
	thumbnailDirectory = flag.String("o", "thumbnails", "directory to store thumbnails in")
)

// Errors
var (
	ErrThumbnailAlreadyExists = errors.New("thumbnail already exists")
)

func getFilename(directory, videoID string) string {
	return path.Join(directory, videoID+".jpg")
}

func downloadThumbnail(videoID, URL, directory string) error {
	s, err := os.Stat(getFilename(directory, videoID))
	if !os.IsNotExist(err) || (s != nil) {
		return ErrThumbnailAlreadyExists
	}

	// Request the image from the thumbnail URL
	res, err := http.Get(URL)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		log.Println("status text was not 200: ", http.StatusText(res.StatusCode))
		return errors.New("status code not 200")
	}
	defer res.Body.Close()

	fmt.Printf("saving to [%s]\n", getFilename(directory, videoID))

	f, err := os.Create(getFilename(directory, videoID))
	defer f.Close()
	if err != nil {
		log.Println("error opening file")
		return err
	}

	_, err = io.Copy(f, res.Body)
	return err
}

func main() {
	flag.Parse()
	archiveFile := flag.Arg(0)
	if archiveFile == "" {
		log.Fatal("please the sqlite3 archive file you want to download thumbnails for")
	}

	db, err := sql.Open("sqlite3", archiveFile)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}

	var total int
	err = db.QueryRow("select count(*) from videos").Scan(&total)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT video_id, thumbnail FROM videos")
	if err != nil {
		log.Fatal("error querying database: ", err)
	}
	defer rows.Close()

	err = os.MkdirAll(*thumbnailDirectory, 0600)
	if err != nil {
		log.Fatal("error creating thumbnail directory: ", err)
	}
	var (
		thumbnail string
		videoID   string
	)
	var count int
	for rows.Next() {
		fmt.Printf("downloading [%d] of [%d] [%s] [%s]\n", count, total, videoID, thumbnail)
		err := rows.Scan(&videoID, &thumbnail)
		count++
		if err != nil {
			log.Println("error scanning: ", err)
		}
		err = downloadThumbnail(videoID, thumbnail, *thumbnailDirectory)
		if err != nil {
			if err == ErrThumbnailAlreadyExists {
				fmt.Printf("skipping, a file with ID [%s] already exists\n", videoID)
			} else {
				log.Println(err)
			}
		}
	}
}
