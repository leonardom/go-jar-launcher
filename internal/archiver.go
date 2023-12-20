package internal

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type archiver struct {
	Source      string
	Target      string
	ExcludeDirs []string
}

func NewArchiver(source, target string, excludeDirs []string) *archiver {
	return &archiver{
		Source:      source,
		Target:      target,
		ExcludeDirs: excludeDirs,
	}
}

func (b *archiver) Archive() error {
	log.Printf("Creating backup file %v...\n", b.Target)
	err := zipFiles(b.Target, b.Source, b.ExcludeDirs)
	if err != nil {
		return err
	}
	return nil
}

func zipFiles(output string, source string, excludeDirs []string) error {
	zipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, excludeDir := range excludeDirs {
			if strings.HasPrefix(path, excludeDir) {
				return nil
			}
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Use relative path for zip file
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	}

	err = filepath.Walk(source, walkFunc)
	return err
}
