package handler

import (
	"context"

	pb "github.com/micro/agent/proto"
	"go-micro.dev/v5/genai"
)

type Agent struct{}

func New() *Agent {
	return new(Agent)
}

func (a *Agent) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	resp, err := genai.DefaultGenAI.Generate(req.Prompt)
	if err != nil {
		return err
	}
	rsp.Answer = resp.Text

	return nil
}
