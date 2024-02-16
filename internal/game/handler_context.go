package game

import "goMud/internal/vm"

type HandlerContext struct {
	player vm.ObjectValue
	room   vm.ObjectValue
}

func newHandlerContext() *HandlerContext {
	return &HandlerContext{}
}

func (h *HandlerContext) GetObjectValueFromContext(key string) *vm.ObjectValue {
	switch key {
	case "player":
		return &h.player
	case "room":
		return &h.room
	}
	return nil
}

func (h *HandlerContext) setPlayer(object *vm.Object) {
	h.player = *vm.NewObjectValue(object)
}

func (h *HandlerContext) setRoom(location *vm.Object) {
	h.room = *vm.NewObjectValue(location)
}
