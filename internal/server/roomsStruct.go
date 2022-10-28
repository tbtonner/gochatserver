package server

import "sync"

type rooms struct {
	rms []*room
	mu  sync.Mutex
}

// method that adds a new room to the rooms struct 
func (rooms *rooms) addRoom(room *room) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	rooms.rms = append(rooms.rms, room)
}

// method that removes a room from the rooms struct
func (rooms *rooms) removeRoom(room *room) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	for i, rm := range rooms.rms {
		if rm.name == room.name {
			rooms.rms = append(rooms.rms[:i], rooms.rms[i+1:]...)
		}
	}
}
