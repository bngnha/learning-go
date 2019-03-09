package main

import (
	"context"
	"fmt"
	"time"

	hello "github.com/bngnha/learn-golang/microservice/greeter/srv/proto/hello"
	micro "github.com/micro/go-micro"
)

// Say structure
type Say struct{}

// Hello function
func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	service.Init()
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
