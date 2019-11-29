package client

type store struct {
	Clients map[int]*Client
	Rooms   map[string]*room
}

func NewStore() *store {
	return &store{
		Clients: make(map[int]*Client),
		Rooms:   make(map[string]*room),
	}
}

func (store *store) AddRoom(r *room) {
	store.Rooms[r.Name] = r
}

func (store *store) RemoveRoom(name string) {
	delete(store.Rooms, name)
}

func (store *store) GetRoom(name string) *room {
	if oldroom, ok := store.Rooms[name]; ok {
		return oldroom
	}
	return nil
}
