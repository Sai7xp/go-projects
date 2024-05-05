package main

import (
	"context"
	"fmt"
	pb "grpc-demo/proto"
	"log"
)

func sendStreamAndReceiveStream(client pb.GreetServiceClient) {
	fmt.Println("BiDirectional Streaming Started")
	fmt.Println("Started sending stream of data to server and listening to server stream as well")

	_, err := client.SayHelloBiDirectionalStreaming(context.Background())
	if err != nil{
		log.Fatalf("cound not send names %v", err)
	}



}
