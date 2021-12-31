package service

import (
	"context"
	pb "github.com/aaronchen2k/deeptest/internal/comm/grpc/proto/greater"
)

type Greeter struct {
	pb.GreeterServer
}

func (c *Greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.GetName()}, nil
}
