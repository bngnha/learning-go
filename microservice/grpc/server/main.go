package main

import (
	"context"
	"log"
	"net"

	hello "github.com/bngnha/learn-golang/microservice/grpc/proto/hello"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("Receive name: %v", req.Name)

	return &hello.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	l, err := net.Listen("tcp", ":8089")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	hello.RegisterGreeterServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failded to serve: %v", err)
	}
}
