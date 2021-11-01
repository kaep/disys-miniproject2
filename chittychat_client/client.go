package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	pb "mp2/chittychat_proto"
	"os"

	"google.golang.org/grpc"
)

//timestamp til brug i lamport timestamp
var timestamp = 0
var name string

func main() {
	f, erro := os.OpenFile("../Logfile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if erro != nil {
		log.Fatalf("Fejl")
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	//close connection after function return
	defer conn.Close()

	//new context
	ctx := context.Background()

	//create the client with the connection
	client := pb.NewChittyChatClient(conn)

	stream, err := client.EstablishConnection(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Klient linje 35 fejl %v", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message %v", err)
			}
			RecieveBroadcast(in)
		}
	}()
	fmt.Println("---------------")
	fmt.Println("Welcome to ChittyChat")
	fmt.Println("Please enter your name")
	fmt.Println("---------------")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()
	fmt.Printf("You chose the name: %v", name)
	fmt.Println()
	for scanner.Scan() {
		go Publish(ctx, client, scanner.Text())
	}

}

func RecieveBroadcast(message *pb.MessageWithLamport) pb.Empty {
	//log
	log.Printf("Klient har modtaget broadcast med følgende besked og timestamp: %v %v", message.GetMessage().Message, message.GetTime().Counter)
	//update timestamp
	timestamp = MaxInt(timestamp, int(message.GetTime().Counter))
	log.Printf("Timestamp opdateret til: %v", timestamp)
	return pb.Empty{}
}

func Publish(ctx context.Context, client pb.ChittyChatClient, message string) {
	timestamp++
	var lamportMessage = &pb.MessageWithLamport{Message: &pb.Message{Message: message}, Time: &pb.Lamport{Counter: int32(timestamp)}}
	log.Printf("Publish kaldt hos klient med timestamp %v: ", timestamp)
	client.Publish(ctx, lamportMessage)

}

//helper function
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func setupLog() {

}
