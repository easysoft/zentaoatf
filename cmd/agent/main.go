package main

import (
	agentCron "github.com/aaronchen2k/deeptest/internal/agent/core/cron"
	agentViper "github.com/aaronchen2k/deeptest/internal/agent/core/viper"
	agentZap "github.com/aaronchen2k/deeptest/internal/agent/core/zap"
	"github.com/aaronchen2k/deeptest/internal/agent/modules/grpc"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

func main() {
	agentCron.NewAgentCron().Init()

	agentViper.Init()
	agentZap.Init()

	injectModule()
}

func injectModule() {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	grpc.NewGrpcModule()
}
