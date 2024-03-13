// client will request with and this server will respond with a message
package main

import (
	"context"
	pb "grpc-demo/proto"
)

func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hey Client, This is the message I want to say"}, nil
}
