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
var idCounter = 0

type Server struct {
	pb.UnimplementedChittyChatServer
	clients []ChatClient
}

type ChatClient struct {
	id     int32
	name   string
	stream pb.ChittyChat_EstablishConnectionServer
}

func main() {
	os.Remove("../Logfile.txt") //Delete the file to ensure a fresh log for every session
	f, erro := os.OpenFile("../Logfile.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if erro != nil {
		log.Fatalf("Logfile error")
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
	log.Printf("ChittyChat server listening on %v", listen.Addr())
	log.Println("---------------")
	if err := server.Serve(listen); err != nil {
		log.Printf("Failed to serve: %v", err)
		log.Println() //overvej at droppe tomme linjer
	}
}

func (s *Server) EstablishConnection(request *pb.ConnectionRequest, stream pb.ChittyChat_EstablishConnectionServer) error {
	timestamp = MaxInt(timestamp, int(request.Lamport))
	var client ChatClient
	client.id = int32(idCounter)
	log.Printf("Client with name '%v' and id %v just joined at time: %v", request.Name, idCounter, timestamp)
	idCounter++
	client.name = request.Name
	client.stream = stream
	//Add the stream to the array of stored streams
	s.clients = append(s.clients, client)

	//Send the id to client, so that it knows of it increment lamport before sending
	var firstMessage = &pb.MessageWithLamport{Message: "Welcome", Lamport: int32(timestamp + 1), Id: int32(client.id)}
	stream.Send(firstMessage)

	//Broadcast that the new participant joined
	var formattedMessage = fmt.Sprintf("'%v' just joined the server at time %v", client.name, timestamp)
	var joinMessage = &pb.MessageWithLamport{
		Message: formattedMessage,
		Lamport: int32(timestamp),
		Id:      1337, //Magic number, not good
	}
	s.Broadcast(stream.Context(), joinMessage)

	//Keep the stream "alive"
	for {
		select {
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *Server) Broadcast(ctx context.Context, message *pb.MessageWithLamport) (*pb.Empty, error) {
	timestamp = MaxInt(timestamp, int(message.Lamport))
	var messageWithUpdatedLamport = &pb.MessageWithLamport{Message: message.Message, Lamport: int32(timestamp), Id: message.Id}
	log.Printf("Logical clock on server incremented because of call to Broadcast()")
	log.Printf("Broadcast called on the server at time: %v", timestamp)
	for i := 0; i < len(s.clients); i++ {
		err := s.clients[i].stream.Send(messageWithUpdatedLamport)
		if err != nil {
			log.Print(err)
		}
	}
	return &pb.Empty{}, nil
}

func (s *Server) Publish(ctx context.Context, message *pb.MessageWithLamport) (*pb.Empty, error) {
	log.Printf("Publish called by client %v with local timestamp %v: ", message.GetId(), message.GetLamport())
	//update timestamp
	timestamp = MaxInt(timestamp, int(message.GetLamport()))

	var newMessage = &pb.MessageWithLamport{Message: message.Message, Lamport: int32(timestamp), Id: message.Id}
	//increment own timestamp before message is sent (call to Broadcast() = local event)
	timestamp++
	log.Printf("Logical clock on server incremented because of call to Publish()")
	s.Broadcast(ctx, newMessage)
	return &pb.Empty{}, nil
}

func (s *Server) Leave(ctx context.Context, request *pb.LeaveRequest) (*pb.Empty, error) {
	timestamp = MaxInt(timestamp, int(request.Lamport))
	log.Printf("Logical clock on server incremented because of call to Leave()")
	var newArray []ChatClient
	var clientName string
	var id = request.GetId()
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i].id == id {
			//Dont add to new list
			clientName = s.clients[i].name
		} else {
			newArray = append(newArray, s.clients[i])
		}
	}
	s.clients = newArray //Smoothest way of updating the array
	log.Printf("XXXXX '%v' has left the building! (at time %v) XXXXX \n", clientName, timestamp)
	//Broadcast that the client has left
	var formattedMessage = fmt.Sprintf("XXXXX '%v' just left the server at time %v XXXXX", clientName, timestamp)
	var leaveMessage = &pb.MessageWithLamport{
		Message: formattedMessage,
		Lamport: int32(timestamp),
		Id:      1337, //Magic number, not good
	}
	s.Broadcast(ctx, leaveMessage)
	return &pb.Empty{}, nil
}

//helper function
func MaxInt(own int, recieved int) int {
	if own >= recieved {
		return own + 1
	}
	log.Printf("Servers logical clock updated to: %v", recieved+1)
	return recieved + 1
}
