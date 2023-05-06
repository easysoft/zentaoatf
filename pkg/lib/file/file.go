package fileUtils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/easysoft/zentaoatf/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	cp "github.com/otiai10/copy"
	"github.com/snowlyg/helper/dir"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
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
	if !IsAbsolutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = AddFilePathSepIfNeeded(pth)

	return pth
}

func IsAbsolutePath(pth string) bool {
	return path.IsAbs(pth) ||
		strings.Index(pth, ":") == 1 // windows
}

func AddUrlPathSepIfNeeded(url string) string {
	sep := "/"

	if strings.LastIndex(url, sep) < len(url)-1 {
		url += sep
	}
	return url
}

func AddFilePathSepIfNeeded(pth string) string {
	sep := consts.FilePthSep

	if strings.LastIndex(pth, sep) < len(pth)-1 {
		pth += sep
	}
	return pth
}
func RemoveFilePathSepIfNeeded(pth string) string {
	sep := consts.FilePthSep

	if strings.LastIndex(pth, sep) == len(pth)-1 {
		pth = pth[:len(pth)-1]
	}
	return pth
}

func GetFilesFromParams(arguments []string) []string {
	ret := make([]string, 0)

	for _, arg := range arguments {
		if strings.Index(arg, "-") != 0 {
			if arg == "." {
				arg = AbsolutePath(".")
			} else if strings.Index(arg, "."+consts.FilePthSep) == 0 {
				arg = AbsolutePath(".") + arg[2:]
			} else if !IsAbsolutePath(arg) {
				arg = AbsolutePath(".") + arg
			}

			ret = append(ret, arg)
		} else {
			break
		}
	}

	return ret
}

func CopyFile(src, dst string) (nBytes int64, err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file", src)
		return
	}

	source, err := os.Open(src)
	if err != nil {
		return
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return
	}
	defer destination.Close()
	nBytes, err = io.Copy(destination, source)
	return
}

func CopyFileAll(src, dst string) (nBytes int64, err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file", src)
		return
	}

	source, err := os.Open(src)
	if err != nil {
		return
	}
	defer source.Close()

	path, _ := filepath.Split(dst)
	err = dir.InsureDir(path)
	if err != nil {
		return
	}
	destination, err := os.Create(dst)
	if err != nil {
		return
	}
	defer destination.Close()
	nBytes, err = io.Copy(destination, source)
	return
}

func CopyDir(src, dest string) (err error) {
	opt := cp.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
	}
	err = cp.Copy(src, dest, opt)

	return
}

func GetFileName(pathOrUrl string) string {
	index := strings.LastIndex(pathOrUrl, consts.FilePthSep)

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

	if index < 0 {
		return ""
	}
	return pathOrUrl[index:]
}
func GetExtNameWithoutDot(pathOrUrl string) string {
	ext := GetExtName(pathOrUrl)
	ext = strings.TrimLeft(ext, ".")

	return ext
}
func GetDirName(pth string) (name string) {
	pth = strings.Trim(pth, consts.FilePthSep)
	index := strings.LastIndex(pth, consts.FilePthSep)
	if index != -1 {
		name = pth[index:]
	}
	name = strings.Trim(name, consts.FilePthSep)

	return name
}

func GetAbsolutePath(pth string) string {
	if !IsAbsolutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = AddSepIfNeeded(pth)

	return pth
}

func AddSepIfNeeded(pth string) string {
	if strings.LastIndex(pth, consts.FilePthSep) < len(pth)-1 {
		pth += consts.FilePthSep
	}
	return pth
}

func GetWorkDir() string { // where we run file in
	dir, _ := os.Getwd()

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	dir, _ = GetUserHome()
	dir += "ztf"
	dir = AddSepIfNeeded(dir)
	return dir
}

func GetZTFDir() (dir string) { // where ztf exe file in
	exeDir, _ := os.Executable()

	if commonUtils.IsRelease() { // release
		dir = filepath.Dir(exeDir)
	} else { // debug mode
		if strings.Index(strings.ToLower(exeDir), "goland") > -1 { // run with ide
			dir = os.Getenv("ZTF_CODE_DIR")
		} else {
			dir = GetWorkDir()
		}
	}

	dir, _ = filepath.Abs(dir)
	dir = AddFilePathSepIfNeeded(dir)

	return
}

func GetUserHome() (dir string, err error) {
	user, err := user.Current()
	if nil == err {
		dir = user.HomeDir
	} else { // cross compile support

		if "windows" == runtime.GOOS { // windows
			dir, err = homeWindows()
		} else { // Unix-like system, so just assume Unix
			dir, err = homeUnix()
		}
	}

	dir = AddSepIfNeeded(dir)

	return
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func Md5(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
