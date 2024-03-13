package main

import (
	"context"
	"fmt"
	pb "grpc-demo/proto"
	"log"
	"time"
)

func MakeUnaryCall(client pb.GreetServiceClient) {
	fmt.Println("--------------------- UNARY ------------------")
	// now let's contact the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.SayHello(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("count not greet the server: %v", err)
	}
	fmt.Printf("Response from grpc server after sending a req from client : %v\n", response.Message)

	fmt.Println("--------------------- UNARY CALL :END ------------------")
}
