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
		channel <- vm.NewMethodCallCommand(h.vmHandlerObject, "HandleLine", []vm.Value{vm.NewStringValue(line)}, map[string]vm.Object{"player": *playerObject})
	}
}

func NewHandler(lineChannel chan string, lineSendingChannel chan string) *Handler {
	handler := &Handler{
		lineChannel,
		lineSendingChannel,
		*vm.NewObject("player_handler"),
	}
	go handler.handleLines()
	return handler
}

func newPlayerObject(lineChannel chan string) *vm.Object {
	class := vm.NewEmptyClass("player")
	fromClass := vm.NewObjectFromClass(*class)
	class.RegisterInternalMethod("Send", 1, 0, func(values []vm.Value) []vm.Value {
		lineChannel <- values[0].(*vm.StringValue).Value
		return []vm.Value{}
	})
	class.RegisterInternalMethod("String", 0, 1, func(values []vm.Value) []vm.Value {
		return []vm.Value{vm.NewStringValue(fromClass.String())}
	})
	return fromClass
}
