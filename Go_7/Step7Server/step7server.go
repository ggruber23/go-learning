package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"step7server/handlers"
	"sync"
	"syscall"
	"time"
)

func main() {

	// support for ctrl-c
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := handlers.NewDispatcher() // actor

	registerHandler := new(handlers.RegisterHandler)
	registerHandler.RegisterChannel = dispatcher.RegisterChannel

	createMsgHandler := new(handlers.CreateMessageHandler)
	createMsgHandler.DataChannel = dispatcher.DataChannel

	var wg4Actor sync.WaitGroup

	wg4Actor.Go(func() { dispatcher.Run(ctx) })

	msgReceiver := http.NewServeMux()

	msgReceiver.Handle("POST /createmessage", createMsgHandler)
	msgReceiver.HandleFunc("GET /getlast10", handlers.GetLast10Messages)

	msgReceiver.Handle("GET /register4messages", registerHandler)

	s := http.Server{
		Addr:         ":1234",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      msgReceiver,
	}

	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGINT {
				fmt.Println("received ctrl-c")
				cancel()  // stop actor
				s.Close() // non-graceful shutdown, causes ListenAndServe to return
				return
			}
		}
	}()

	fmt.Println("Start listening with address: ", s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe error: " + err.Error())
		if err != http.ErrServerClosed {
			fmt.Println("panic following")
			panic(err)
		}
	}

	// test showed: web socket connection(s) not yet closed here if not closed explicitly
	// fmt.Println("wait before shutdown main")
	// time.Sleep(30 * time.Second)

	// fmt.Println("wait for actor ...")
	wg4Actor.Wait()

	fmt.Println("shutdown main")

}
