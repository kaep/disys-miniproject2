package main

import (
	"log"
	pb "mp2/chittychat_proto"

	"google.golang.org/grpc"
)

//der skal selvfølgelig laves noget lamport-dreng

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	//close connection after function return
	defer conn.Close()

	//create the client with the connection
	c := pb.NewChittyChatClient(conn)

	Publish(c)
}

func Publish(c pb.ChittyChatClient) {
	//create the message to publish
	message := pb.Message{}

}
