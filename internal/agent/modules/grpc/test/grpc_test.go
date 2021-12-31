package kvm

import (
	pb "github.com/aaronchen2k/deeptest/internal/comm/grpc/proto/greater"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"testing"
	"time"
)

const ()

func TestGrpc(t *testing.T) {
	// Set up a connection to the server.
	creds, err := credentials.NewClientTLSFromFile("cert/test.pem", "*.deeptest.com")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost:8086", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "aaron"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
