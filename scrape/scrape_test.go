package scrape_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Necroforger/youtubearchive/scrape"
)

func dumpJSON(i interface{}) string {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetIndent("", "  ")
	e.Encode(i)
	return string(b.Bytes())
}

func TestChannelInfo(t *testing.T) {
	info, err := scrape.GetChannelInfo("https://www.youtube.com/user/Hiiragi230/about")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", dumpJSON(info))
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
