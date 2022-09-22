package main

import (
	. "MAJORITYACK/URBMarjorityAck"
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify at least tree address:port!")
		fmt.Println("go run main.go 127.0.0.1:5001  127.0.0.1:6001 127.0.0.1:7001")
		fmt.Println("go run main.go 127.0.0.1:6001  127.0.0.1:5001 127.0.0.1:7001")
		fmt.Println("go run main.go 127.0.0.1:7001  127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	urb := URBMarjorityAck_Module{
		Req: make(chan URBMarjorityAck_Message),
		Ind: make(chan URBMarjorityAck_Message),
	}

	urb.Init(addresses[0], addresses)

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
			}

			req := URBMarjorityAck_Message{
				From:    addresses[0],
				Message: msg,
				ID:      generateId(),
			}
			urb.Req <- req
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-urb.Ind
			fmt.Printf("          Message from %v: %v\n", in.From, in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}

func generateId() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, 20)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}
