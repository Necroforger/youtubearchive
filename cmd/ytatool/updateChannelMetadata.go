package main

import (
	"fmt"
	"log"

	"github.com/Necroforger/youtubearchive"

	"github.com/jinzhu/gorm"
)

// Channel holds channel information
type Channel struct {
	Name string
	URL  string
}

func updateChannelMetadataCmd(db *gorm.DB) {
	var channels []Channel
	rows, err := db.Raw("SELECT uploader, uploader_url FROM terminated_channels WHERE not terminated").Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Scam all results into array of channels
	for rows.Next() {
		var t = Channel{}
		err = rows.Scan(&t.Name, &t.URL)
		if err != nil {
			log.Fatal(err)
		}
		channels = append(channels, t)
	}

	// Archive each channel's metadata
	for _, v := range channels {
		fmt.Println("archiving: ", v.Name)
		youtubearchive.ArchiveChannelMetadata(db, v.URL)
	}
}
