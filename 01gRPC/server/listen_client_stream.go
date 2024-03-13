package main

import (
	"fmt"
	pb "grpc-demo/proto"
	"io"
	"log"
)

func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string

	/// start listening to client stream
	for {
		clientReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessagesList{Messages: messages})
		}
		if err != nil {
			log.Fatalf("error while listening to client stream data")
			return err
		}
		fmt.Println("New req received from client with Name : ", clientReq.Name)
		messages = append(messages, "Hello "+clientReq.Name)
	}
}
