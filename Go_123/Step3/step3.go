package main

import (
	"fmt"
	"net/http"
	"step3/handlers"
	"time"
)

func main() {

	msgReceiver := http.NewServeMux()

	msgReceiver.HandleFunc("POST /step3/createmessage", handlers.CreateMessageHandler)

	msgReceiverWithMiddleware := handlers.TraceIDMiddleware(msgReceiver)

	s := http.Server{
		Addr:         ":8083",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      msgReceiverWithMiddleware,
	}

	fmt.Println("Start listening with address: ", s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
