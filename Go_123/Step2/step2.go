package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	// support for ctrl-c
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGINT {
				fmt.Print("received ctrl-c")
				os.Exit(0)
			} else {
				fmt.Println("Caught:", sig)
			}
		}
	}()

	// functionality ...
	var nFlag = flag.Int("u", 1234, "help message for flag u(userid)")
	var msgFlag = flag.String("m", "nomsg", "help message for flag m(message)")

	flag.Parse()

	ctx := context.Background()

	ctx = context.WithValue(ctx, "TraceID", "33")

	f, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Error opening file.", "error", err)
		return
	}

	defer f.Close()

	line := strconv.Itoa(*nFlag) + " | " + *msgFlag

	_, err2 := f.WriteString(line + "\n")
	if err2 != nil {
		slog.Error("Error writing to file.", "error", err2)
		return
	}

	logInfoWithTraceID(ctx, "Write to file.", "line", line)

	f.Seek(0, 0)

	lines, shouldReturn := readLines(f)
	if shouldReturn {
		return
	}

	printLines(ctx, lines)

	for {
		fmt.Println("wait for ctrl-c")
		time.Sleep(10 * time.Second)
	}

}

func readLines(f *os.File) ([]string, bool) {

	lines := make([]string, 0, 4)

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			//slog.Info("EOF")		// goes here
			if len(line) != 0 { // does not go here
				lines = append(lines, line)
			}
			break
		}

		if err != nil {
			slog.Error("Error reading file.", "error", err)
			return nil, true
		}

		lines = append(lines, line)
	}

	return lines, false
}

func logInfoWithTraceID(ctx context.Context, msg string, args ...any) {

	if traceid, ok := ctx.Value("TraceID").(string); ok {
		args = append(args, "TraceID", traceid)
		slog.Info(msg, args...)
	} else {
		slog.Info(msg)
	}
}

// print the last 10 lines
func printLines(ctx context.Context, lines []string) int32 {

	logInfoWithTraceID(ctx, "Print the last 10 lines.")

	var count int32
	startIdx := max(len(lines)-10, 0)

	for _, line2 := range lines[startIdx:] {
		fmt.Print(line2)
		count++
	}

	return count
}
