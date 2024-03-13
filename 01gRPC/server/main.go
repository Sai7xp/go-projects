/*
* Created on 12 March 2024
* @author Sai Sumanth
 */
package main

import (
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-demo/proto"
	"log"
	"net"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	fmt.Println("GRPC Server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	log.Printf("Server listening at %v", lis.Addr())
	// Serve accepts incoming connections on the listener lis
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve & accept incoming connections on listener %v", err)
	}
}
