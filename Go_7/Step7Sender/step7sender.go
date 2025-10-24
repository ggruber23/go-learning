package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserMessage struct {
	UserID  string `json:"userid"`
	Message string `json:"message"`
}

func main() {

	fmt.Println("User please provide pairs of user-id and message text as prompted. Finish loop with input q.")
	usermessages := make([]UserMessage, 0, 4)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("UserID: ")
	for {
		scanner.Scan()
		strUser := scanner.Text()

		if strUser == "q" {
			break
		}

		//fmt.Println(strUser) // Println will add back the final '\n'

		fmt.Print("Message: ")
		scanner.Scan()
		strMessage := scanner.Text()

		if strMessage == "q" {
			break
		}

		//fmt.Println(strMessage) // Println will add back the final '\n'

		um := UserMessage{UserID: strUser, Message: strMessage}
		usermessages = append(usermessages, um)

		fmt.Print("UserID: ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return
	}

	//return

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for _, um := range usermessages {

		msg, err := json.Marshal(um)
		if err != nil {
			log.Println("Could not marshal UserMessage to json string")
			return
		}

		reader := strings.NewReader(string(msg))

		req, err2 := http.NewRequestWithContext(context.Background(),
			http.MethodPost, "http://localhost:1234/createmessage", reader)

		if err2 != nil {
			panic(err2)
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusOK {
			panic(fmt.Sprintf("http unexpected status: got %v", res.Status))
		}
	}

}
