/*
* Created on 12 March 2024
* @author Sai Sumanth
 */
package main

import (
	"context"
	"fmt"
	pb "grpc-demo/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	fmt.Println("GRPC Client")
	// set up a connection to the server
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

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
