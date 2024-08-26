package main

import (
	"github.com/hashicorp/go-plugin"

	zapPlugin "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/plugin"
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
	zapShared "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/shared"
)

func main() {
	zapPlugin := zapPlugin.ZapPlugin{
		Impl: &zapService.ZapService{},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: zapShared.Handshake,
		Plugins: map[string]plugin.Plugin{
			zapShared.PluginNameZap: &zapPlugin,
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})
}
