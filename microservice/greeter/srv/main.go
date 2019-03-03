package main

import(
	"github.com/micro/go-micro"
)

type Say struct{}

func(s *Say)Hello(ctx context.Context, req*hello.Request, rsp*hello.Response) error{
	rsp.Msg="Hello "+req.Name
	return nil
}

func main() {
service:=micro.NewService(
	micro.Name("go.micro.srv.greeter"),
	micro.RegistryTTL(time.Second*30),
	micro.RegisterInterval(time.Second*10)
)
service.Init()
hello.RegisterSayHandler()service.Server(), new(Say)
}
