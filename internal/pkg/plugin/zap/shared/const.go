package zapShared

import (
	zapPlugin "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/plugin"
	"github.com/hashicorp/go-plugin"
)

const (
	PluginNameZap = "ZAP"
)

var PluginMap = map[string]plugin.Plugin{
	PluginNameZap: &zapPlugin.ZapPlugin{},
}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "DEEPTEST_PLUGIN",
	MagicCookieValue: "hello",
}
