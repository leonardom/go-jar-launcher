package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type unpacker struct {
	ZipFile         string
	DestinationPath string
}

func NewUnpacker(zipFile, destinationPath string) *unpacker {
	return &unpacker{
		ZipFile:         zipFile,
		DestinationPath: destinationPath,
	}
}

func (u *unpacker) Unpack() error {
	fmt.Printf("Unpacking file %v...\n", u.ZipFile)
	archive, err := zip.OpenReader(u.ZipFile)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, f := range archive.File {
		filePath := filepath.Join(u.DestinationPath, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(u.DestinationPath)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}
		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
