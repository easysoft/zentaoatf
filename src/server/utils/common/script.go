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

func Download(uri string, dst string) error {
	logUtils.PrintTo(fmt.Sprintf("download file from %s.\n", uri))
	res, err := http.Get(uri)
	if err != nil {
		logUtils.PrintTo(err.Error())
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logUtils.PrintTo(err.Error())
	}
	logUtils.PrintTof("size of download: %d\n", len(bytes))

	err = ioutil.WriteFile(dst, bytes, 0666)
	if err != nil {
		logUtils.PrintTof("download fail, error: %s.\n", err.Error())
	} else {
		logUtils.PrintTof("download %s to %s.\n", uri, dst)
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
