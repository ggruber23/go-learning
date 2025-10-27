package handlers

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type RegisterHandler struct {
	number          int // used to create artifical names for the websocket connections
	RegisterChannel chan RegisterCommand
}

func (rh *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("Connection from:", r.Host)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to websocket connection.", "Error", err)
		return
	}

	if rh.RegisterChannel == nil {
		fmt.Println("RegisterHandler: RegisterChannel not set")
		return //TODO: error
	}

	rh.number++
	rh.RegisterChannel <- RegisterCommand{Name: strconv.Itoa(rh.number), Connection: ws}

}
