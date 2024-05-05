package main

import (
	"fmt"
	pb "grpc-demo/proto"
	"io"
)

func (s *helloServer) SayHelloBiDirectionalStreaming(stream pb.GreetService_SayHelloBiDirectionalStreamingServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		fmt.Println("Got a new request from client with Name :", req.Name)

		if err := stream.Send(&pb.HelloResponse{Message: "Hey " + req.Name}); err != nil {
			return nil
		}
	}
}
