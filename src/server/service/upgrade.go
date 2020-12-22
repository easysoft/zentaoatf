package service

import (
	"fmt"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/mholt/archiver/v3"
	"strconv"
	"strings"
)

var ()

type UpgradeService struct {
}

func NewUpgradeService() *UpgradeService {
	return &UpgradeService{}
}

func (s *UpgradeService) CheckUpgrade() {
	pth := vari.AgentLogDir + "version.txt"
	serverUtils.Download(serverConst.AgentUpgradeURL, pth)

	content := strings.TrimSpace(fileUtils.ReadFile(pth))
	version, _ := strconv.ParseFloat(content, 64)
	if vari.Config.Version < version {
		s.Upgrade(version)
	}
}

func (s *UpgradeService) Upgrade(ver float64) (err error) {
	version := fmt.Sprintf("%.1f", ver)

	os := commonUtils.GetOs()
	if commonUtils.IsWin() {
		os = fmt.Sprintf("%s%d", os, strconv.IntSize)
	}
	url := fmt.Sprintf(serverConst.AgentDownloadURL, version, os)

	dir := vari.AgentLogDir + version
	pth := dir + ".zip"
	err = serverUtils.Download(url, pth)

	if err == nil {
		fileUtils.MkDirIfNeeded(dir)
		archiver.Unarchive(pth, dir)
	}

	return
}
