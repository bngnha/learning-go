package main

import (
	"context"
	"fmt"

	"github.com/bngnha/learn-golang/microservice/greeter/srv/proto/hello"
	micro "github.com/micro/go-micro"
)

func main() {
	service := micro.NewService()

	service.Init()

	cl := hello.NewSayService("go.micro.srv.greeter", service.Client())

	rsp, err := cl.Hello(context.Background(), &hello.Request{Name: "John"})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
