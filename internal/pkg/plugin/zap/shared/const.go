package zapShared

import (
	"github.com/hashicorp/go-plugin"
)

const (
	PluginNameZap = "ZAP"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "DEEPTEST_PLUGIN",
	MagicCookieValue: "hello",
}
