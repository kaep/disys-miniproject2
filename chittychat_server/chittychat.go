package main

import (
	"context"
	"fmt"
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

func (s *Server) Publish(ctx context.Context, in *pb.MessageWithLamport) (*pb.MessageWithLamport, error) {
	//logik her
	fmt.Println("Publish kaldt på serveren")
	log.Println("Publish")
	return &pb.MessageWithLamport{Message: &pb.Message{Message: "10hi"}, Time: &pb.Lamport{Counter: int32(54)}}, nil
}

func (s *Server) Broadcast(ctx context.Context, in *pb.MessageWithLamport) (*pb.BroadcastReply, error) {
	//husk logging

	//bestem om eget timestamp er større end den der kommer ind
	var timeToReport int
	if counter > int(in.Time.Counter) {
		timeToReport = counter
	} else if counter < int(in.Time.Counter) {
		timeToReport = int(in.Time.Counter)
	}
	var message = &pb.MessageWithLamport{Message: in.GetMessage(), Time: &pb.Lamport{Counter: int32(timeToReport)}}
	//for alle klient in klienter: broadcast(besked, timestamp)
	for i := 0; i < len(clients); i++ {
		clients[i].Broadcast(ctx, message)
	}
	return &pb.BroadcastReply{}, nil
}

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
