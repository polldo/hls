package hls

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/grafov/m3u8"
)

type Result struct {
	Err error
}

func Download(videoURL string, outfile string) (int, chan Result, error) {
	mediaURL, err := url.Parse(videoURL)
	if err != nil {
		return 0, nil, fmt.Errorf("parsing video URL %s: %w", videoURL, err)
	}

	b, err := download(videoURL)
	if err != nil {
		return 0, nil, fmt.Errorf("downloading video: %w", err)
	}

	media, err := decode(b)
	if err != nil {
		return 0, nil, fmt.Errorf("decoding video: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(outfile), 0777)
	if err != nil {
		return 0, nil, fmt.Errorf("creating directory to store media segments: %w", err)
	}

	f, err := os.Create(outfile)
	if err != nil {
		return 0, nil, fmt.Errorf("creating video file: %w", err)
	}

	res := make(chan Result, 100)

	go func() {
		for _, s := range media.Segments {
			if s == nil {
				res <- Result{}
				continue
			}

			sURL, err := mediaURL.Parse(s.URI)
			if err != nil {
				res <- Result{Err: errors.New("URL not valid")}
				return
			}

			b, err := download(sURL.String())
			if err != nil {
				res <- Result{Err: errors.New("cannot download segment")}
				return
			}

			_, err = f.Write(b)
			if err != nil {
				res <- Result{Err: errors.New("cannot write segment to file")}
				return
			}

			res <- Result{}
		}
	}()
	return int(media.Count()), res, nil
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func decode(b []byte) (*m3u8.MediaPlaylist, error) {
	p, listType, err := m3u8.DecodeFrom(bytes.NewReader(b), false)
	if err != nil {
		return nil, fmt.Errorf("cannot decode bytes: %w", err)
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		return mediapl, nil
	}

	return nil, errors.New("format not supported")
}
