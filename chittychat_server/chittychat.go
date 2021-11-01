package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "mp2/chittychat_proto"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

var timestamp = 0

type Server struct {
	pb.UnimplementedChittyChatServer
	clients []pb.ChittyChat_EstablishConnectionServer //might make sense as a map
}

func main() {
	os.Remove("../Logfile") //Delete the file to ensure a fresh log for every session
	f, erro := os.OpenFile("../Logfile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if erro != nil {
		log.Fatalf("Fejl")
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
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

func (s *Server) EstablishConnection(Empty *pb.Empty, stream pb.ChittyChat_EstablishConnectionServer) error {
	//Add the stream to our stored streams
	s.clients = append(s.clients, stream)
	fmt.Println("New user joined!") //overvej at smække navne på
	//forsøger at holde liv i stream
	for {
		select {
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *Server) Broadcast(ctx context.Context, message *pb.MessageWithLamport) (*pb.Empty, error) {
	log.Printf("Broadcast kaldt på serveren med timestamp: %v", timestamp)
	for i := 0; i < len(s.clients); i++ {
		err := s.clients[i].Send(message)

		if err != nil {
			log.Print(err)
			//atm an error is logged everytime the server tries to broadcast to a client that has left -> make a leave method
		}
	}

	return &pb.Empty{}, nil
}

func (s *Server) Publish(ctx context.Context, message *pb.MessageWithLamport) (*pb.Empty, error) {
	//debugging
	//fmt.Printf("Publish kaldt på serveren: %v %v", message.GetMessage(), message.GetTime())
	//fmt.Println()
	log.Printf("Publish kaldt på server med timestamp: %v", message.GetTime().Counter)
	//update timestamp
	timestamp = MaxInt(timestamp, int(message.GetTime().Counter))
	log.Printf("Serverens timestamp opdateret til: %v", timestamp)
	//increment timestamp (modtagelse)
	timestamp++
	log.Printf("Serverens timestamp incrementet til %v pga. modtagelse ", timestamp)
	//increment timestamp (afsendelse)
	timestamp++
	log.Printf("Serverens timestamp incrementet til %v pga. afsendelse", timestamp)

	//TJEK OP PÅ OM TIMESTAMPS OPDATERES KORREKT
	var newMessage = &pb.MessageWithLamport{Message: &pb.Message{Message: message.GetMessage().Message}, Time: &pb.Lamport{Counter: int32(timestamp)}}
	//log.Printf("Publish kaldt på serveren %v", newMessage.GetTime())
	s.Broadcast(ctx, newMessage)
	return &pb.Empty{}, nil
}

//helper function
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
