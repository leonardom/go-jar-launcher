package internal

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type downloader struct {
	URL             string
	DestinationPath string
	Timeout         int
}

func NewDownloader(url, destination string, timeout int) *downloader {
	return &downloader{
		URL:             url,
		DestinationPath: destination,
		Timeout:         timeout,
	}
}

func (d *downloader) Download() error {
	err := d.prepareDestinationPath()
	if err != nil {
		return err
	}
	log.Printf("Downloading %v to %v...\n", d.URL, d.DestinationPath)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(d.Timeout))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.URL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
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
