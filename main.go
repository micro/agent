package main

import (
	"github.com/micro/agent/handler"
	pb "github.com/micro/agent/proto"
	"go-micro.dev/v5"
)

func main() {
	service := micro.New("agent")

	service.Init()

	pb.RegisterAgentHandler(service.Server(), handler.New())

	service.Run()
}
