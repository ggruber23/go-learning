package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Dispatcher struct {
	Conns           map[string](*websocket.Conn)
	RegisterChannel chan RegisterCommand
	DataChannel     chan UserMessage
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		Conns:           map[string](*websocket.Conn){},
		RegisterChannel: make(chan RegisterCommand, 1),
		DataChannel:     make(chan UserMessage, 1),
	}
}

func (disp *Dispatcher) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				for _, conn := range disp.Conns {
					conn.WriteControl(websocket.CloseMessage, nil, time.Time{})
				}
				fmt.Println("finished dispatcher loop")
				return
			}
		case regCmd := <-disp.RegisterChannel:
			disp.Conns[regCmd.Name] = regCmd.Connection
			fmt.Println("dispatcher: registration received: " + regCmd.Name)
		case userMsg := <-disp.DataChannel:
			{
				fmt.Println("dispatcher: message received")
				for _, conn := range disp.Conns {
					// send messages to websocket client
					msgStr := userMsg.UserID + " | " + userMsg.Message
					conn.WriteMessage(websocket.TextMessage, []byte(msgStr))
					fmt.Println("dispatcher: message sent")
				}
			}
		}
	}

}
