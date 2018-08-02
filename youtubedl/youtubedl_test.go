package youtubedl_test

import (
	"fmt"
	"testing"

	"github.com/Necroforger/youtubearchive/youtubedl"
)

func TestCountVideos(t *testing.T) {
	videos := make(chan youtubedl.FlatVideo)
	go func() {
		for v := range videos {
			fmt.Println(v.Title)
		}
	}()

	count, err := youtubedl.EnumerateVideos("https://www.youtube.com/user/Diremagic", videos)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	fmt.Println(count)
}

func TestGetVideoInfo(t *testing.T) {
	fmt.Println("Counting videos")
	count, err := youtubedl.EnumerateVideos("https://www.youtube.com/user/Hiiragi230/videos", nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(count, "videos found")

	videos := make(chan youtubedl.Video)
	errc := make(chan error)

	go func() {
		for v := range videos {
			fmt.Println(v.Title)
		}
	}()
	go func() {
		for v := range errc {
			fmt.Println(v)
			err = v
		}
	}()

	fmt.Println("Extracting playlist info")
	youtubedl.ExtractPlaylistInfo(
		"https://www.youtube.com/user/Hiiragi230/videos",
		&youtubedl.ExtractOpts{
			Procs: 30,
			Count: count,
		},
		videos,
		errc,
	)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
	}
	fmt.Println("Done")
}
