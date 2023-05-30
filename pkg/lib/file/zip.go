package fileUtils

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

func ZipDir(dist string, dir string) error {
	dir = RemoveFilePathSepIfNeeded(dir)

	paths := make([]string, 0)
	paths = append(paths, dir)

	zip := archiver.NewZip()
	err := zip.Archive(paths, dist)

	return err
}

func ZipFiles(dist string, dir string, files []string) error {
	dir = AddFilePathSepIfNeeded(dir)

	paths := make([]string, 0)
	for _, file := range files {
		path := dir + file
		paths = append(paths, path)
	}

	zip1 := archiver.NewZip()
	err := zip1.Archive(paths, dist)
	return err
}

func Unzip(zipPath, dstDir string) error {
	MkDirIfNeeded(filepath.Dir(dstDir))

	// open zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}
func unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
