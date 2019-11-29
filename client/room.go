package client

import (
	"time"
)

type room struct {
	Time    chan time.Time
	Name    string
	Admin   int
	Clients map[int]*Client
}

func NewRoom() *room {
	return &room{
		Clients: make(map[int]*Client),
	}
}

func (r *room) AddClient(cl *Client) {
	r.Clients[cl.Id] = cl
}

func (r *room) RemoveClient(cl *Client) {
	delete(r.Clients, cl.Id)
}

func (r *room) GetClient(id int) *Client {
	if oldcl, ok := r.Clients[id]; ok {
		return oldcl
	}
	return nil
}
