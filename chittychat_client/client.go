package main

import (
	"bufio"
	"context"
	"log"
	pb "mp2/chittychat_proto"
	"os"

	"google.golang.org/grpc"
)

//counter til brug i lamport timestamp
var counter int = 0

func main() {
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go Publish(ctx, client, scanner.Text())
	}

}

func RecieveBroadcast(message *pb.MessageWithLamport) pb.Empty {
	log.Printf("%v %v", message.GetMessage(), message.GetTime())

	return pb.Empty{}
}

func Publish(ctx context.Context, client pb.ChittyChatClient, message string) {
	var lamportMessage = &pb.MessageWithLamport{Message: &pb.Message{Message: message}, Time: &pb.Lamport{Counter: int32(counter)}}
	client.Publish(ctx, lamportMessage)
}
