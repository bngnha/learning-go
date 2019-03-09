package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bngnha/learn-golang/microservice/grpc/proto/hello"

	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial("localhost:8089", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer con.Close()

	c := hello.NewGreeterClient(con)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &hello.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
