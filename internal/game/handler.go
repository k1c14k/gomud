package game

import (
	"goMud/internal/vm"
)

type Handler struct {
	lineChannel        chan string
	lineSendingChannel chan string
	vmHandlerObject    vm.Object
}

func (h Handler) handleLines() {
	channel := vm.GetCommandChannel()
	playerObject := newPlayerObject(h.lineSendingChannel)
	for {
		line := <-h.lineChannel
		channel <- vm.NewMethodCallCommand(h.vmHandlerObject, "HandleLine", []vm.Value{vm.NewStringValue(line)}, map[string]vm.Object{"player": playerObject})
	}
}

func NewHandler(lineChannel chan string, lineSendingChannel chan string) *Handler {
	handler := &Handler{
		lineChannel,
		lineSendingChannel,
		vm.NewObject("player_handler"),
	}
	go handler.handleLines()
	return handler
}

type playerObject struct {
	playerClass playerClass
}

func (p playerObject) GetClass() vm.Class {
	return p.playerClass
}

type playerClass struct {
	playerSendMethod playerSendMethod
}

func (p playerClass) GetStringPool() []string {
	return []string{}
}

func (p playerClass) GetMethod(_ string) vm.Method {
	return p.playerSendMethod
}

type playerSendMethod struct {
	lineChannel chan string
}

func (p playerSendMethod) Execute(ef *vm.ExecutionFrame) {
	value := ef.PopValue()
	switch value.(type) {
	case *vm.StringValue:
		p.lineChannel <- value.(*vm.StringValue).Value
	}
}

func newPlayerObject(lineChannel chan string) vm.Object {
	return playerObject{playerClass{playerSendMethod{lineChannel: lineChannel}}}
}
