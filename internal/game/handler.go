package game

import (
	"goMud/internal/vm"
	"log"
)

type Handler struct {
	lineChannel        chan string
	lineSendingChannel chan string
	vmHandlerObject    vm.Object
	context            HandlerContext
}

func (h *Handler) handleLines() {
	channel := vm.GetCommandChannel()
	h.prepareContext()
	log.Println("Handler started")
	for {
		line := <-h.lineChannel
		channel <- vm.NewMethodCallCommand(h.vmHandlerObject, "HandleLine", []vm.Value{vm.NewStringValue(line)}, &h.context)
	}
}

func (h *Handler) prepareContext() {
	playerObject := h.newPlayerObject(h.lineSendingChannel)
	h.context.setPlayer(playerObject)
	h.context.setRoom(loadStartingLocation())
}

func NewHandler(lineChannel chan string, lineSendingChannel chan string) *Handler {
	handler := &Handler{
		lineChannel,
		lineSendingChannel,
		*vm.NewObject("player_handler"),
		*newHandlerContext(),
	}
	go handler.handleLines()
	return handler
}

func (h *Handler) newPlayerObject(lineChannel chan string) *vm.Object {
	class := vm.NewEmptyClass("<player>")
	fromClass := vm.NewObjectFromClass(*class)
	class.RegisterInternalMethod("Send", 1, 0, func(values []vm.Value) []vm.Value {
		lineChannel <- values[0].(*vm.StringValue).Value
		return []vm.Value{}
	})
	class.RegisterInternalMethod("String", 0, 1, func(values []vm.Value) []vm.Value {
		return []vm.Value{vm.NewStringValue(fromClass.String())}
	})
	class.RegisterInternalMethod("MoveTo", 1, 0, func(values []vm.Value) []vm.Value {
		room := values[0].(*vm.StringValue).Value
		h.context.setRoom(getOrInitializeRoom(room))
		return []vm.Value{}

	})
	return fromClass
}

func loadStartingLocation() *vm.Object {
	return vm.NewObject("locations/room_a")
}
