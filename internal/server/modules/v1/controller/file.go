package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type FileCtrl struct {
	FileService *service.FileService `inject:""`
}

func NewFileCtrl() *FileCtrl {
	return &FileCtrl{}
}

// Upload 上传文件
func (c *FileCtrl) Upload(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	defer f.Close()

	data, err := c.FileService.UploadFile(ctx, fh)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// ListDir 列出目录
func (c *FileCtrl) ListDir(ctx iris.Context) {
	parentDir := ctx.URLParam("parentDir")

	if parentDir == "" {
		var err error
		parentDir, err = fileUtils.GetUserHome()
		if err != nil {
			return
		}
	}

	data, err := c.FileService.LoadDirs(parentDir)

	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}
