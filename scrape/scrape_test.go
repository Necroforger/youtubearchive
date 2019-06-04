package scrape_test

import (
	"testing"

	"github.com/Necroforger/youtubearchive/scrape"
)

func TestChannelInfo(t *testing.T) {
	info, err := scrape.GetChannelInfo("https://www.youtube.com/user/Hiiragi230/about")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", info)
}

func TestChannelPlaylists(t *testing.T) {
	links, err := scrape.GetChannelPlaylists("https://www.youtube.com/user/Hiiragi230/playlists")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range links {
		t.Log(v.Name, v.URL)
	}
}
