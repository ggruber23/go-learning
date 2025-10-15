package main

import (
	"context"
	"fmt"
	"testing"
)

func Test_printLines_1(t *testing.T) {
	fmt.Println("test printLines 1")

	msgs := []string{"msg1\n", "msg2\n", "msg3\n", "msg4\n", "msg5\n"}

	result := printLines(context.Background(), msgs)

	if result != 5 {
		t.Error("incorrect result: expected 5 got", result)
	}
}

func Test_printLines_2(t *testing.T) {
	fmt.Println("test printLines 2")

	msgs := []string{"msg1\n", "msg2\n", "msg3\n", "msg4\n", "msg5\n", "msg6\n", "msg7\n", "msg8\n", "msg9\n", "msg10\n", "msg11\n", "msg12\n"}

	result := printLines(context.Background(), msgs)

	if result != 10 {
		t.Error("incorrect result: expected 10 got", result)
	}
}
