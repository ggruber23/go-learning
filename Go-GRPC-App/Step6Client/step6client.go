package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "step6client/generated6"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStep6RPCClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetLast10(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get data: %v", err)
	} else {
		log.Printf("Number of messages: %d", len(r.Messages))
	}

	_, err = c.Save(ctx, &pb.MyMessage{UserID: "1800", Message: "mountains falling,"})
	if err != nil {
		log.Fatalf("could not save message")
	} else {
		log.Printf("saved message")
	}

}
