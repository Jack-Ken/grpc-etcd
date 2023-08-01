package client

import (
	"context"
	"fmt"
	"test/idl"
)

type DemoService struct {
	*idl.UnimplementedDemoServiceServer
}

func (*DemoService) SayHi(ctx context.Context, req *idl.HiRequest) (*idl.HiResponse, error) {
	word := fmt.Sprintf("hello, %s, this is a test!", req.Name)
	return &idl.HiResponse{
		Message: word,
	}, nil
}
