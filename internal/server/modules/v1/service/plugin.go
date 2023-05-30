package service

import (
	"errors"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zapPlugin "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/plugin"
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
	"github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/shared"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	"github.com/hashicorp/go-plugin"
	"path/filepath"
	"strings"
)

const (
	ZapPath = "/Users/aaron/rd/project/zentao/go/ztf/internal/pkg/plugin/zap-plugin"
)

type PluginService struct {
	zapClient    *plugin.Client
	zapRpcClient plugin.ClientProtocol
}

func (s *PluginService) Exec() (err error) {
	return
}

func (s *PluginService) Cancel() (err error) {
	return
}

func (s *PluginService) Start() (err error) {
	s.zapClient = plugin.NewClient(&plugin.ClientConfig{
		Plugins: map[string]plugin.Plugin{
			zapShared.PluginNameZap: &zapPlugin.ZapPlugin{},
		},
		Cmd:              shellUtils.GetCmd(ZapPath),
		HandshakeConfig:  zapShared.Handshake,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	s.zapRpcClient, err = s.zapClient.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	// Request the plugin
	raw, err := s.zapRpcClient.Dispense(zapShared.PluginNameZap)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	zapService := raw.(zapService.ZapInterface)

	err = zapService.Put("key", []byte("Set Msg"))
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	result, err := zapService.Get("key")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	fmt.Println(string(result))

	return
}

func (s *PluginService) Stop() (err error) {
	s.zapClient.Kill()
	return
}

func (s *PluginService) Install(req commDomain.PluginInstallReq) (err error) {
	name := req.Name
	version := req.Version

	if version == "" {
		version = "latest"
	}

	url := fmt.Sprintf(commConsts.ZapDownloadUrl, version)

	zipPath, extractDir := s.getPluginDownloadPath(name, version)

	err = fileUtils.Download(url, zipPath)
	if err != nil {
		return
	}

	err = fileUtils.Download(url+".md5", zipPath+".md5")
	if err != nil {
		return
	}

	success := s.checkMd5(zipPath)
	if !success {
		err = errors.New("md5 check fail")
		return
	}

	fileUtils.Unzip(zipPath, extractDir)

	return
}

func (s *PluginService) Uninstall() (err error) {

	return
}

func (s *PluginService) getPluginDownloadPath(name, version string) (zipPath, binDir string) {
	dir := filepath.Join(commConsts.WorkDir, commConsts.PluginDir, name)

	zipPath = filepath.Join(dir, commConsts.DownloadDir, fmt.Sprintf("%s-%s.zip", name, version))
	binDir = filepath.Join(dir, commConsts.BinDir)

	return
}

func (s *PluginService) checkMd5(filePth string) (pass bool) {
	md5Path := filePth + ".md5"

	if !fileUtils.FileExist(filePth) {
		return false
	}

	expectVal := fileUtils.ReadFile(md5Path)
	actualVal, _ := fileUtils.GetMd5(filePth)

	return strings.TrimSpace(actualVal) == strings.TrimSpace(expectVal)
}
