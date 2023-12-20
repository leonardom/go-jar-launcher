package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/leonardom/go-jar-launcher/configs"
	cp "github.com/otiai10/copy"
)

type updater struct {
	AppName        string
	Config         *configs.Config
	ServerEndpoint string
	UpdateDir      string
}

func NewUpdater(appName string, config *configs.Config) updater {
	return updater{
		AppName:        appName,
		Config:         config,
		ServerEndpoint: config.CheckUpdate,
		UpdateDir:      "./update",
	}
}

func (u *updater) CheckUpdate() {
	u.prepareUpdateDir()
	log.Printf("Checking for update on %v\n", u.ServerEndpoint)
	localChecksum := u.getLocalChecksum()
	serverChecksum := u.getServerChecksum()
	log.Printf("Local checksum:  [%v]\n", localChecksum)
	log.Printf("Server checksum: [%v]\n", serverChecksum)
	if localChecksum == serverChecksum {
		log.Printf("App is up-to-date!\n")
		return
	}
	log.Printf("New version avaiable!\n")
	err := u.archive()
	if err != nil {
		log.Printf("ERROR: Not possible to create archive: %v. Update aborted!\n", err)
		return
	}
	zipfile, err := u.download()
	if err != nil {
		log.Printf("ERROR: Not possible to download update: %v. Update aborted!\n", err)
		return
	}
	err = u.unpack(zipfile)
	if err != nil {
		log.Printf("ERROR: Not possible to unpack update: %v. Update aborted!\n", err)
		return
	}
	err = u.replaceFiles()
	if err != nil {
		log.Printf("ERROR: Not possible to apply update: %v!\n", err)
		return
	}
}

func (u *updater) prepareUpdateDir() error {
	os.RemoveAll(u.UpdateDir)
	if _, err := os.Stat(u.UpdateDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(u.UpdateDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *updater) download() (string, error) {
	updateFilename := fmt.Sprintf("%v.zip", u.AppName)
	url := fmt.Sprintf("%v/%v", u.ServerEndpoint, updateFilename)
	downloader := NewDownloader(url, u.UpdateDir)
	err := downloader.Download()
	if err != nil {
		return "", err
	}
	zipfile := path.Join(u.UpdateDir, updateFilename)
	return zipfile, nil
}

func (u *updater) archive() error {
	archiveName := fmt.Sprintf("backup-%v-%v.zip", u.AppName, getTimeStamp())
	excludeDirs := []string{"backup-", "backups", strings.Replace(u.UpdateDir, "./", "", -1), strings.Replace(u.Config.JavaHome, "./", "", -1)}
	archiver := NewArchiver(".", archiveName, excludeDirs)
	err := archiver.Archive()
	return err
}

func (u *updater) unpack(zipfile string) error {
	unpacker := NewUnpacker(zipfile, u.UpdateDir)
	err := unpacker.Unpack()
	if err != nil {
		return err
	}
	os.Remove(zipfile)
	return nil
}

func (u *updater) replaceFiles() error {
	err := cp.Copy(u.UpdateDir, "./")
	return err
}

func (u *updater) getLocalChecksum() string {
	checksumFilename := fmt.Sprintf("%v.checksum", u.AppName)
	checksum := readFile(checksumFilename)
	return strings.TrimSpace(string(checksum))
}

func (u *updater) getServerChecksum() string {
	checksumFilename := fmt.Sprintf("%v.checksum", u.AppName)
	url := fmt.Sprintf("%v/%v", u.ServerEndpoint, checksumFilename)
	downloader := NewDownloader(url, u.UpdateDir)
	err := downloader.Download()
	if err != nil {
		return ""
	}
	return readFile(path.Join(u.UpdateDir, checksumFilename))
}

func readFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func getTimeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
}
