package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/shared"
	"io/ioutil"

	"github.com/hashicorp/go-plugin"
)

type Zap struct{}

func (Zap) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\nWritten from plugin-go-grpc", string(value)))
	return ioutil.WriteFile("zap_"+key, value, 0644)
}

func (Zap) Get(key string) ([]byte, error) {
	return ioutil.ReadFile("zap_" + key)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: zapShared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"ZAP": &zapShared.ZapPlugin{Impl: &Zap{}},
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})
}
