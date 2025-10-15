package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserMessage struct {
	UserID  string `json:"userid"`
	Message string `json:"message"`
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)

	if err != nil || data == nil {
		fmt.Println("Error reading body of HTTP request.")
		return
	}

	var um UserMessage

	json.Unmarshal(data, &um)

	fmt.Printf("User: %s - Message: %s\n", um.UserID, um.Message)

	ctx := r.Context()
	if myVal, ok := ctx.Value(traceId{}).(int64); ok {
		fmt.Println("Found TraceID: ", myVal)
	} else {
		fmt.Println("Error: No TraceID in request context.")
	}

}
