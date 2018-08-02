package youtubedl

import (
	"bufio"
	"encoding/json"
	"io"
	"os/exec"
	"strconv"
	"sync"
)

const (
	ytCommand       = "youtube-dl"
	ytFlatPlaylist  = "--flat-playlist"
	ytDumpJSON      = "--dump-json"
	ytPlaylistStart = "--playlist-start"
	ytPlaylistEnd   = "--playlist-end"
)

// doYoutubledl executes youtubedl and calls the supplied callback function with the
// stdout and stderr
func doYoutubedl(args []string, fn func(io.ReadCloser, io.ReadCloser)) error {
	cmd := exec.Command(ytCommand, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	fn(stdout, stderr)
	return cmd.Wait()
}

// EnumerateVideos enumerates the videos in a playlist
// URL is the URL of the playlist you would like to enumerate
// it can be a youtube playlist or youtube channel
// videos is a channel of FlatVideo. This parameter can be left as nil
func EnumerateVideos(URL string, videos chan FlatVideo) (int, error) {
	var (
		count      int
		finalError error
	)

	err := doYoutubedl([]string{ytFlatPlaylist, ytDumpJSON, URL}, func(stdin io.ReadCloser, stderr io.ReadCloser) {
		defer func() {
			if videos != nil {
				close(videos)
			}
		}()
		reader := bufio.NewReader(stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					finalError = err
				}
				break
			}

			count++

			if videos != nil {
				var fv = FlatVideo{}
				err := json.Unmarshal([]byte(line), &fv)
				if err != nil {
					finalError = err
					break
				}
				videos <- fv
			}
		}
	})
	if err != nil {
		return -1, err
	}

	return count, finalError
}

// ExtractOpts are extraction options
type ExtractOpts struct {
	// Number of videos in playlist
	Count int

	// Number of goroutines to use
	Procs int

	// Start and end specify start and end offsets of playlist
	// Extraction. Do not use this with count and procs
	Start int
	End   int
}

// NewExtractOpts returns the default extract options
func NewExtractOpts() *ExtractOpts {
	return &ExtractOpts{
		Count: -1,
		Procs: 1,
	}
}

// ExtractPlaylistInfo is like EnumerateVideos but extracts more detailed information, like the Video description and
// Thumbnail.
func ExtractPlaylistInfo(URL string, opts *ExtractOpts, videos chan Video, errc chan error) {
	if opts == nil {
		opts = NewExtractOpts()
	}

	parallel(opts.Count, opts.Procs, func(start, end int) {
		// Construct youtube-dl arguments
		args := createExtractArgs(URL, opts, start, end)
		err := doYoutubedl(args, func(stdin, stderr io.ReadCloser) {
			for {
				// Decode video JSON
				v := Video{}
				err := json.NewDecoder(stdin).Decode(&v)
				if err != nil {
					if err != io.EOF && errc != nil {
						errc <- err
					}
					return
				}
				if videos != nil {
					videos <- v
				}
			}
		})
		if err != nil && errc != nil {
			errc <- err
		}
	})

	// Close channels
	if videos != nil {
		close(videos)
	}
	if errc != nil {
		close(errc)
	}
}

func createExtractArgs(URL string, opts *ExtractOpts, start, end int) []string {
	args := []string{ytDumpJSON}
	if opts.Count > 1 {
		args = append(args, // If the playlist is a definite size, specify start and end
			ytPlaylistStart, strconv.Itoa(start+1),
			ytPlaylistEnd, strconv.Itoa(end),
		)
	}
	if opts.Start > 0 {
		args = append(args, ytPlaylistStart, strconv.Itoa(opts.Start+1))
	}
	if opts.End > 0 {
		args = append(args, ytPlaylistEnd, strconv.Itoa(opts.End))
	}
	args = append(args, URL)
	return args
}

// helper function for splitting work up into separate goroutines
func parallel(count int, procs int, fn func(int, int)) {
	var blockSize int
	if count > 0 && procs > 0 {
		blockSize = count / procs
	} else {
		blockSize = 1
	}
	if blockSize <= 0 {
		blockSize = 1
	}

	idx := count
	var wg sync.WaitGroup
	for idx > 0 {
		start := idx - blockSize
		end := idx
		if start < 0 {
			start = 0
		}
		idx -= blockSize

		wg.Add(1)
		go func() {
			fn(start, end)
			wg.Done()
		}()
	}
	wg.Wait()
}
