package client

import (
	"golang.org/x/net/websocket"
)

type server struct {
	store    *store
	handlers *map[Kind]func(*Client, *Message)
}

func NewServer(handlers *map[Kind]func(*Client, *Message)) *server {
	return &server{
		store:    NewStore(),
		handlers: handlers,
	}
}

func (server *server) Handle(conn *websocket.Conn) {
	newClient(server.store).handle(conn, server.handlers)
}
