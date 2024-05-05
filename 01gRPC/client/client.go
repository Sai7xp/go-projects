/*
* Created on 12 March 2024
* @author Sai Sumanth
 */
package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-demo/proto"
	"log"
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

	///
	// MakeUnaryCall(client)

	/// server stream call
	namesList := &pb.NamesList{
		Names: []string{
			"Sai",
			"Alice",
			"Bob",
			"Jonathan",
		},
	}
	// GetStreamDataFromServer(client, namesList)

	SendStreamToServer(client, namesList)

	/// Bi Directional Streaming
	sendStreamAndReceiveStream(client)
}
