package fileUtils

import (
	"archive/zip"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/err"
	"github.com/mholt/archiver/v3"
	"strings"
)

func ZipFiles(dist string, dir string, files []string) error {
	dir = AddPathSepIfNeeded(dir)

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
