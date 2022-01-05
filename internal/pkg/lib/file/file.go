package fileUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	_commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = _commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := filepath.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func MkDirIfNeeded(dir string) error {
	if !FileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		return err
	}

	return nil
}
func RmDir(dir string) error {
	if FileExist(dir) {
		err := os.RemoveAll(dir)
		return err
	}

	return nil
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func AbsolutePath(pth string) string {
	if !IsAbosutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = AddPathSepIfNeeded(pth)

	return pth
}

func IsAbosutePath(pth string) bool {
	return path.IsAbs(pth) ||
		strings.Index(pth, ":") == 1 // windows
}

func AddPathSepIfNeeded(pth string) string {
	sep := consts.PthSep

	if strings.LastIndex(pth, sep) < len(pth)-1 {
		pth += sep
	}
	return pth
}

func GetFilesFromParams(arguments []string) []string {
	ret := make([]string, 0)

	for _, arg := range arguments {
		if strings.Index(arg, "-") != 0 {
			if arg == "." {
				arg = AbsolutePath(".")
			} else if strings.Index(arg, "."+consts.PthSep) == 0 {
				arg = AbsolutePath(".") + arg[2:]
			} else if !IsAbosutePath(arg) {
				arg = AbsolutePath(".") + arg
			}

			ret = append(ret, arg)
		} else {
			break
		}
	}

	return ret
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetFileName(pathOrUrl string) string {
	index := strings.LastIndex(pathOrUrl, consts.PthSep)

	name := pathOrUrl[index+1:]
	return name
}

func GetFileNameWithoutExt(pathOrUrl string) string {
	name := GetFileName(pathOrUrl)
	index := strings.LastIndex(name, ".")
	return name[:index]
}

func GetExtName(pathOrUrl string) string {
	index := strings.LastIndex(pathOrUrl, ".")

	return pathOrUrl[index:]
}
func GetDirName(pth string) (name string) {
	pth = strings.Trim(pth, consts.PthSep)
	index := strings.LastIndex(pth, consts.PthSep)
	name = pth[index:]
	name = strings.Trim(name, consts.PthSep)

	return name
}

func GetAbsolutePath(pth string) string {
	if !IsAbosutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = AddSepIfNeeded(pth)

	return pth
}

func AddSepIfNeeded(pth string) string {
	if strings.LastIndex(pth, consts.PthSep) < len(pth)-1 {
		pth += consts.PthSep
	}
	return pth
}

func GetWorkDir() string { // where we run file in
	dir, _ := os.Getwd()

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	return dir
}

func GetExeDir(workDir string) string { // where zd.exe file in
	var dir string
	arg1 := strings.ToLower(os.Args[0])

	name := filepath.Base(arg1)
	if strings.Index(name, "zd") == 0 && strings.Index(arg1, "go-build") < 0 {
		p, _ := exec.LookPath(os.Args[0])
		if strings.Index(p, string(os.PathSeparator)) > -1 {
			dir = p[:strings.LastIndex(p, string(os.PathSeparator))]
		}
	} else { // debug
		dir = workDir
	}

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	return dir
}
