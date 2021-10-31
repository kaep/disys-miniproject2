package main

import (
	"context"
	"log"
	pb "mp2/chittychat_proto"

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
	ctx, cancel := context.Background()
	defer cancel()

	//create the client with the connection
	client := pb.NewChittyChatClient(conn)

}
