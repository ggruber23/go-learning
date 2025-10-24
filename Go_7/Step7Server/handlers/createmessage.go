package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "step7server/generated6"
)

type CreateMessageHandler struct {
	DataChannel chan UserMessage
}

func (cmh CreateMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Error("Failed to read body of HTTP request.", "Error", err)
		return
	}

	var um UserMessage

	err = json.Unmarshal(data, &um)

	if err != nil {
		slog.Error("Failed to unmarshal HTTP request body into UserMessage struct.", "Error", err)
		http.Error(w, "Failed to parse HTTP request body as expected.", http.StatusBadRequest)
		/*
			Missing json field does not cause error but leads to usage of default value in the struct.
			Only body which is no valid json leads to an error here.
		*/
		return
	}

	fmt.Printf("User: %s - Message: %s\n", um.UserID, um.Message)

	// save via grpc interface (grpc server implemented in step6)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Could not create GRPC Client.", "Error", err)
	}
	defer conn.Close()
	c := pb.NewStep6RPCClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.Save(ctx, &pb.MyMessage{UserID: um.UserID, Message: um.Message})
	if err != nil {
		slog.Error("Could not save message via GRPC connection.", "Error", err)
		http.Error(w, "Server could not save message in datastore.", http.StatusInternalServerError)
	} else {
		slog.Info("Saved message via GRPC connection")

		if cmh.DataChannel == nil {
			//TODO: error
		}

		cmh.DataChannel <- um
	}

}
