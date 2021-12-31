package kvm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	pb "github.com/aaronchen2k/deeptest/internal/comm/grpc/proto/greater"
	"log"
	"net/http"
	"testing"
)

const ()

func TestHttp(t *testing.T) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(pb.HelloRequest{Name: "aaron"})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Post("https://local.deeptest.com:8086/helloworld.Greeter/SayHello", "application/json", buf)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var reply pb.HelloReply
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Greeting: %s", reply.GetMessage())
}
