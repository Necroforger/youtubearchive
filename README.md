# youtubearchive

<!-- TOC -->

- [youtubearchive](#youtubearchive)
	- [TODO](#todo)
	- [Installing](#installing)
	- [example usage](#example-usage)
	- [Downloading all videos from your subscribed channels](#downloading-all-videos-from-your-subscribed-channels)
	- [Flags](#flags)
- [Server](#server)
	- [Installing](#installing-1)
	- [Usage example](#usage-example)
	- [Flags](#flags-1)
	- [Screenshots](#screenshots)

<!-- /TOC -->

Archive youtube video metadata of channels to an SQL database.

## TODO
- Scraping user subscriptions
- Scraping user playlists
- Downloading image thumbnails as base64

## Installing
`go get -u github.com/Necroforger/youtubearchive/cmd/youtube-archive`

## example usage

`youtube-archive -procs 30 https://www.youtube.com/user/Diremagic/videos`

This will download the metadata (title, description etc...) for all videos uploaded by Diremagic to 
a sqlite3 database.

The number of simultaneous processes is 30, so the playlist downloading will be split over 30 or less
youtube-dl processes each downloading at a different position.
I recommend using multiple processes because it speeds things up a lot.

on subsequent calls using the same database, if any of the videos contain a different title or description, it will save them, otherwise the `last_scanned` column is updated to the current time.


## Downloading all videos from your subscribed channels
To download all of the videos from the channels you are subscribed to you need to get your youtube from google takeout.

You can do so here 

https://takeout.google.com/settings/takeout?pli=1

Make sure your selected format for subscriptions is json
![img](https://i.imgur.com/foAUN8t.png)

Then run `youtube-archive -subs path/to/subscriptions.json`

## Flags
| Flag       | Description                                    | Default      |
|------------|------------------------------------------------|--------------|
| subs       | subscriptions json file                        | ""           |
| subs-start | start index in subscriptions file              | -1           |
| subs-end   | end index in subscriptions file                | -1           |
| o          | name of the output Sqlite3 file                | "archive.db" |
| procs      | number of parallel youtube-dl processes to run | 1            |


# Server

## Installing
`go get -u github.com/Necroforger/youtubearchive/cmd/youtube-archive-server`

## Usage example
```sh
yadir="$GOPATH/src/github.com/Necroforger/youtubearchive"
templates="$yadir/templates"
static="$yadir/static"
dbfile="archive.db" # Change this to your database file
addr="0.0.0.0:80"   # bind to any ip on port 80

go run youtube-archive-server \ 
-templates "$templates" -static "$static" -db "$dbfile" -addr "$addr"
```

This will host a webserver on port 80 serving the database.
You may also want to consider using an sqlite browser such as
https://sqlitebrowser.org/ to view and query the saved data.

Once its running connect to http://localhost:80 in your browser

## Flags

| Flag      | Description                                            |
|-----------|--------------------------------------------------------|
| db        | Path to database file                                  |
| templates | directory in which templates are stored                |
| static    | directory containing static files (default "./static") |
| addr      | address to bind to (default ":80")                     |

## Screenshots

![img](https://i.imgur.com/vH34Q7u.png)
![img](https://i.imgur.com/z12m6u4.png)
![img](https://i.imgur.com/I7pcv2u.png)