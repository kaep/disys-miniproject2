package main

import (
	"bufio"
	"context"
	"log"
	pb "mp2/chittychat_proto"
	"os"

	"google.golang.org/grpc"
)

//timestamp til brug i lamport timestamp
var timestamp = 0

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
	//log
	log.Printf("Recieved broadcast: %v %v", message.GetMessage(), message.GetTime())

	//update timestamp
	timestamp = MaxInt(timestamp, int(message.GetTime().Counter))

	return pb.Empty{}
}

func Publish(ctx context.Context, client pb.ChittyChatClient, message string) {
	timestamp++
	var lamportMessage = &pb.MessageWithLamport{Message: &pb.Message{Message: message}, Time: &pb.Lamport{Counter: int32(timestamp)}}
	client.Publish(ctx, lamportMessage)
}

//helper function
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
