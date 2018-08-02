# youtubearchive

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
