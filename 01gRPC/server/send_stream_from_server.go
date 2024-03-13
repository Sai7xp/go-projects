// server streaming

package main

import (
	pb "grpc-demo/proto"
	"time"
)

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	println("Received a req from client. Let me send stream data")

	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hey " + name,
		}

		if err := stream.Send(res); err != nil {
			return err
		}
		// 2 second delay to simulate a long running process
		time.Sleep(2 * time.Second)
	}
	return nil
}
