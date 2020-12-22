package service

import (
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
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
		s.Upgrade()
	}
}

func (s *UpgradeService) Upgrade() {

}
