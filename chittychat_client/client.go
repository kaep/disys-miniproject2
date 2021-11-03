package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	pb "mp2/chittychat_proto"
	"os"
	"time"

	"github.com/thecodeteam/goodbye"
	"google.golang.org/grpc"
)

var timestamp = 0
var name string
var id int

func main() {
	f, erro := os.OpenFile("../Logfile.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if erro != nil {
		log.Fatalf("Fejl")
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	var conn *grpc.ClientConn
	log.Print("Trying to connect to ChittyChat")
	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	//close connection after function return
	defer conn.Close()

	//new context
	ctx := context.Background()

	//create the client with the connection
	client := pb.NewChittyChatClient(conn)

	fmt.Println("---------------")
	fmt.Println("Welcome to ChittyChat")
	fmt.Println("Please enter your name")
	fmt.Println("---------------")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()

	//We want to send a message, so we increment clock
	timestamp++
	var request = &pb.ConnectionRequest{Name: name, Lamport: int32(timestamp)}
	stream, err := client.EstablishConnection(ctx, request)
	if err != nil {
		log.Fatalf("Client error %v", err)
	}
	firstmessage, err := stream.Recv()
	id = int(firstmessage.GetId())
	timestamp = MaxInt(timestamp, int(firstmessage.Lamport))
	if err != nil {
		log.Fatalf("Client error %v", err)
	}
	fmt.Println()
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("Failed to recieve message %v", err)
			}
			RecieveBroadcast(in)
		}
	}()
	fmt.Printf("You chose the name: %v", name)
	fmt.Println()
	fmt.Println("---------------")
	fmt.Println("To leave the chat service, type '/leave' at any time")
	fmt.Println("---------------")
	//for the cmd ctrl+c stuff, from https://github.com/thecodeteam/goodbye
	defer goodbye.Exit(ctx, -1)
	goodbye.Notify(ctx)

	goodbye.RegisterWithPriority(func(ctx context.Context, sig os.Signal) {
		//noget
	}, 0)
	goodbye.RegisterWithPriority(func(ctx context.Context, sig os.Signal) {
		//noget andet
	}, 1)
	goodbye.RegisterWithPriority(func(ctx context.Context, sig os.Signal) {
		var request = &pb.LeaveRequest{Id: int32(id), Lamport: int32(timestamp) + 1}
		client.Leave(ctx, request)
	}, 5)

	for scanner.Scan() {
		if scanner.Text() == "/leave" {
			var request = &pb.LeaveRequest{Id: int32(id), Lamport: int32(timestamp + 1)}
			client.Leave(ctx, request)
			conn.Close()
			os.Exit(0)
		} else {
			go Publish(ctx, client, scanner.Text())
		}
	}

}

func RecieveBroadcast(message *pb.MessageWithLamport) pb.Empty {
	timestamp = MaxInt(timestamp, int(message.GetLamport()))
	if message.Id == int32(1337) {
		var msg = fmt.Sprintf("Client %v: %v", id, message.Message)
		log.Print(msg)
	} else if message.Id == int32(id) {
		log.Printf("Client %v: Recieved own message '%v' at timestamp: %v", id, message.GetMessage(), timestamp)
	} else {
		log.Printf("Client %v: Recieved message '%v' from client %v at timestamp: %v", id, message.GetMessage(), message.Id, timestamp)
	}
	return pb.Empty{}
}

func Publish(ctx context.Context, client pb.ChittyChatClient, message string) {
	timestamp++
	var lamportMessage = &pb.MessageWithLamport{Message: message, Lamport: int32(timestamp), Id: int32(id)}
	client.Publish(ctx, lamportMessage)
}

//helper function
func MaxInt(own int, recieved int) int {
	if own >= recieved {
		return own + 1
	}
	log.Printf("Client %v logical clock updated to: %v", id, recieved+1)
	return recieved + 1
}
