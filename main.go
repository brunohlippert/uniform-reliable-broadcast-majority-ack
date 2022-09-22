package main

import (
	. "MAJORITYACK/URBMarjorityAck"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify at least tree address:port!")
		fmt.Println("go run main.go chat 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001")
		fmt.Println("go run main.go chat 127.0.0.1:6001 127.0.0.1:5001 127.0.0.1:7001")
		fmt.Println("go run main.go chat 127.0.0.1:7001 127.0.0.1:6001 127.0.0.1:5001")
		fmt.Println("go run main.go fail 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001")
		fmt.Println("go run main.go fail 127.0.0.1:6001 127.0.0.1:5001 127.0.0.1:7001")
		fmt.Println("go run main.go fail 127.0.0.1:7001 127.0.0.1:6001 127.0.0.1:5001")
		fmt.Println("go run main.go stress 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001")
		fmt.Println("go run main.go stress 127.0.0.1:6001 127.0.0.1:5001 127.0.0.1:7001")
		fmt.Println("go run main.go stress 127.0.0.1:7001 127.0.0.1:6001 127.0.0.1:5001")
		return
	}

	addresses := os.Args[2:]
	teste := os.Args[1]

	fmt.Println(addresses)

	urb := URBMarjorityAck_Module{
		Req: make(chan URBMarjorityAck_Message),
		Ind: make(chan URBMarjorityAck_Message),
	}

	urb.Init(addresses[0], addresses)

	switch teste {
	case "chat":
		simpleChat(addresses, urb)
		break
	case "fail":
		simpleChatWithFail(addresses, urb)
		break
	case "stress":
		stressTest(addresses, urb)
		break
	}

}

func simpleChat(addresses []string, urb URBMarjorityAck_Module) {
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

func simpleChatWithFail(addresses []string, urb URBMarjorityAck_Module) {
	scanner := bufio.NewScanner(os.Stdin)
	var msg string

	if scanner.Scan() {
		msg = scanner.Text()
	}

	req := URBMarjorityAck_Message{
		From:    addresses[0],
		Message: msg,
		ID:      "000000",
	}
	urb.Req <- req

	fmt.Printf("Falha no processo\n")

	// Apenas para garantir que a mensagem ja foi enviada antes do programa sair
	time.Sleep(3 * time.Second)
}

func stressTest(addresses []string, urb URBMarjorityAck_Module) {
	N_MSGS := 100000

	// receptor de broadcasts
	go func() {
		i := 1
		messages := map[string]int{}
		for {
			in := <-urb.Ind
			// fmt.Printf("          Message from %v: %v\n", in.From, in.Message)
			messages[in.From] += 1

			if i%10000 == 0 {
				for j := 0; j < len(addresses); j++ {
					fmt.Printf(addresses[j] + ": " + strconv.Itoa(messages[addresses[j]]) + " |")
				}
				fmt.Printf("\n")
			}

			i++
		}
	}()

	//Sleep apenas para nao comecar a mandar msgs antes de todos os terminais estarem rodando
	time.Sleep(3 * time.Second)

	go func() {
		i := 0
		for {
			msg := "LOREM IPSUM..."

			req := URBMarjorityAck_Message{
				From:    addresses[0],
				Message: msg,
				ID:      generateId(),
			}
			urb.Req <- req

			if i > N_MSGS {
				break
			}
			i++
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
