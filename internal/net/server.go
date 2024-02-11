package net

import (
	"fmt"
	"goMud/internal/game"
	"goMud/internal/vm"
	"net"
	"os"
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":2323")
	s.listener = listener

	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}

	fmt.Println("Server listening on port 2323")
	virtialMachine := vm.GetVirtualMachine()
	go virtialMachine.Run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		connectionChannel := make(chan ConnectionCommand)
		readChannel := make(chan ConnectionRead)

		tConnection := NewConnection(conn, connectionChannel, readChannel)
		lineHandlerChannel := make(chan string)
		lineSenderChannel := make(chan string)
		tTelnet := NewTelnet(connectionChannel, readChannel, lineHandlerChannel, lineSenderChannel)

		go tConnection.HandleConnection()
		go tTelnet.HandleConnection()
		game.NewHandler(lineHandlerChannel, lineSenderChannel)
	}
}
