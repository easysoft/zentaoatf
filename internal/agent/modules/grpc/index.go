package grpc

import (
	agentConsts "github.com/aaronchen2k/deeptest/internal/agent/consts"
	"github.com/aaronchen2k/deeptest/internal/agent/modules/grpc/service"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"google.golang.org/grpc/credentials"
	"net"

	// [...]
	pb "github.com/aaronchen2k/deeptest/internal/comm/grpc/proto/greater"
	"google.golang.org/grpc"
)

func NewGrpcModule() (grpcServer *grpc.Server) {
	listen, err := net.Listen("tcp", agentConsts.CONFIG.System.Addr)
	if err != nil {
		logUtils.Errorf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("cert/test.pem", "cert/test.key")
	if err != nil {
		logUtils.Errorf("failed to generate credentials %v", err.Error())
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	//opts = append(opts, grpc.UnaryInterceptor(interceptor))

	grpcServer = grpc.NewServer(opts...)
	grpcServer = grpc.NewServer(grpc.Creds(creds))

	// register service
	pb.RegisterGreeterServer(grpcServer, &service.Greeter{})
	logUtils.Infof("start agent")
	grpcServer.Serve(listen)

	return
}
