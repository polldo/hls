package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/grafov/m3u8"
)

// func download(url, dir string, threads int) (chan struct{}, int) {
// 	hlsDL := hls.New(url, nil, dir, threads, false)
// 	return hlsDL.Download()
// }

func download(url string) ([]byte, error) {
	// Get the data
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
		for _, s := range mediapl.Segments {
			if s != nil {
				fmt.Println(s.URI)
			}
		}
		fmt.Println(mediapl.Key)
		return mediapl, nil
	}

	return nil, errors.New("format not supported")
}

func open() {
	f, err := os.Open("v2.m3u")
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		for _, s := range mediapl.Segments {
			if s != nil {
				fmt.Println(s.URI)
			}
		}
		fmt.Println(mediapl.Key)
	}
}
