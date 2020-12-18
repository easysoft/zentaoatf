package serverUtils

import (
	"archive/zip"
	"fmt"
	errUtils "github.com/easysoft/zentaoatf/src/utils/err"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"net/http"
	"strings"
)

func Download(url string, dst string) {
	if d, err := HTTPDownload(url); err == nil {
		logUtils.PrintTo(fmt.Sprintf("downloaded %s.\n", url))
		if WriteDownloadFile(dst, d) == nil {
			logUtils.PrintTo(fmt.Sprintf("saved %s as %s\n", url, dst))
		}
	}
}
func HTTPDownload(uri string) ([]byte, error) {
	logUtils.PrintTo(fmt.Sprintf("download file from %s.\n", uri))
	res, err := http.Get(uri)
	if err != nil {
		logUtils.PrintTo(err.Error())
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logUtils.PrintTo(err.Error())
	}
	logUtils.PrintTo(fmt.Sprintf("ReadFile: Size of download: %d\n", len(d)))
	return d, err
}

func WriteDownloadFile(dst string, d []byte) error {
	logUtils.PrintTo(fmt.Sprintf("WriteFile: Size of download: %d\n", len(d)))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		logUtils.PrintTo(err.Error())
	}
	return err
}

func GetZipSingleDir(path string) string {
	folder := ""
	z := archiver.Zip{}
	err := z.Walk(path, func(f archiver.File) error {
		if f.IsDir() {
			zfh, ok := f.Header.(zip.FileHeader)
			if ok {
				logUtils.PrintTo("file: " + zfh.Name)

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
