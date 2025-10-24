package handlers

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	pb "step7server/generated6"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetLast10Messages(w http.ResponseWriter, r *http.Request) {

	log.Println("Connection from:", r.Host)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to websocket connection.", "Error", err)
		return
	}

	// read via grpc interface (grpc server implemented in step6)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Could not create GRPC Client.", "Error", err)
	}
	defer conn.Close()
	c := pb.NewStep6RPCClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r2, err2 := c.GetLast10(ctx, &emptypb.Empty{})
	if err2 != nil {
		slog.Error("Could not read messages via GRPC connection.", "Error", err2)
		ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Could not read messages from datastore"),
			time.Time{})
		return
	}

	for _, msg := range r2.Messages {
		// send messages to websocket client
		msgStr := msg.UserID + " | " + msg.Message
		ws.WriteMessage(websocket.TextMessage, []byte(msgStr))
	}

}
