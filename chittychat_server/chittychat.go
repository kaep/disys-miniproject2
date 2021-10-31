package main

import (
	"log"
	pb "mp2/chittychat_proto"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

var counter = 0
var clients = make([]pb.ChittyChatClient, 5) //starter i 5, vokser måske?

type Server struct {
	pb.UnimplementedChittyChatServer
}

//server start, tid = 0
//klient jointer, tid = 0
//klient.publish() (tid+1)
//server: hey, 1 er større end 0, tid = 1
//server: broadcast(tid+1 = 2)
//klient: hey, 2 er større end min 1

//start og lyt :)
func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v, %v", port, err)
		log.Println()
	}
	server := grpc.NewServer()
	pb.RegisterChittyChatServer(server, &Server{})
	log.Printf("Server listening at %v", listen.Addr())
	log.Println()
	if err := server.Serve(listen); err != nil {
		log.Printf("Failed to serve: %v", err)
		log.Println() //overvej at droppe tomme linjer
	}
}
