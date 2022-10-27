package server

import "sync"

type rooms struct {
	rms []*room
	mu  sync.Mutex
}

func (rooms *rooms) addRoom(room *room) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	rooms.rms = append(rooms.rms, room)
}

func (rooms *rooms) removeRoom(room *room) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	for i, rm := range rooms.rms {
		if rm.name == room.name {
			rooms.rms = append(rooms.rms[:i], rooms.rms[i+1:]...)
		}
	}
}
