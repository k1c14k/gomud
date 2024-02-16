package game

import "goMud/internal/vm"

type RoomRepository struct {
	rooms map[string]*vm.Object
}

var instance *RoomRepository

func getRoomRepository() *RoomRepository {
	if instance == nil {
		instance = &RoomRepository{
			rooms: make(map[string]*vm.Object),
		}
	}
	return instance
}

func getOrInitializeRoom(name string) *vm.Object {
	r := getRoomRepository()
	room, ok := r.rooms[name]
	if !ok {
		room = vm.NewObject(name)
		r.rooms[name] = room
	}
	return room
}
