package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/jinzhu/gorm"
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

// downloadThumbnailCmd handles downloading thumbnails from the database
func downloadThumbnails(db *gorm.DB, directory string, procs int) {
	tx := db.Begin()
	defer tx.Commit()

	rows, err := tx.Raw("select thumbnail_url, video_id from videos group by video_id").Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		thumbURL string
		videoID  string
	)

	// Fill the semaphore with 'procs' tokens
	semaphore := make(chan struct{}, procs)
	for i := 0; i < procs; i++ {
		semaphore <- struct{}{}
	}

	for rows.Next() {
		if err := rows.Scan(&thumbURL, &videoID); err != nil {
			log.Fatal("error scanning row: ", err)
		}

		<-semaphore
		go func(videoID, thumbURL string) {
			downloadThumbnail(videoID, thumbURL, directory)
			semaphore <- struct{}{}
		}(videoID, thumbURL)
	}

	// Drain all tokens from the semaphore
	for i := 0; i < procs; i++ {
		<-semaphore
	}
}
