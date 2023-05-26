package main

import (
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
	"github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/shared"
	"github.com/hashicorp/go-plugin"
)

func main() {
	zapPlugin := zapShared.ZapPlugin{
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
