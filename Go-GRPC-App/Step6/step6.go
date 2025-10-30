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
	fileStore *datastore.FileStore
}

func (s *server) Save(ctx context.Context, msg *pb.MyMessage) (*emptypb.Empty, error) {

	fmt.Println("Save called")

	if s.fileStore == nil || !s.fileStore.IsOpen() {
		return nil, status.Errorf(codes.Internal, "filestore not open")
	}

	s.fileStore.AddMessage(msg.UserID + " | " + msg.Message)

	return new(emptypb.Empty), nil
}

func (s *server) GetLast10(ctx context.Context, dummy *emptypb.Empty) (*pb.MyMessageList, error) {

	fmt.Println("GetLast10 called")

	if s.fileStore == nil || !s.fileStore.IsOpen() {
		return nil, status.Errorf(codes.Internal, "filestore not open")
	}

	lines := s.fileStore.GetLast10Messages()

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var fs = new(datastore.FileStore)
	fs.Filename = "messages.txt"
	fs.OpenFile()
	defer fs.Close()

	rpcServer := grpc.NewServer()
	pb.RegisterStep6RPCServer(rpcServer, &server{fileStore: fs})
	log.Printf("server listening at %v", listener.Addr())
	if err := rpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
