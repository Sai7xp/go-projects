package main

import (
	"bufio"
	"log"
	"net"
)

func CreateTCPServer() {
	// net.Dial()   // Dial connects to a server
	// net.Listen() // Listen() creates a server and listens on specified PORT

	listener, err := net.Listen("tcp", TCP_PORT)
	if err != nil {
		log.Fatalln("Failed to listen on Port :  ", TCP_PORT, " Err: ", err)
	}
	log.Println("TCP listening on PORT ", TCP_PORT)
	defer listener.Close()

	// start listening to requests
	for {
		conn, err := listener.Accept()
		if err != nil {

		}
		log.Println("New client connected:", conn.RemoteAddr().String())
		go handleNewConnection(conn)
	}
}

func handleNewConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Hi\n")) // say Hi to new Client

	// One way of listening to Client messages
	// buffer := make([]byte, 1024)
	// for {
	// 	conn.Read(buffer)
	// 	log.Println("Local Addr String: ", conn.RemoteAddr().String(), "Message: ", string(buffer))
	// }

	// Another way of listening to client messages using Reader
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected: ", conn.RemoteAddr().String())
			break
		}
		log.Println("Message from ", conn.RemoteAddr().String(), " -> ", msg)
	}
}
