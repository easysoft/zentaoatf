package service

import (
	"context"
	"fmt"
	pb "github.com/aaronchen2k/deeptest/internal/comm/grpc/proto/greater"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

type RpcService struct{}

func NewRpcService() *RpcService {
	return &RpcService{}
}

func (s *RpcService) SayHello(nodeIp string, nodePort int, req *pb.HelloRequest) (
	reply *pb.HelloReply) {

	// Set up a connection to the server.
	creds, err := credentials.NewClientTLSFromFile("cert/test.pem", "*.deeptest.com")
	if err != nil {
		logUtils.Errorf("credentials.NewClientTLSFromFile fail, error %s", err.Error())
		return
	}

	url := fmt.Sprintf("%s:%d", nodeIp, nodePort)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		logUtils.Errorf("grpc.Dial %s fail, error %s", url, err.Error())
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err = c.SayHello(ctx, req)
	if err != nil {
		logUtils.Errorf("SayHello fail, error %s", url, err.Error())
	}

	return
}
