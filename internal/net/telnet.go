package net

import "log"

type TelnetState int

const (
	SendDataState TelnetState = iota
	ProcessCommandState
	ProcessOptionState
	SubNegotiationState
)

type Telnet struct {
	conn_command  chan ConnectionCommand
	conn_read     chan ConnectionRead
	state         TelnetState
	buffer        []byte
	commandBuffer []byte
	line_handler  chan string
	line_sender   chan string
}

type TelnetCommandByte byte

const (
	SubNegotiationEnd  TelnetCommandByte = 240
	NoOperation        TelnetCommandByte = 241
	SubNegotiation     TelnetCommandByte = 250
	Will               TelnetCommandByte = 251
	Wont               TelnetCommandByte = 252
	Do                 TelnetCommandByte = 253
	Dont               TelnetCommandByte = 254
	InterpretAsCommand TelnetCommandByte = 255
)

type TelnetOptionByte byte

const (
	MCCP2 TelnetOptionByte = 86
)

func NewTelnet(conn_command chan ConnectionCommand, conn_read chan ConnectionRead, line_handler chan string, line_sender chan string) *Telnet {
	return &Telnet{
		conn_command:  conn_command,
		conn_read:     conn_read,
		state:         SendDataState,
		buffer:        make([]byte, 0),
		commandBuffer: make([]byte, 0),
		line_handler:  line_handler,
		line_sender:   line_sender,
	}
}

func (t Telnet) HandleConnection() {
	go t.handleAsyncRead()
	go t.handleLines()
	t.conn_command <- ConnectionCommand{
		command: SendData,
		data: []byte{
			byte(InterpretAsCommand),
			byte(Will),
			byte(MCCP2),
		},
	}
}

func (t Telnet) handleAsyncRead() {
	for {
		select {
		case read := <-t.conn_read:
			for _, b := range read.data {
				switch t.state {
				case SendDataState:
					switch b {
					case byte(InterpretAsCommand):
						t.state = ProcessCommandState
						t.commandBuffer = make([]byte, 1)
						t.commandBuffer[0] = b
					case 0, 13:
						// ignore null and carriage return
					case 10:
						// convert buffer to string and handle in handleLine method and clear buffer
						t.handleLine(string(t.buffer))
						t.buffer = make([]byte, 0)
					default:
						// add byte to buffer
						t.buffer = append(t.buffer, b)
					}
				case ProcessCommandState:
					switch b {
					case byte(Will), byte(Wont), byte(Do), byte(Dont):
						t.state = ProcessOptionState
						t.commandBuffer = append(t.commandBuffer, b)
					case byte(SubNegotiation):
						t.state = SubNegotiationState
						t.commandBuffer = append(t.commandBuffer, b)
					case byte(NoOperation):
						t.state = SendDataState
					case byte(InterpretAsCommand):
						t.state = SendDataState
						t.buffer = append(t.buffer, b)
					default:
						t.state = SendDataState
					}
				case ProcessOptionState:
					t.state = SendDataState
					t.commandBuffer = append(t.commandBuffer, b)
					t.processOption(TelnetCommandByte(t.commandBuffer[1]), TelnetOptionByte(t.commandBuffer[2]))
				case SubNegotiationState:
					// pass
				}
			}
		}
	}
}

func (t Telnet) handleLine(s string) {
	if t.line_handler != nil {
		t.line_handler <- s
	} else {
		log.Println("No line handler set for telnet, line sent:", s)
	}
}

func (t Telnet) processOption(commandByte TelnetCommandByte, optionByte TelnetOptionByte) {
	switch commandByte {
	case Will, Do:
		switch optionByte {
		case MCCP2:
			t.conn_command <- ConnectionCommand{
				command: SendData,
				data: []byte{
					byte(InterpretAsCommand),
					byte(SubNegotiation),
					byte(MCCP2),
					byte(InterpretAsCommand),
					byte(SubNegotiationEnd),
				},
			}
			t.conn_command <- ConnectionCommand{command: StartCompression}
		}
	case Wont, Dont:
		switch optionByte {
		case MCCP2:
			t.conn_command <- ConnectionCommand{command: StopCompression}
		}
	}
}

func (t Telnet) SendLine(line string) {
	t.conn_command <- ConnectionCommand{
		command: SendData,
		data:    []byte(line + "\n"),
	}
}

func (t Telnet) handleLines() {
	for {
		line := <-t.line_sender
		t.SendLine(line)
	}
}
