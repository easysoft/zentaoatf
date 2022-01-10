package service

import (
	"errors"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"go.uber.org/zap"
)

var (
	ErrEmpty = errors.New("请上传正确的文件")
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx iris.Context, fh *multipart.FileHeader) (iris.Map, error) {
	filename, err := GetFileName(fh.Filename)
	if err != nil {

		return nil, err
	}
	path := filepath.Join(dir.GetCurrentAbPath(), "static", "upload", "images")
	err = dir.InsureDir(path)
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("dir.InsureDir", err.Error()))
		return nil, err
	}
	_, err = ctx.SaveFormFile(fh, filepath.Join(path, filename))
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("ctx.SaveFormFile", "保存文件到本地"))
		return nil, err
	}

	return iris.Map{"local": path}, nil
}

func (s *FileService) LoadDirs(dir string) (asset serverDomain.TestAsset, err error) {
	if !fileUtils.FileExist(dir) || !fileUtils.IsDir(dir) {
		logUtils.Errorf("dir %s not exist", dir)
		return
	}

	asset = serverDomain.TestAsset{Path: dir, Title: fileUtils.GetDirName(dir), IsDir: true, Slots: iris.Map{"icon": "icon"}}
	s.GetAllChildren(dir, &asset)

	return
}

func (s *FileService) GetAllChildren(childPath string, parent *serverDomain.TestAsset) (err error) {
	if !fileUtils.IsDir(childPath) { // is file
		return
	}

	childPath = fileUtils.AddPathSepIfNeeded(fileUtils.AbsolutePath(childPath))

	list, err := ioutil.ReadDir(childPath)
	if err != nil {
		return err
	}

	for _, grandson := range list {
		name := grandson.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		childPath := childPath + name
		if grandson.IsDir() { // 目录, 递归遍历
			dirNode := s.addDir(childPath, parent)

			s.GetAllChildren(childPath, dirNode)
		}
	}

	return
}

func (s *FileService) addDir(pth string, parent *serverDomain.TestAsset) (dirNode *serverDomain.TestAsset) {
	dirNode = &serverDomain.TestAsset{Path: pth, Title: fileUtils.GetDirName(pth),
		IsDir: true, Slots: iris.Map{"icon": "icon"}}
	parent.Children = append(parent.Children, dirNode)

	return
}

// GetFileName 获取文件名称
func GetFileName(name string) (string, error) {
	fns := strings.Split(strings.TrimLeft(name, "./"), ".")
	if len(fns) != 2 {
		logUtils.Errorf("文件上传失败", zap.String("trings.Split", name))
		return "", ErrEmpty
	}
	ext := fns[1]
	md5, err := dir.MD5(name)
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("dir.MD5", err.Error()))
		return "", err
	}
	return str.Join(md5, ".", ext), nil
}

// GetPath 获取文件路径
func (s *FileService) GetPath(filename string) string {
	return filepath.Join("upload", "images", filename)
}
