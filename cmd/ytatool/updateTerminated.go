package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Necroforger/youtubearchive/models"
	"github.com/jinzhu/gorm"
)

var (
	terminatedIndicators = []string{
		"This account has been terminated because we received multiple third-party claims of copyright infringement regarding material the user posted.",
		"This account has been terminated because",
	}
)

// query youtube and see if it responds with a 404 to determine if a channel is terminated
func testChannelTerminated(URL string) (bool, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return true, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	body := string(b)
	for _, v := range terminatedIndicators {
		if strings.Contains(body, v) {
			return true, nil
		}
	}

	return false, nil
}

func updateTerminatedCmd(db *gorm.DB) {
	rows := []models.Video{}
	err := db.Select("uploader_url, uploader").Group("uploader_url").Find(&rows).Error
	if err != nil {
		log.Fatal(err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatal("error beginning transaction: ", err)
	}

	err = tx.Exec(`
		DROP TABLE IF EXISTS terminated_channels;
		CREATE TABLE IF NOT EXISTS terminated_channels(
			ID            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			uploader      TEXT NOT NULL,
			uploader_url  TEXT NOT NULL UNIQUE,
			terminated    INTEGER NOT NULL
		)
	`).Error
	if err != nil {
		log.Fatal("error executing sql: ", err)
	}

	// Create a buffered channel of `n` allowing up to `n` http processes to be executing concurrently.
	semaphore := make(chan struct{}, *updateTerminatedProcs)
	for i := 0; i < *updateTerminatedProcs; i++ {
		semaphore <- struct{}{}
	}

	// Iterate over all the rows, scan each of them ensuring they are not terminated, and insert them into the database.
	for _, v := range rows {
		<-semaphore

		go func(v models.Video) {
			defer func() {
				semaphore <- struct{}{}
			}()

			if v.UploaderURL == "" {
				log.Println(v.Uploader, "does not have a channel url to scan")
				return
			}

			terminated, err := testChannelTerminated(v.UploaderURL)
			if err != nil {
				tx.Rollback()
				log.Fatal("could not test if channel: ", v.Uploader, " is terminated: ", err)
			}

			fmt.Printf("%t\t%s\t%s\n", terminated, v.Uploader, v.UploaderURL)

			err = tx.Exec(`
					INSERT INTO terminated_channels(uploader, uploader_url, terminated) VALUES(?, ?, ?);
				`, v.Uploader, v.UploaderURL, terminated).Error

			if err != nil {
				tx.Rollback()
				log.Fatal("could not insert record into database: ", err)
			}
		}(v)
	}

	// drain the semaphore ensuring all processes are complete
	for i := 0; i < *updateTerminatedProcs; i++ {
		<-semaphore
	}

	if tx.Commit().Error != nil {
		log.Fatal("error committing transaction: ", err)
	}
}
