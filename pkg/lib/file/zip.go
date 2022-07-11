package fileUtils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	errUtils "github.com/easysoft/zentaoatf/pkg/lib/err"
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

func GetZipSingleDir(path string) string {
	folder := ""
	z := archiver.Zip{}
	err := z.Walk(path, func(f archiver.File) error {
		if f.IsDir() {
			zfh, ok := f.Header.(zip.FileHeader)
			if ok {
				fmt.Println("file: ", zfh.Name)

				if folder == "" && zfh.Name != "__MACOSX" {
					folder = zfh.Name
				} else {
					if strings.Index(zfh.Name, folder) != 0 {
						return errUtils.New("found more than one folder")
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return folder
}

func Unzip(zipPath, dstDir string) error {
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
