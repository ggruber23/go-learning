package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"step6/datastore"
	pb "step6/generated6"
)

// implement Server here

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement Step6RPCServer
type server struct {
	pb.UnimplementedStep6RPCServer
}

func (s *server) Save(ctx context.Context, msg *pb.MyMessage) (*emptypb.Empty, error) {

	fmt.Println("Save called")

	var store datastore.FileStore
	store.Filename = "messages.txt"
	if !store.OpenFile() {
		return nil, status.Errorf(codes.Internal, "could not open datastore")
	}
	defer store.Close()

	store.AddMessage(msg.UserID + " | " + msg.Message)

	return new(emptypb.Empty), nil
}

func (s *server) GetLast10(ctx context.Context, dummy *emptypb.Empty) (*pb.MyMessageList, error) {

	fmt.Println("GetLast10 called")

	var store datastore.FileStore
	store.Filename = "messages.txt"
	if !store.OpenFile() {
		return nil, status.Errorf(codes.Internal, "could not open datastore")
	}
	defer store.Close()

	lines := store.GetLast10Messages()

	output := new(pb.MyMessageList)
	output.Messages = make([]*pb.MyMessage, 0, 10)

	for _, line := range lines {
		strs := strings.Split(line, "|")
		if len(strs) == 2 {
			output.Messages = append(output.Messages, &pb.MyMessage{UserID: strs[0], Message: strs[1]})
		}
	}

	return output, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStep6RPCServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
