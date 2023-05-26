package zapPlugin

import (
	"context"
	zapProto "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/proto"
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
)

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type ZapPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin

	// Concrete implementation, written in Go. This is only used for plugins that are written in Go.
	Impl zapService.ZapInterface
}

func (p *ZapPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	zapProto.RegisterZapServer(s, &GRPCServer{Impl: p.Impl})

	return nil
}

func (p *ZapPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{Client: zapProto.NewZapClient(c)}, nil
}
