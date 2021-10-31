package main

import (
	"context"
	"log"
	pb "mp2/chittychat_proto"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//create the client with the connection
	client := pb.NewChittyChatClient(conn)
	//client.RegisterClient(ctx, &pb.Client{}, client)

	//register klienten hos serveren

	//client.Publish(ctx, )
	//publish message, compiler happy
	//Publish(client)
	var message = &pb.MessageWithLamport{Message: &pb.Message{Message: string("Hey bro")}, Time: &pb.Lamport{Counter: int32(42)}}
	client.Publish(ctx, message)

}

func RecieveBroadcastClient(ctx context.Context, in *pb.MessageWithLamport) {
	//Denne metode kaldes fra serveren når der broadcastes, så alt logges jf. krav R4
	log.Printf("%v %v", in.GetMessage(), in.GetTime())

	//Opdater counter til serverens værdi
	counter = MaxInt(counter, int(in.GetTime().Counter))
}

func RegisterClient(ctx context.Context, client pb.ChittyChatClient) {
	client.RegisterClient(ctx, &pb.Client{})
}

//Helper method
func MaxInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
