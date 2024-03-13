package main

import (
	"context"
	pb "grpc-demo/proto"
	"io"
	"log"
)

func GetStreamDataFromServer(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Println("Names will be sent to server and server will start a stream by greeting each person")

	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Could not send names to server :%v", err)
	}

	/// start listening to stream
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error occurred while listening to server stream data %v", err)
		}
		log.Println(message)
	}

	log.Printf("Streaming Finished")
}
