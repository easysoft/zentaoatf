package zapShared

import (
	"context"
	zapProto "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/proto"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "DEEPTEST_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginNameZap: &ZapPlugin{},
}

// Zap is the interface that we're exposing as a plugin.
type Zap interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type ZapPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Zap
}

func (p *ZapPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	zapProto.RegisterZapServer(s, &GRPCServer{Impl: p.Impl})

	return nil
}

func (p *ZapPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: zapProto.NewZapClient(c)}, nil
}
