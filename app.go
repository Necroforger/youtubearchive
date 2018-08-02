package youtubearchive

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Necroforger/youtubearchive/youtubedl"
	"github.com/jinzhu/gorm"
)

// ErrorGroup is a group of errors
type ErrorGroup []error

// Add an error to the ErrorGroup
func (e *ErrorGroup) Add(err error) {
	*e = append(*e, err)
}

// Error returns a string representation of this error group
func (e ErrorGroup) Error() (err string) {
	for _, v := range e {
		err += v.Error() + "; "
	}
	return
}

// AppOptions are options passed to NewApp
type AppOptions struct {
	Procs int
	Log   io.Writer
}

// App contains application data
type App struct {
	DB    *gorm.DB
	Procs int

	logMu sync.Mutex
	Log   io.Writer
}

func (a *App) log(v ...interface{}) {
	if a.Log == nil {
		return
	}

	a.logMu.Lock()
	fmt.Fprintln(a.Log, v...)
	a.logMu.Unlock()
}

func (a *App) logMessage(v ...interface{}) {
	a.log(
		append([]interface{}{"[message]: "}, v...)...,
	)
}

func (a *App) logVideo(v Video) {
	a.log("[video]: ", v.Title)
}

func (a *App) logError(e error) {
	a.log("[error]: ", e.Error())
}

// NewApp returns a pointer to a new app
func NewApp(DB *gorm.DB, opts *AppOptions) *App {
	if opts == nil {
		opts = &AppOptions{}
	}

	DB.AutoMigrate(&Video{})
	DB.AutoMigrate(&Tag{})

	return &App{
		DB:    DB,
		Procs: opts.Procs,
		Log:   opts.Log,
	}
}

// InsertVideo inserts a video into the database
func (a *App) InsertVideo(v youtubedl.Video) error {
	m := VideoFromYoutubedl(v)
	m.LastScanned = time.Now()

	var found Video
	err := a.DB.Where("video_id = ?", m.VideoID).Order("created_at DESC").First(&found).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logMessage("creating record")
			return a.DB.Create(&m).Error
		}
		return err
	}

	if !videoEqual(found, m) {
		a.logMessage("video is different from previously saved video: saving another copy")
		return a.DB.Create(&m).Error
	}

	a.logMessage("updating LastScanned for last saved video")
	found.LastScanned = time.Now()
	return a.DB.Save(&found).Error
}

// DownloadURL downloads metadata information from the given URL
func (a *App) DownloadURL(URL string) error {
	count := 1

	a.logMessage("downloading from URL: ", URL)
	a.logMessage("processes: ", a.Procs)

	// If we are using multiple youtube-dl processes, retrieve the total
	// Number of videos in the playlist so that we may split the work
	// Amongst count/procs goroutines.
	if a.Procs > 1 {
		a.logMessage("enumerating playlist videos")
		c, err := youtubedl.EnumerateVideos(URL, nil)
		if err != nil {
			return err
		}
		count = c
		a.logMessage(count, "videos found")
	}

	// Group of any errors returned by the ExtractPlaylistInfo function
	errors := ErrorGroup{}
	vidc := make(chan youtubedl.Video)
	errc := make(chan error)

	go func() {
	done:
		for {
			select {
			case v, ok := <-vidc:
				if !ok {
					break done
				}
				a.logVideo(VideoFromYoutubedl(v))
				err := a.InsertVideo(v)
				if err != nil {
					a.logError(err)
				}
			case v, ok := <-errc:
				if !ok {
					break done
				}
				a.logError(v)
				errors.Add(v)
			}
		}
		a.logMessage("done")
	}()

	a.logMessage("extracting playlist information...")
	youtubedl.ExtractPlaylistInfo(URL, &youtubedl.ExtractOpts{
		Count: count,
		Procs: a.Procs,
	}, vidc, errc)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// videoEqual returns if two videos are equal
func videoEqual(a, b Video) bool {
	return (a.Description == b.Description) &&
		(a.Title == b.Title)
}
