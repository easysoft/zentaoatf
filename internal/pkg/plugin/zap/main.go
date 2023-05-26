package main

import (
	zapPlugin "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/plugin"
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
	"github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/shared"
	"github.com/hashicorp/go-plugin"
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
