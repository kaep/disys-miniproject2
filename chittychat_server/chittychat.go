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

type Server struct {
	pb.UnimplementedChittyChatServer
}

func (s *Server) Publish(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	//logik her
	fmt.Println("Publish kaldt på serveren")
	log.Println("Publish")
	return &pb.Message{Message: "10hi"}, nil
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
