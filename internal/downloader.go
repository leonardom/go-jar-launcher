package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type downloader struct {
	URL             string
	DestinationPath string
}

func NewDownloader(url, destination string) *downloader {
	return &downloader{
		URL:             url,
		DestinationPath: destination,
	}
}

func (d *downloader) Download() error {
	err := d.prepareDestinationPath()
	if err != nil {
		return err
	}
	log.Printf("Downloading %v to %v...\n", d.URL, d.DestinationPath)
	resp, err := http.Get(d.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	filename := path.Base(resp.Request.URL.String())
	out, err := os.Create(path.Join(d.DestinationPath, filename))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func (d *downloader) prepareDestinationPath() error {
	if _, err := os.Stat(d.DestinationPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(d.DestinationPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
