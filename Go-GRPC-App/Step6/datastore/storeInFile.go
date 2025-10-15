package datastore

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileStore struct {
	Filename   string
	fileHandle *os.File
}

func (fs *FileStore) OpenFile() bool {

	fh, err := os.OpenFile(fs.Filename, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Error opening file.", "error", err)
		return false
	}

	fs.fileHandle = fh
	return true
}

// func (fs *FileStore) CloseFile() {
// 	if fs.fileHandle != nil {
// 		fs.fileHandle.Close()
// 	}
// }

func (fs *FileStore) Close() error {
	if fs.fileHandle != nil {
		return fs.fileHandle.Close()
	} else {
		return status.Errorf(codes.Unavailable, "file was not open")
	}
}

func (fs *FileStore) AddMessage(message string) {

	if fs.fileHandle == nil {
		slog.Error("Error writing to file.", "error", "file not open")
	}

	_, err := fs.fileHandle.WriteString(message + "\n")
	if err != nil {
		slog.Error("Error writing to file.", "error", err)
		return
	}
}

func (fs *FileStore) GetLast10Messages() []string {

	if fs.fileHandle == nil {
		slog.Error("Error reading from file.", "error", "file not open")
	}

	fh := fs.fileHandle

	fh.Seek(0, 0)
	lines := make([]string, 0, 4)

	r := bufio.NewReader(fh)

	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			if len(line) != 0 {
				lines = append(lines, line)
			}
			break
		}

		if err != nil {
			fmt.Printf("error reading file %s", err)
			return nil
		}
		lines = append(lines, line)
	}

	startIdx := max(len(lines)-10, 0)

	last10lines := lines[startIdx:]

	return last10lines
}
