package main

import (
	"context"
	"log"
	pb "mp2/chittychat_proto"
	"time"

	"google.golang.org/grpc"
)

//counter til brug i lamport
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//create the client with the connection
	client := pb.NewChittyChatClient(conn)
	//client.Publish(ctx, )
	//publish message, compiler happy
	//Publish(client)
	var message = &pb.MessageWithLamport{Message: &pb.Message{Message: string("Hey bro")}, Time: &pb.Lamport{Counter: int32(42)}}
	client.Publish(ctx, message)
}

func Publish(c pb.ChittyChatClient) {

}

//klient implementation af broadcast skal være
//at skrive besked og timestamp til log?
