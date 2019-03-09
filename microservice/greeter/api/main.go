package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/bngnha/learn-golang/microservice/greeter/srv/proto/hello"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

// Say struct
type Say struct {
	Client hello.SayService
}

// Hello function
func (s *Say) Hello(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Receive Say.Hello API request")

	name, ok := req.Get["name"]

	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.greeter", "Name cannot be blank")
	}

	response, err := s.Client.Hello(ctx, &hello.Request{
		Name: strings.Join(name.Values, " "),
	})

	if err != nil {
		return err
	}
	rsp.StatusCode = 200

	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})

	rsp.Body = string(b)

	return nil
}

func main() {
	service := micro.NewService(micro.Name("go.micro.api.greeter"))

	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(
			&Say{Client: hello.NewSayService("go.micro.srv.greeter", service.Client())},
		),
	)
}
