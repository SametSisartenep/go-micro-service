package main

import (
	"fmt"
	proto "github.com/SametSisartenep/go-micro-service/proto"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
	"os"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, res *proto.HelloResponse) error {
	res.Greeting = "Hello " + req.Name
	return nil
}

func RunClient(service micro.Service) {
	greeter := proto.NewGreeterClient("greeter", service.Client())

	res, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(res)
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),

		micro.Flags(
			cli.BoolFlag{
				Name:  "client",
				Usage: "Run the client",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			if c.Bool("client") {
				RunClient(service)
				os.Exit(0)
			}
		}),
	)

	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
