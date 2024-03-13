package main

import (
	"context"
	"fmt"
	pb "grpc-demo/proto"
	"log"
	"time"
)

func SendStreamToServer(client pb.GreetServiceClient, names *pb.NamesList) {
	fmt.Println("Client Streaming Started")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send stream of names to server %v", err)
	}

	for _, name := range names.Names {

		if err := stream.Send(&pb.HelloRequest{Name: name}); err != nil {
			log.Fatalf("Error while sending stream %v", err)
		}
		log.Println(name, "said hi to server")
		time.Sleep(time.Millisecond * 1000)
	}

	/// close the stream and receive the message from server
	res, err := stream.CloseAndRecv()
	fmt.Println("Finished streaming from client")
	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}
	log.Printf("Message received from Server %v", res.Messages)

}
