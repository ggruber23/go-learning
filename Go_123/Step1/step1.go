package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {

	var nFlag = flag.Int("u", 1234, "help message for flag u(userid)")
	var msgFlag = flag.String("m", "nomsg", "help message for flag m(message)")

	flag.Parse()

	// fmt.Printf("Hello, %d\n", *nFlag)
	// fmt.Print(*msgFlag)

	f, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer f.Close()

	s := strconv.Itoa(*nFlag) + " | " + *msgFlag + "\n"

	f.WriteString(s)

	f.Seek(0, 0)

	lines := make([]string, 0, 4)

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			if len(line) != 0 {
				//fmt.Println(line)
				lines = append(lines, line)
			}
			break
		}

		if err != nil {
			fmt.Printf("error reading file %s", err)
			return
		}
		//fmt.Print(line)
		lines = append(lines, line)
	}

	startIdx := max(len(lines)-10, 0)

	// for idx, line2 := range lines {

	// 	if idx >= startIdx {

	// 		strs := strings.Split(line2, "|")
	// 		countparts := len(strs)
	// 		if countparts == 2 {
	// 			fmt.Print(strs[1])
	// 		} else {
	// 			fmt.Print(strings.Join(strs[1:], "|"))
	// 		}

	// 		//fmt.Print(line2)
	// 	}
	// }

	for _, line2 := range lines[startIdx:] {
		fmt.Print(line2)
	}

}
