package URBMarjorityAck

import (
	. "MAJORITYACK/BEB"
	"fmt"
	"strings"
)

type URBMarjorityAck_Message struct {
	From    string
	Message string
	ID      string
}

type URBMarjorityAck_Module struct {
	Ind chan URBMarjorityAck_Message
	Req chan URBMarjorityAck_Message

	// Usamos maps para ter um acesso constante aos dados
	Delivered map[URBMarjorityAck_Message]bool
	Pending   map[URBMarjorityAck_Message]bool
	Ack       map[URBMarjorityAck_Message]int

	Addresses []string

	beb BestEffortBroadcast_Module
	dbg bool
}

func (module *URBMarjorityAck_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ URB msg : " + s + " ]")
	}
}

func (module *URBMarjorityAck_Module) Init(address string, addresses []string) {
	module.InitD(address, addresses, true)
}

func (module *URBMarjorityAck_Module) InitD(address string, addresses []string, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init URB!")
	module.beb = BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message)}

	module.Addresses = addresses

	module.Delivered = map[URBMarjorityAck_Message]bool{}
	module.Pending = map[URBMarjorityAck_Message]bool{}
	module.Ack = map[URBMarjorityAck_Message]int{}

	module.beb.Init(address)
	module.Start()
}

func (module *URBMarjorityAck_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.beb.Ind:
				module.Deliver(BEB2URB(y))
			}
		}
	}()
}

func (module *URBMarjorityAck_Module) Broadcast(message URBMarjorityAck_Message) {
	module.Pending[message] = true
	msg := URB2BEB(message, module.Addresses)
	module.beb.Req <- msg
}

func (module *URBMarjorityAck_Module) Deliver(message URBMarjorityAck_Message) {
	module.Ack[message] += 1

	if !module.Pending[message] {
		module.Broadcast(message) // O Pending eh adicionado dentro da funcao broadcast
	}

	// Verifica se pode fazer delivery
	isPending := module.Pending[message]
	isDelivered := module.Delivered[message]
	canDeliver := module.Ack[message] > len(module.Addresses)/2

	if isPending && !isDelivered && canDeliver {
		module.Delivered[message] = true
		module.Ind <- message
	}

}

func URB2BEB(message URBMarjorityAck_Message, addresses []string) BestEffortBroadcast_Req_Message {
	// Como as camadas inferiores (BEB e PP2PLink) substituem a origem da mensagem vamos salvar
	// o remetente dentro do corpo da mensagem a ser enviado juntamente com seu ID para a nao duplicacao
	// dentro do conjunto de delivered, para o caso do usuario mandar duas mensagens iguais.
	originalMessage := message.Message + "@" + message.From + "@" + message.ID

	return BestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message:   originalMessage}

}

func BEB2URB(message BestEffortBroadcast_Ind_Message) URBMarjorityAck_Message {
	s := strings.Split(message.Message, "@")

	return URBMarjorityAck_Message{
		Message: s[0],
		From:    s[1],
		ID:      s[2],
	}
}

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message) }
	mod.Init(addresses[0])

	msg := BestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}
*/
